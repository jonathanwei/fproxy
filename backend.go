package main

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"

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

TODO: rewrite this.

server: this field is used to store all configuration for the RPC server that
we run. This includes the address it's located at, the certificate it presents
to clients (and the associated private key), and the CA it uses to verify
client certificates. A sample clause:
  server {
    addr: ":10000"
    cert_file: "/certs/backend/backend.pem"
    key_file: "/certs/backend/backend-key.pem"
    client_ca_file: "/certs/frontend/ca.pem"
  }

serve_path: this required field lets you set the root path to serve files from.
A sample clause:
  serve_path: "/home/foo/bar/"
`),
	Action: func(c *cli.Context) {
		runBe(defaultConfigPath(c, "backend.textproto"))
	},
}

func runBe(configPath string) {
	var config pb.BackendConfig
	readConfig(configPath, &config)

	if config.ServePath == "" {
		glog.Fatal("No serving path was given in config.serve_path.")
	}
	srvConfig := config.GetServer()

	glog.Infof("Serving path %v", config.ServePath)

	l, err := net.Listen("tcp", srvConfig.Addr)
	if err != nil {
		glog.Fatalf("Failed to listen on %v: %v", srvConfig.Addr, err)
	}
	defer l.Close()

	glog.Infof("Listening for requests on %v", l.Addr())

	server := &http.Server{
		Handler: getBackendHTTPMux(&config),
	}

	if t := srvConfig.GetTls(); t != nil {
		l = tls.NewListener(l, BackendTLSConfigOrDie(t))
	} else if srvConfig.GetInsecure() {
		PrintServerInsecureWarning()
	} else {
		glog.Fatalf("The config must specify one of 'insecure' or 'tls'")
	}

	glog.Fatal(server.Serve(l))
}

func getBackendHTTPMux(config *pb.BackendConfig) http.Handler {
	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir(config.ServePath))
	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		glog.Infof("Got request from user: %v", req.Header.Get("User"))
		fs.ServeHTTP(rw, req)
	})

	return mux
}

func PrintInsecureWarning(msg string) {
	warningMsg := "******* WARNING ******"
	glog.Errorf("\n%s\n%s\n%s", warningMsg, msg, warningMsg)
}

func PrintServerInsecureWarning() {
	PrintInsecureWarning("Serving plaintext HTTP; this is dangerous!")
}

func PrintClientInsecureWarning() {
	PrintInsecureWarning("Connecting to backend over plaintext HTTP; this is dangerous!")
}
