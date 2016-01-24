package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	goauth2 "google.golang.org/api/oauth2/v2"

	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
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
	crypter := cookieCrypter{
		aead:     NewAEADOrDie(config.AuthCookieKey),
		insecure: config.AuthCookieInsecure,
	}

	var oauthCfg = &oauth2.Config{
		ClientID:     config.OauthConfig.ClientId,
		ClientSecret: config.OauthConfig.ClientSecret,

		Endpoint: google.Endpoint,

		RedirectURL: config.OauthConfig.RedirectUrl,

		Scopes: []string{"email"},
	}

	backendURL, err := url.Parse(config.GetBackend().Addr)
	if err != nil {
		glog.Fatalf("Backend url is invalid: %v", err)
	}

	backendProxy := httputil.NewSingleHostReverseProxy(backendURL)

	// TODO: configure client certs on outbound connection to backend.
	if t := config.GetBackend().GetTls(); t != nil {
		backendProxy.Transport = &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			TLSClientConfig:     FrontendClientTLSConfigOrDie(t),
			TLSHandshakeTimeout: 10 * time.Second,
		}
	} else if config.GetBackend().GetInsecure() {
		PrintClientInsecureWarning()
	} else {
		glog.Fatalf("The config must specify one of 'insecure' or 'tls'")
	}

	mux := http.NewServeMux()
	mux.Handle("/", &feHandler{
		crypter:  crypter,
		oauthCfg: oauthCfg,
		backend:  backendProxy,
	})
	mux.Handle("/oauth2Callback", oauthHandler{
		crypter:       crypter,
		cfg:           oauthCfg,
		emailToUserId: config.EmailToUserId,
	})
	mux.HandleFunc("/unauthorized", func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
	})

	return mux
}

type feHandler struct {
	crypter  cookieCrypter
	oauthCfg *oauth2.Config

	// TODO: make this a list of backends.
	backend http.Handler
}

func (f *feHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	cookie := f.crypter.GetAuthCookie(req)
	if cookie == nil {
		// TODO: Sign state.
		state := req.URL.Path
		url := f.oauthCfg.AuthCodeURL(state)
		http.Redirect(rw, req, url, http.StatusFound)
		return
	}

	req.Header.Add("User", cookie.User)
	f.backend.ServeHTTP(rw, req)
}

type oauthHandler struct {
	crypter       cookieCrypter
	cfg           *oauth2.Config
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
		// TODO: verify 'state' signature.
		http.Redirect(rw, req, req.FormValue("state"), http.StatusFound)
		return
	}

	http.Redirect(rw, req, "/unauthorized", http.StatusFound)
}
