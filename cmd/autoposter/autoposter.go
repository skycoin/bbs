package main

import "gopkg.in/urfave/cli.v1"

var (
	HTTPPort = 7410
)

func main() {
	app := cli.NewApp()
	app.Name = "Skycoin BBS AutoPoster"
	app.Usage = "Used for testing Skycoin BBS"
	app.Flags = cli.FlagsByName{
		cli.IntFlag{
			Name: "http-port,p",
			Destination: &HTTPPort,
			Value: HTTPPort,
		},
	}
	app.Commands = cli.Commands{
		{

		},
	}
	app.Action = cli.ActionFunc(func(ctx *cli.Context) error {
		return nil
	})
}
