package fileserver

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/golang/glog"
)

type DirPolicy int

const (
	None DirPolicy = iota
	AllowListing
	AllowArchives
)

type Handler struct {
	Base      string
	DirPolicy DirPolicy

	// TODO: Configure URL signing.
}

func httpError(rw http.ResponseWriter, code int) {
	http.Error(rw, http.StatusText(code), code)
}

func (h *Handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	err := h.serveHTTP(rw, req)
	if os.IsNotExist(err) {
		glog.Warningf("Path was not found: %v", err)
		http.NotFound(rw, req)
		return
	} else if err != nil {
		glog.Warningf("Got error %v for path %v.", err, req.URL.Path)
		httpError(rw, http.StatusInternalServerError)
		return
	}
}

func open(path string) (*os.File, os.FileInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}

	finfo, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}

	return f, finfo, nil
}

func (h *Handler) serveHTTP(rw http.ResponseWriter, req *http.Request) error {
	reqPath := path.Join(h.Base, path.Clean("/"+req.URL.Path))

	f, finfo, err := open(reqPath)
	if err != nil {
		return err
	}

	if !finfo.IsDir() {
		http.ServeContent(rw, req, finfo.Name(), finfo.ModTime(), f)
		return nil
	}

	if h.DirPolicy < AllowListing {
		glog.Warningf("Disallowing ls on %v.", reqPath)
		http.NotFound(rw, req)
		return nil
	}

	title := finfo.Name()
	if reqPath == h.Base {
		title = "/"
	}

	if !strings.HasSuffix(req.URL.Path, "/") {
		rw.Header().Set("Location", finfo.Name()+"/")
		rw.WriteHeader(http.StatusMovedPermanently)
		return nil
	}

	r := getRenderer(getRenderType(h.DirPolicy, req.URL), reqPath)
	return render(r, rw, title, reqPath)
}
