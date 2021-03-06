package main

import (
	"strings"
	"sync"

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

TODO: rewrite.

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

oauth_config: this field lets you configure the oauth settings for
authentication. It requires the client ID, client secret, and redirect URL to
be set. Currently, only Google oauth is supported. Note that the redirect URL
must always use the path '/oauth2Callback'. A sample clause:
  oauth_config {
    client_id: "thisisanid.apps.googleusercontent.com"
    client_secret: "thisis-a-clientsecret"
    redirect_url: "http://localhost:8000/oauth2Callback"
  }

auth_cookie_key: this field is used as a key for encrypting and authenticating
auth cookies. It must be a 16 byte string. A sample clause:
  auth_cookie_key: "1234567890123456"

email_to_user_id: this repeated field is used as a mapping between a verified
email to the user ID for use in ACLs. A sample clause:
  email_to_user_id {
    key: "email@example.com",
    value: "user1",
  }

backend: this field is used to store all configuration for the backend server
we connect to. This includes the address it's located at, the root CA to use
for it, the client certificate and key with which to identify the frontend
server, and the server name for the backend. A sample clause:

backend {
  addr: ":10000"
  root_ca_file: "/certs/backend/ca.pem"
  client_cert_file: "/certs/frontend/frontend.pem"
  client_key_file: "/certs/frontend/frontend-key.pem"
  server_name: "backend.fproxy"
}
`),
	Action: func(c *cli.Context) {
		runFe(defaultConfigPath(c, "frontend.textproto"))
	},
}

func runFe(configPath string) {
	var config pb.FrontendConfig
	readConfig(configPath, &config)

	if len(config.EmailToUserId) == 0 {
		glog.Error("No email to user mappings were given. The HTTP server will reject all users.")
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		runTCPProxy(config.TcpProxyRoute)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runHttpServer(&config)
	}()

	wg.Wait()
}
