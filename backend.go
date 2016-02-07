package main

import (
	"crypto/cipher"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"github.com/jonathanwei/fproxy/fileserver"
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

	go updateFrontendPort(NewAEADOrDie(config.PortKey), l.Addr(), config.PortUpdateUrl)

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

func updateFrontendPort(aead cipher.AEAD, addr net.Addr, url string) {
	update := &pb.PortUpdate{
		Port: int32(addr.(*net.TCPAddr).Port),
	}

	url = url + "/" + EncryptProto(aead, update, nil)
	for {
		glog.Infof("Posting update: %v", update)
		err := getURL(url)
		if err != nil {
			// TODO: Change this to exponential backoff.
			glog.Warningf("Got error updating port: %v", err)
			time.Sleep(5 * time.Second)
		}
	}
}

func getURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("Got non-200 response: %v", resp)
	}

	_, err = io.Copy(ioutil.Discard, resp.Body)
	return err
}

func getBackendHTTPMux(config *pb.BackendConfig) http.Handler {
	mux := http.NewServeMux()

	fs := fileserver.Handler{Base: config.ServePath, ListDirs: true}
	mux.Handle("/", &fs)

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
