package main

import (
	"os"
	"time"

	"github.com/flaccid/j2xrp/proxy"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "j2xrp"
	app.Version = "0.0.1"
	app.Usage = "A reverse proxy that converts a JSON request to an XML request"
	app.Compiled = time.Now()
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Chris Fordham",
			Email: "chris@fordham-nagy.id.au",
		},
	}
	app.Copyright = "(c) 2017 Chris Fordham"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "port, p",
			Value:  "9090",
			Usage:  "listen port",
			EnvVar: "PORT",
		},
		cli.StringFlag{
			Name:  "scheme, s",
			Value: "http",
			Usage: "protocol scheme, http or https",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.NewExitError("usage: j2xrp [command options] [host:port]", 1)
		}
		name := c.Args()[0]
		proxy.Serve(c.String("scheme"), name, c.String("port"))
		return nil
	}

	app.Run(os.Args)
}
