package main

import (
	"net/http"
	"text/template"

	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
)

func runHttpServer(serverAddr string) {
	mux := http.NewServeMux()
	mux.Handle("/", &feHandler{})
	glog.Warning(http.ListenAndServe(serverAddr, mux))
}

type feHandler struct {
}

var tempData = &pb.GetNodeResponse{
	Node: &pb.Node{
		Name: "foo",
		Kind: pb.Node_DIR,
		Child: []*pb.Node{
			{
				Name:      "bar",
				Kind:      pb.Node_FILE,
				SizeBytes: 1000,
			},
			{
				Name: "baz",
				Kind: pb.Node_DIR,
			},
		},
	},
}

const tmplText = `
<!doctype html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Node.Name}}</title>
	</head>
	<body>
		{{range .Node.Child}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
	</body>
</html>`

var tmpl = template.Must(template.New("listing").Parse(tmplText))

func (f *feHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	err := tmpl.Execute(rw, tempData)
	if err != nil {
		glog.Infof("Failed to write HTTP response: %v", err)
	}
}
