package main

import (
	"os"

	"github.com/codegangsta/cli"
)

func main() {
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

	app.Run(os.Args)
}
