package main

import (
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
)

var feCommand = cli.Command{
	Name:    "frontend",
	Aliases: []string{"fe"},
	Usage:   "Run the frontend.",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "tcp-proxy-config",
			Value: filepath.Join(os.Getenv("HOME"), ".config/fproxy/tcp.json"),
			Usage: "Configuration for TCP proxying in JSON format.",
		},
	},
	Action: func(c *cli.Context) {
		runFe(c.String("tcp-proxy-config"))
	},
}

func runFe(configPath string) {
	runTCPProxy(configPath)
}
