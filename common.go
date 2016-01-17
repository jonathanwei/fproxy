package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/codegangsta/cli"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

func defaultConfigPath(c *cli.Context, file string) string {
	path := filepath.Join(os.Getenv("HOME"), ".config/fproxy", file)
	if args := c.Args(); len(args) > 0 {
		path = args[0]
	}
	return path
}

// Returns the first element of args, or defaultArg if args is an empty slice.
func firstOrDefault(args []string, defaultArg string) string {
	if len(args) > 0 {
		return args[0]
	}
	return defaultArg
}

func readConfig(configPath string, msg proto.Message) {
	protoBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		glog.Fatalf("Could not read config at path %v: %v", configPath, err)
	}

	err = proto.UnmarshalText(string(protoBytes), msg)
	if err != nil {
		glog.Fatalf("Could not parse config at path %v: %v", configPath, err)
	}
}
