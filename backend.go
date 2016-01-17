package main

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
)

var beCommand = cli.Command{
	Name:      "backend",
	Aliases:   []string{"be"},
	Usage:     "Run the backend.",
	ArgsUsage: "[path/to/config/proto]",
	Description: strings.TrimSpace(`
The default path for the config proto is "~/.config/fproxy/backend.textproto".

Supported fields in the config are listed below.

server_addr: this field lets you set the address to listen on. A sample clause:
  server_addr: ":10000"
`),
	Action: func(c *cli.Context) {
		runBe(defaultConfigPath(c, "backend.textproto"))
	},
}

func runBe(configPath string) {
	var config pb.BackendConfig
	readConfig(configPath, &config)

	l, err := net.Listen("tcp", config.ServerAddr)
	if err != nil {
		glog.Fatalf("Failed to listen on %v: %v", config.ServerAddr, err)
	}
	defer l.Close()

	glog.Infof("Listening for requests on %v", l.Addr())

	s := grpc.NewServer()
	pb.RegisterBackendServer(s, &backendServer{})
	err = s.Serve(l)
	if err != nil {
		glog.Fatalf("Failed to serve on %v: %v", config.ServerAddr, err)
	}
}

type backendServer struct{}

func (b *backendServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	node, err := getNodeFromPath(req.Path)
	if err != nil {
		return nil, err
	}

	return &pb.GetNodeResponse{Node: node}, nil
}

func getNodeFromPath(path string) (*pb.Node, error) {
	path = filepath.Clean("/" + path)

	// TODO: take the base directory from the config instead of always using ".".
	path = filepath.Join(".", path)

	if lg := glog.V(2); lg {
		lg.Infof("Fetching path %q", path)
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	finfo, err := f.Stat()
	if err != nil {
		return nil, err
	}

	node, err := getNodeFromFileInfo(finfo)
	if err != nil || node.Kind == pb.Node_FILE {
		return node, err
	}

	chinfos, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}

	for _, chinfo := range chinfos {
		chnode, err := getNodeFromFileInfo(chinfo)
		if err != nil {
			glog.Warningf("Got error traversing dir: %v", err)
			continue
		}

		node.Child = append(node.Child, chnode)
	}

	return node, nil
}

func getNodeFromFileInfo(finfo os.FileInfo) (*pb.Node, error) {
	if finfo.Mode().IsRegular() {
		return &pb.Node{
			Name:      finfo.Name(),
			Kind:      pb.Node_FILE,
			SizeBytes: finfo.Size(),
		}, nil
	}

	if finfo.Mode().IsDir() {
		return &pb.Node{
			Name: finfo.Name(),
			Kind: pb.Node_DIR,
		}, nil
	}

	return nil, fmt.Errorf("File wasn't a dir or a regular file: %v", finfo.Mode())
}
