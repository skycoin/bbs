package main

import (
	"github.com/skycoin/bbs/src/rpc"
	"github.com/skycoin/bbs/src/store/object"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

var (
	Address = "127.0.0.1:8996"
)

func main() {
	app := cli.NewApp()
	app.Name = "bbscli"
	app.Usage = "a command-line interface to interact with a Skycoin BBS node"
	app.Flags = cli.FlagsByName{
		cli.StringFlag{
			Name:        "address,a",
			Usage:       "rpc address of bbs node",
			EnvVar:      "BBS_ADDRESS",
			Value:       Address,
			Destination: &Address,
		},
	}
	app.Commands = cli.Commands{
		{
			Name:  "connections",
			Usage: "manages connections of the node",
			Subcommands: cli.Commands{
				{
					Name:  "list",
					Usage: "lists all connections",
					Action: func(ctx *cli.Context) error {
						return do(rpc.GetConnections())
					},
				},
				{
					Name:  "new",
					Usage: "adds a new connection",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "address, a",
							Usage: "address to add",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.NewConnection(&object.ConnectionIO{
							Address: ctx.String("address"),
						}))
					},
				},
				{
					Name:  "del",
					Usage: "removes a connection",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "address, a",
							Usage: "address to remove",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.DeleteConnection(&object.ConnectionIO{
							Address: ctx.String("address"),
						}))
					},
				},
			},
		},
		{
			Name:  "subscriptions",
			Usage: "manages subscriptions of the node",
			Subcommands: cli.Commands{
				{
					Name:  "list",
					Usage: "lists all subscriptions",
					Action: func(ctx *cli.Context) error {
						return do(rpc.GetSubscriptions())
					},
				},
				{
					Name:  "new",
					Usage: "adds a new subscription",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name: "public-key, pk",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.NewSubscription(&object.BoardIO{
							PubKeyStr: ctx.String("public-key"),
						}))
					},
				},
				{
					Name:  "delete",
					Usage: "removes a subscription",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name: "public-key, pk",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.DeleteSubscription(&object.BoardIO{
							PubKeyStr: ctx.String("public-key"),
						}))
					},
				},
			},
		},
		{
			Name:  "content",
			Usage: "manages boards and their content",
			Subcommands: cli.Commands{
				{
					Name:  "new_board",
					Usage: "creates a new board",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "name, n",
							Usage: "name of the board",
						},
						cli.StringFlag{
							Name:  "body, b",
							Usage: "body of the board",
						},
						cli.StringFlag{
							Name:  "seed, s",
							Usage: "seed to generate key pair of the board",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.NewBoard(&object.NewBoardIO{
							Name: ctx.String("name"),
							Body: ctx.String("body"),
							Seed: ctx.String("seed"),
						}))
					},
				},
				{
					Name:  "delete_board",
					Usage: "deletes a board",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "public-key, pk",
							Usage: "public key of the board to delete",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.DeleteBoard(&object.BoardIO{
							PubKeyStr: ctx.String("public-key"),
						}))
					},
				},
				{
					Name:  "export_board",
					Usage: "exports a board",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "public-key, pk",
							Usage: "public key of the board to export",
						},
						cli.StringFlag{
							Name:  "file-name, fn",
							Usage: "name of file to export board to",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.ExportBoard(&object.ExportBoardIO{
							PubKeyStr: ctx.String("public-key"),
							Name:      ctx.String("file-name"),
						}))
					},
				},
				{
					Name:  "import_board",
					Usage: "imports a board",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "public-key, pk",
							Usage: "public key of the board to import data to",
						},
						cli.StringFlag{
							Name:  "file-name, fn",
							Usage: "name of file to import board data from",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.ImportBoard(&object.ExportBoardIO{
							PubKeyStr: ctx.String("public-key"),
							Name:      ctx.String("file-name"),
						}))
					},
				},
				{
					Name:  "get_boards",
					Usage: "gets a list of hosted boards on the node",
					Action: func(ctx *cli.Context) error {
						return do(rpc.GetBoards())
					},
				},
				{
					Name:  "get_board",
					Usage: "gets a single board",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "board-public-key, bpk",
							Usage: "public key of the board to obtain",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.GetBoard(&object.BoardIO{
							PubKeyStr: ctx.String("board-public-key"),
						}))
					},
				},
				{
					Name:  "get_board_page",
					Usage: "gets a view of a board and it's threads",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "board-public-key, bpk",
							Usage: "public key of the board to obtain",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.GetBoardPage(&object.BoardIO{
							PubKeyStr: ctx.String("board-public-key"),
						}))
					},
				},
				{
					Name: "get_thread_page",
					Usage: "gets a view of a board's thread and it's posts",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name: "board-public-key, bpk",
							Usage: "the public key of the board in which the thread resides",
						},
						cli.StringFlag{
							Name: "thread-hash, th",
							Usage: "the hash of the thread in which to obtain thread page",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.GetThreadPage(&object.ThreadIO{
							BoardPubKeyStr: ctx.String("board-public-key"),
							ThreadRefStr: ctx.String("thread-hash"),
						}))
					},
				},
				{
					Name: "get_follow_page",
					Usage: "gets a view of users that the specified user is following/avoiding",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name: "board-public-key, bpk",
							Usage: "public key of board in which to obtain follow page",
						},
						cli.StringFlag{
							Name: "user-public-key, upk",
							Usage: "public key of user to get follow page of",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(rpc.GetFollowPage(&object.UserIO{
							BoardPubKeyStr: ctx.String("board-public-key"),
							UserPubKeyStr: ctx.String("user-public-key"),
						}))
					},
				},
			},
		},
	}
	if e := app.Run(os.Args); e != nil {
		log.Println(e)
	}
}

func do(method string, in interface{}) error {
	log.Println(rpc.Send(Address)(method, in))
	return nil
}
