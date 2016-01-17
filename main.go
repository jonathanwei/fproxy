package main

import (
	"flag"
	"os"

	"github.com/codegangsta/cli"

	// Redirect grpc logs to glog.
	_ "google.golang.org/grpc/grpclog/glogger"
)

func main() {
	flag.Parse()

	app := cli.NewApp()
	app.Name = "fproxy"
	app.Usage = "A file proxy server with SSH and web access."
	app.Authors = []cli.Author{
		{Name: "Jonathan Wei,"},
		{Name: "Sanjay Menakuru"},
	}
	app.Version = "0.1"
	app.Commands = []cli.Command{
		feCommand,
		beCommand,
	}

	var flags []string
	flags = append(flags, os.Args[0])
	flags = append(flags, flag.Args()...)

	app.Run(flags)
}
