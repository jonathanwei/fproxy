package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/golang/protobuf/proto"
	pb "github.com/jonathanwei/fproxy/proto"
)

var feCommand = cli.Command{
	Name:      "frontend",
	Aliases:   []string{"fe"},
	Usage:     "Run the frontend.",
	ArgsUsage: "[path/to/config/proto]",
	Description: strings.TrimSpace(`
The default path for the config proto is "~/.config/fproxy/frontend.textproto".
`),
	Action: func(c *cli.Context) {
		runFe(readFeConfig(c.Args()))
	},
}

func runFe(config *pb.FrontendConfig) {
	runTCPProxy(config.TcpProxyRoute)
}

func readFeConfig(args []string) *pb.FrontendConfig {
	configPath := filepath.Join(os.Getenv("HOME"), ".config/fproxy/frontend.textproto")
	if len(args) > 0 {
		configPath = args[0]
	}

	protoBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Could not read frontend config: %v", err)
	}

	var config pb.FrontendConfig

	err = proto.UnmarshalText(string(protoBytes), &config)
	if err != nil {
		log.Fatalf("Could not parse frontend config: %v", err)
	}

	return &config
}
