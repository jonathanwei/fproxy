package fileserver

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/golang/glog"
)

type Handler struct {
	Base     string
	ListDirs bool

	// TODO: Configure URL signing.
}

func httpError(rw http.ResponseWriter, code int) {
	http.Error(rw, http.StatusText(code), code)
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

// TODO: Use actual stylesheet.
var dirTmpl = template.Must(template.New("dir").Parse(`
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
		<style>
			* {
				box-sizing: border-box;
			}
			body {
				font-family: monospace;
			}
		</style>
	</head>
	<body>
		{{- range .Files -}}
			<a href="{{.}}">{{.}}</a><br>
		{{- else -}}
			<strong>No files.</strong>
		{{- end -}}
	</body>
</html>`))

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

	if !h.ListDirs {
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

	finfos, err := f.Readdir(-1)
	if err != nil {
		return err
	}

	var files []string
	for _, chinfo := range finfos {
		files = append(files, chinfo.Name())
	}

	data := struct {
		Title string
		Files []string
	}{
		title,
		files,
	}

	var b bytes.Buffer
	err = dirTmpl.Execute(&b, data)
	if err != nil {
		return err
	}

	_, err = io.Copy(rw, &b)
	return err
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
