package main

import (
	"net/http"
	"path"
	"path/filepath"
	"text/template"

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

const tmplText = `
<!doctype html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Node.Name}}</title>
	</head>
	<body>
		{{if ne .ParentDir ""}}
			<a href="{{.ParentDir}}">
				{{if eq .ParentDir "/"}}
					..
				{{else}}
					{{.ParentDir}}
				{{end}}
			</a>/{{.Node.Name}}
		{{end}}
		{{if eq .Node.Kind 1}}
			({{.Node.SizeBytes}} bytes)
		{{else if eq .Node.Kind 2}}
			{{range .Node.Child}}
				<div style="margin-left: 16px">
					<a href="{{pathClean (print $.Path "/" .Name)}}">{{.Name}} ({{.Kind}})</a>
				</div>
			{{else}}
				<div><strong>no children</strong></div>
			{{end}}
		{{end}}
	</body>
</html>`

var tmpl = template.Must(template.New("listing").
	Funcs(template.FuncMap{
		"pathClean": path.Clean,
	}).
	Parse(tmplText))

func (f *feHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	beReq := &pb.GetNodeRequest{
		Path: path.Clean("/" + req.URL.Path),
	}
	res, err := f.client.GetNode(context.TODO(), beReq)
	if err != nil {
		glog.Warningf("GetNode RPC failed. Request: %v, error: %v", beReq, err)
		http.Error(rw, "Failed.", http.StatusInternalServerError)
		return
	}

	parentDir := filepath.Dir(beReq.Path)

	// If we're already at the root, then there is no parent dir.
	if parentDir == beReq.Path {
		parentDir = ""
	}

	tmplData := struct {
		Node      *pb.Node
		Path      string
		ParentDir string
	}{res.Node, beReq.Path, parentDir}

	err = tmpl.Execute(rw, tmplData)
	if err != nil {
		glog.Infof("Failed to write HTTP response: %v", err)
	}
}
