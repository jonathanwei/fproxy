package main

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
)

func runHttpServer(serverAddr string, client pb.BackendClient) {
	mux := http.NewServeMux()
	mux.Handle("/", &feHandler{client: client})
	glog.Warning(http.ListenAndServe(serverAddr, mux))
}

type feHandler struct {
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

	http.FileServer(grpcFs{ctx, f.client}).ServeHTTP(rw, req)
}
