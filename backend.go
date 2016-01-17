package main

import (
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

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

func (b *backendServer) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	glog.Infof("Got request: %v", req)
	return &pb.HelloResponse{Greeting: "Hello " + req.Thingy}, nil
}

func (b *backendServer) GetNode(ctx context.Context, req *pb.GetNodeRequest) (*pb.GetNodeResponse, error) {
	return nil, grpc.Errorf(
		codes.Unimplemented,
		"Unimplemented.",
	)
}
