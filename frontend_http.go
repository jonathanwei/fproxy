package main

import (
	"crypto/tls"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	goauth2 "google.golang.org/api/oauth2/v2"

	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
)

func runHttpServer(config *pb.FrontendConfig, client pb.BackendClient) {
	crypter := cookieCrypter{
		aead: MustNewAEAD(config.AuthCookieKey),
	}

	var oauthCfg = &oauth2.Config{
		ClientID:     config.OauthConfig.ClientId,
		ClientSecret: config.OauthConfig.ClientSecret,

		Endpoint: google.Endpoint,

		RedirectURL: config.OauthConfig.RedirectUrl,

		Scopes: []string{"email"},
	}

	mux := http.NewServeMux()
	mux.Handle("/", &feHandler{
		crypter:  crypter,
		client:   client,
		oauthCfg: oauthCfg,
	})
	mux.Handle("/oauth2Callback", oauthHandler{
		crypter:       crypter,
		cfg:           oauthCfg,
		emailToUserId: config.EmailToUserId,
	})
	mux.HandleFunc("/unauthorized", func(rw http.ResponseWriter, req *http.Request) {
		http.Error(rw, "Unauthorized", http.StatusUnauthorized)
	})

	server := &http.Server{
		Addr:    config.HttpAddr,
		Handler: mux,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	glog.Fatal(server.ListenAndServeTLS(config.HttpCertFile, config.HttpKeyFile))
}

type feHandler struct {
	crypter  cookieCrypter
	client   pb.BackendClient
	oauthCfg *oauth2.Config
}

func (f *feHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	// Enable end-to-end cancellation.
	if c, ok := rw.(http.CloseNotifier); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)
		ch := c.CloseNotify()

		go func() {
			select {
			case <-ctx.Done():
			case <-ch:
				cancel()
			}
		}()
	}

	cookie := f.crypter.GetAuthCookie(req)
	if cookie == nil {
		// TODO: Sign state.
		state := req.URL.Path
		url := f.oauthCfg.AuthCodeURL(state)
		http.Redirect(rw, req, url, http.StatusFound)
		return
	}
	ctx = WithAuthCookie(ctx, cookie)

	http.FileServer(grpcFs{ctx, f.client}).ServeHTTP(rw, req)
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
