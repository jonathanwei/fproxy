package main

import (
	"net/http"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	goauth2 "google.golang.org/api/oauth2/v2"

	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
)

var oauthCfg = &oauth2.Config{
	Endpoint: google.Endpoint,

	// TODO: get this from config.
	RedirectURL: "http://localhost:8000/oauth2Callback",

	Scopes: []string{"email"},
}

func runHttpServer(serverAddr string, client pb.BackendClient) {
	crypter := newCookieCrypter([]byte("0000000000000000")) // TODO: get this from config

	mux := http.NewServeMux()
	mux.Handle("/", &feHandler{
		crypter: crypter,
		client:  client,
	})
	mux.Handle("/oauth2Callback", oauthHandler{crypter})
	glog.Warning(http.ListenAndServe(serverAddr, mux))
}

type feHandler struct {
	crypter cookieCrypter

	client pb.BackendClient
}

func (f *feHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ctx := context.Background()

	// Enable end-to-end cancellation.
	if c, ok := rw.(http.CloseNotifier); ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithCancel(ctx)

		go func() {
			select {
			case <-ctx.Done():
			case <-c.CloseNotify():
				cancel()
			}
		}()
	}

	cookie := f.crypter.GetAuthCookie(req)
	if cookie == nil {
		url := oauthCfg.AuthCodeURL(req.URL.Path)
		http.Redirect(rw, req, url, http.StatusFound)
		return
	}
	ctx = WithAuthCookie(ctx, cookie)

	http.FileServer(grpcFs{ctx, f.client}).ServeHTTP(rw, req)
}

type oauthHandler struct {
	crypter cookieCrypter
}

func (o oauthHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	code := req.FormValue("code")

	t, err := oauthCfg.Exchange(context.TODO(), code)
	if err != nil {
		glog.Warningf("Got error: %v", err)
		panic("foo")
	}

	httpClient := oauthCfg.Client(context.TODO(), t)

	oauth2Service, err := goauth2.New(httpClient)
	if err != nil {
		glog.Warningf("Got error getting http client: %v", err)
		panic("bar")
	}

	tokInfo, err := oauth2Service.
		Tokeninfo().
		Context(context.TODO()).
		AccessToken(t.AccessToken).
		Do()
	if err != nil {
		glog.Warningf("Got error getting token info: %v", err)
		panic("baz")
	}

	// TODO: check VerifiedEmail.

	o.crypter.SetAuthCookie(rw, &pb.AuthCookie{User: tokInfo.Email})
	http.Redirect(rw, req, req.FormValue("state"), http.StatusFound)
}
