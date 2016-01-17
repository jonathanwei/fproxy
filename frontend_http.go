package main

import (
	"net/http"
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
	beReq := &pb.GetNodeRequest{}
	res, err := f.client.GetNode(context.TODO(), &pb.GetNodeRequest{})
	if err != nil {
		glog.Warningf("GetNode RPC failed. Request: %v, error: %v", beReq, err)
		http.Error(rw, "Failed.", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(rw, res)
	if err != nil {
		glog.Infof("Failed to write HTTP response: %v", err)
	}
}
