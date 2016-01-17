package main

import (
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	pb "github.com/jonathanwei/fproxy/proto"
)

var feCommand = cli.Command{
	Name:      "frontend",
	Aliases:   []string{"fe"},
	Usage:     "Run the frontend.",
	ArgsUsage: "[path/to/config/proto]",
	Description: strings.TrimSpace(`
The default path for the config proto is "~/.config/fproxy/frontend.textproto".

Supported fields in the config are listed below.

tcp_proxy_route: this repeated field lets you configure the frontend to perform
TCP proxying. Each clause must set a single listening address and set a target
dialing address for proxying. A sample clause:
  tcp_proxy_route {
    listen: ":8080"
    dial: "example.com:80"
  }

http_addr: this field lets you set the address to listen on for serving HTTP
requests. A sample clause:
  http_addr: ":8000"
`),
	Action: func(c *cli.Context) {
		runFe(defaultConfigPath(c, "frontend.textproto"))
	},
}

func runFe(configPath string) {
	var config pb.FrontendConfig
	readConfig(configPath, &config)

	var wg sync.WaitGroup
	defer wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runTCPProxy(config.TcpProxyRoute)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runTestClient(config.BackendAddr)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runHttpServer(config.HttpAddr)
	}()
}

func runTestClient(backendAddr string) {
	conn, err := grpc.Dial(backendAddr, grpc.WithInsecure())
	if err != nil {
		glog.Warningf("Couldn't connect to backend: %v", err)
	}
	defer conn.Close()

	client := pb.NewBackendClient(conn)

	for {
		resp, err := client.Hello(context.Background(), &pb.HelloRequest{Thingy: "client"})
		if err != nil {
			glog.Warningf("Got error from server: %v", err)
		} else {
			glog.Infof("Got response from server: %v", resp)
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func runHttpServer(serverAddr string) {
	mux := http.NewServeMux()
	mux.Handle("/", &feHandler{})
	glog.Warning(http.ListenAndServe(serverAddr, mux))
}

type feHandler struct {
}

func (f *feHandler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	io.WriteString(rw, "Hello World!")
}
