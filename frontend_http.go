package main

import (
	"crypto/cipher"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	goauth2 "google.golang.org/api/oauth2/v2"

	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
	"github.com/unrolled/secure"
)

func runHttpServer(config *pb.FrontendConfig) {
	srvConfig := config.GetServer()

	l, err := net.Listen("tcp", srvConfig.Addr)
	if err != nil {
		glog.Fatalf("Failed to listen on %v: %v", srvConfig.Addr, err)
	}
	defer l.Close()

	glog.Infof("Listening for requests on %v", l.Addr())

	server := &http.Server{
		Handler: getFrontendHTTPMux(config),
	}

	if t := srvConfig.GetTls(); t != nil {
		l = tls.NewListener(l, FrontendTLSConfigOrDie(t))
	} else if srvConfig.GetInsecure() {
		PrintServerInsecureWarning()
	} else {
		glog.Fatalf("The config must specify one of 'insecure' or 'tls'")
	}

	glog.Fatal(server.Serve(l))
}

func getFrontendHTTPMux(config *pb.FrontendConfig) http.Handler {
	hostname := config.GetServer().Hostname
	if hostname == "" {
		glog.Fatal("Must provide a hostname for the frontend.")
	}

	crypter := cookieCrypter{
		aead:     NewAEADOrDie(config.AuthCookieKey),
		insecure: config.AuthCookieInsecure,
	}

	oauthAEAD := NewAEADOrDie(config.GetOauthConfig().StateKey)

	var oauthCfg = &oauth2.Config{
		ClientID:     config.OauthConfig.ClientId,
		ClientSecret: config.OauthConfig.ClientSecret,

		Endpoint: google.Endpoint,

		RedirectURL: config.OauthConfig.RedirectUrl,

		Scopes: []string{"email"},
	}

	backendMux := http.NewServeMux()
	portUpdateMux := http.NewServeMux()

	var backendPaths []string
	for _, backendCfg := range config.Backend {
		if strings.Contains(backendCfg.Name, "/") {
			glog.Fatalf("Backend name must not contain slashes: %v", backendCfg.Name)
		} else if backendCfg.Name == "" {
			glog.Fatal("Backend name must be non-empty.")
		}

		transport := http.DefaultTransport
		if t := backendCfg.GetTls(); t != nil {
			transport = &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSClientConfig:     FrontendClientTLSConfigOrDie(t),
				TLSHandshakeTimeout: 10 * time.Second,
			}
		} else if backendCfg.GetInsecure() {
			PrintClientInsecureWarning()
		} else {
			glog.Fatalf("The config must specify one of 'insecure' or 'tls'")
		}

		backendProxyHandler := &backendProxyHandler{
			config:    backendCfg,
			transport: transport,
		}

		portHandler := &portHandler{
			b:    backendProxyHandler,
			aead: NewAEADOrDie(backendCfg.PortKey),
		}

		backendPath := "/" + backendCfg.Name + "/"
		backendMux.Handle(backendPath, http.StripPrefix(backendPath, backendProxyHandler))
		backendPaths = append(backendPaths, backendPath)

		// TODO: Find out why we can't have two levels of StripPrefix.
		portUpdateMux.Handle(hostname+"/portupdate"+backendPath, http.StripPrefix("/portupdate"+backendPath, portHandler))
	}

	authedHandler := func(next http.Handler) http.Handler {
		return &authHandler{
			crypter:   crypter,
			oauthCfg:  oauthCfg,
			oauthAEAD: oauthAEAD,
			next:      next,
		}
	}

	mux := http.NewServeMux()
	mux.Handle(hostname+"/", authedHandler(&feHandler{
		backendPaths: backendPaths,
		backendMux:   backendMux,
	}))
	mux.Handle(hostname+"/oauth2Callback", oauthHandler{
		crypter:       crypter,
		cfg:           oauthCfg,
		aead:          oauthAEAD,
		emailToUserId: config.EmailToUserId,
	})
	mux.Handle(hostname+"/portupdate/", portUpdateMux)
	mux.HandleFunc(hostname+"/unauthorized", func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
	})

	secureMiddleware := secure.New(secure.Options{
		STSSeconds:            60 * 60 * 24 * 365, // One year.
		STSIncludeSubdomains:  true,
		STSPreload:            true,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: "default-src 'self'; style-src 'self' 'unsafe-inline'",
		IsDevelopment:         config.GetServer().GetInsecure(),
	})
	return secureMiddleware.Handler(mux)
}

