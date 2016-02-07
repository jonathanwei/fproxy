package fileserver

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
)

type renderer interface {
	Init(w io.Writer, title string)
	Walk(path string, finfo os.FileInfo, err error) error
	Flush() error
}

func render(r renderer, w io.Writer, title, reqPath string) error {
	r.Init(w, title)

	err := filepath.Walk(reqPath, r.Walk)
	if err != nil {
		return err
	}

	return r.Flush()
}

type relativeRenderer struct {
	renderer

	root string
}

func (r *relativeRenderer) Walk(path string, finfo os.FileInfo, err error) error {
	if r.root != "" {
		return r.renderer.Walk(strings.TrimPrefix(path, r.root), finfo, err)
	}

	if err != nil {
		return err
	}

	r.root = path + "/"
	return nil
}

type renderType int

const (
	renderTypeHTML renderType = iota
	renderTypeTAR
	renderTypeZIP
)

func getRenderType(d DirPolicy, url *url.URL) renderType {
	if d < AllowArchives {
		return renderTypeHTML
	}

	mode := url.Query().Get("m")
	if mode == "" {
		return renderTypeHTML
	}

	switch mode {
	case "html":
		return renderTypeHTML
	case "tar":
		return renderTypeTAR
	case "zip":
		return renderTypeZIP
	default:
		glog.Warningf("Unknown renderer %v", mode)
		return renderTypeHTML
	}
}

func getRenderer(r renderType, root string) renderer {
	switch r {
	case renderTypeHTML:
		return &relativeRenderer{renderer: &htmlRenderer{}}
	case renderTypeTAR:
		return nil
	case renderTypeZIP:
		return &relativeRenderer{renderer: &zipRenderer{root: root}}
	default:
		err := fmt.Errorf("Unknown renderer %v", r)
		glog.Error(err)
		panic(err)
	}
}