type backendProxyHandler struct {
	mu   sync.Mutex
	next http.Handler

	port      int32
	config    *pb.FrontendConfig_Backend
	transport http.RoundTripper
}

func (b *backendProxyHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	b.mu.Lock()
	next := b.next
	b.mu.Unlock()

	if next == nil {
		http.Error(rw, "Backend unavailable", http.StatusInternalServerError)
		return
	}

	next.ServeHTTP(rw, req)
}

func (b *backendProxyHandler) UpdatePort(port int32) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.port == port {
		return
	}
	b.port = port

	backendURL, err := url.Parse(fmt.Sprintf("%v:%v", b.config.Host, b.port))
	if err != nil {
		glog.Fatalf("Backend url is invalid: %v", err)
	}

	backendProxy := httputil.NewSingleHostReverseProxy(backendURL)
	backendProxy.Transport = b.transport
	b.next = backendProxy
}

type portHandler struct {
	b    *backendProxyHandler
	aead cipher.AEAD
}

func (p *portHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var update pb.PortUpdate
	ok := DecryptProto(p.aead, req.URL.Path, nil, &update)
	if !ok {
		glog.Warningf("Got malformed port update: %v", req.URL.Path)
		http.Error(rw, "Error", http.StatusInternalServerError)
		return
	}

	p.b.UpdatePort(update.Port)

	fmt.Fprintf(rw, "Alive.\n")
	for {
		time.Sleep(2 * time.Minute)

		_, err := fmt.Fprintf(rw, "Still alive!\n")
		if err != nil {
			return
		}

		if f, ok := rw.(http.Flusher); ok {
			f.Flush()
		}
	}
}

type authHandler struct {
	crypter   cookieCrypter
	oauthCfg  *oauth2.Config
	oauthAEAD cipher.AEAD
	next      http.Handler
}

func (a *authHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	cookie := a.crypter.GetAuthCookie(req)
	if cookie == nil {
		state := EncryptProto(a.oauthAEAD, &pb.OAuthState{Path: req.URL.Path}, nil)
		url := a.oauthCfg.AuthCodeURL(state)
		http.Redirect(rw, req, url, http.StatusFound)
		return
	}

	req.Header.Add("User", cookie.User)
	req.Header.Del("Cookie")
	a.next.ServeHTTP(rw, req)
}

type feHandler struct {
	backendPaths []string
	backendMux   http.Handler
}

func (f *feHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	reqPath := req.URL.Path
	for _, backendPath := range f.backendPaths {
		if strings.HasPrefix(reqPath, backendPath) ||
			reqPath == backendPath[:len(backendPath)-1] {
			f.backendMux.ServeHTTP(rw, req)
			return
		}
	}

	if reqPath != "/" {
		http.NotFound(rw, req)
		return
	}
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(rw, "<pre>\n")
	for _, path := range f.backendPaths {
		fmt.Fprintf(rw, "<a href=\"%s\">%s</a>\n", path[1:], path[1:])
	}
	fmt.Fprintf(rw, "</pre>\n")
}

type oauthHandler struct {
	crypter       cookieCrypter
	cfg           *oauth2.Config
	aead          cipher.AEAD
	emailToUserId map[string]string
}

func (o oauthHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	code := req.FormValue("code")

	t, err := o.cfg.Exchange(ctx, code)
	if err != nil {
		glog.Warningf("Got error: %v", err)
		panic("foo")
	}

	httpClient := o.cfg.Client(ctx, t)

	oauth2Service, err := goauth2.New(httpClient)
	if err != nil {
		glog.Warningf("Got error getting http client: %v", err)
		panic("bar")
	}

	tokInfo, err := oauth2Service.
		Tokeninfo().
		Context(ctx).
		AccessToken(t.AccessToken).
		Do()
	if err != nil {
		glog.Warningf("Got error getting token info: %v", err)
		panic("baz")
	}

	if uid, ok := o.emailToUserId[tokInfo.Email]; ok && tokInfo.VerifiedEmail {
		o.crypter.SetAuthCookie(rw, &pb.AuthCookie{User: uid})
		var state pb.OAuthState
		redirectPath := "/"
		if DecryptProto(o.aead, req.FormValue("state"), nil, &state) {
			redirectPath = state.Path
		}
		http.Redirect(rw, req, redirectPath, http.StatusFound)
		return
	}

	http.Redirect(rw, req, "/unauthorized", http.StatusFound)
}
