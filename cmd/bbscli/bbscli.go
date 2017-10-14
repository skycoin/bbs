package main

import (
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/bbs/src/rpc"
	"github.com/skycoin/bbs/src/store/object"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
)

var (
	Address = "[::]:8996"
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
			Name:  "tools",
			Usage: "cryptography tools",
			Subcommands: cli.Commands{
				{
					Name:  "generate_seed",
					Usage: "generates a random unique seed",
					Action: func(ctx *cli.Context) error {
						return do(keys.GenerateSeed())
					},
				},
				{
					Name:  "generate_key_pair",
					Usage: "generates a public, private key pair",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "seed, s",
							Usage: "seed to generate key pair with",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(keys.GenerateKeyPair(&keys.GenerateKeyPairIn{
							Seed: ctx.String("seed"),
						}))
					},
				},
				{
					Name:  "sum_sha256",
					Usage: "finds the SHA256 hash sum of given data",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "data, d",
							Usage: "data to hash",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(keys.SumSHA256(&keys.SumSHA256In{
							Data: ctx.String("data"),
						}))
					},
				},
				{
					Name:  "sign_hash",
					Usage: "generates a signature of a hash with given secret key",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "hash, h",
							Usage: "hash to be signed",
						},
						cli.StringFlag{
							Name:  "secret-key, sk",
							Usage: "secret key to sign hash with",
						},
					},
					Action: func(ctx *cli.Context) error {
						return do(keys.SignHash(&keys.SignHashIn{
							Hash:   ctx.String("hash"),
							SecKey: ctx.String("secret_key"),
						}))
					},
				},
			},
		},
		{
			Name:  "connections",
			Usage: "manages connections of the node",
			Subcommands: cli.Commands{
				{
					Name:  "list",
					Usage: "lists all connections",
					Action: func(ctx *cli.Context) error {
						return call(rpc.GetConnections())
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
						return call(rpc.NewConnection(&object.ConnectionIO{
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
						return call(rpc.DeleteConnection(&object.ConnectionIO{
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
						return call(rpc.GetSubscriptions())
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
						return call(rpc.NewSubscription(&object.BoardIO{
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
						return call(rpc.DeleteSubscription(&object.BoardIO{
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
						return call(rpc.NewBoard(&object.NewBoardIO{
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
						return call(rpc.DeleteBoard(&object.BoardIO{
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
						return call(rpc.ExportBoard(&object.ExportBoardIO{
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
						return call(rpc.ImportBoard(&object.ExportBoardIO{
							PubKeyStr: ctx.String("public-key"),
							Name:      ctx.String("file-name"),
						}))
					},
				},
				{
					Name:  "get_boards",
					Usage: "gets a list of hosted boards on the node",
					Action: func(ctx *cli.Context) error {
						return call(rpc.GetBoards())
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
						return call(rpc.GetBoard(&object.BoardIO{
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
						return call(rpc.GetBoardPage(&object.BoardIO{
							PubKeyStr: ctx.String("board-public-key"),
						}))
					},
				},
				{
					Name:  "get_thread_page",
					Usage: "gets a view of a board's thread and it's posts",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "board-public-key, bpk",
							Usage: "the public key of the board in which the thread resides",
						},
						cli.StringFlag{
							Name:  "thread-hash, th",
							Usage: "the hash of the thread in which to obtain thread page",
						},
					},
					Action: func(ctx *cli.Context) error {
						return call(rpc.GetThreadPage(&object.ThreadIO{
							BoardPubKeyStr: ctx.String("board-public-key"),
							ThreadRefStr:   ctx.String("thread-hash"),
						}))
					},
				},
				{
					Name:  "get_follow_page",
					Usage: "gets a view of users that the specified user is following/avoiding",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "board-public-key, bpk",
							Usage: "public key of board in which to obtain follow page",
						},
						cli.StringFlag{
							Name:  "user-public-key, upk",
							Usage: "public key of user to get follow page of",
						},
					},
					Action: func(ctx *cli.Context) error {
						return call(rpc.GetFollowPage(&object.UserIO{
							BoardPubKeyStr: ctx.String("board-public-key"),
							UserPubKeyStr:  ctx.String("user-public-key"),
						}))
					},
				},
				{
					Name:  "new_thread",
					Usage: "submits a new thread to specified board",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "board-public-key, bpk",
							Usage: "public key of the board in which to submit the thread",
						},
						cli.StringFlag{
							Name:  "name, n",
							Usage: "name of the thread",
						},
						cli.StringFlag{
							Name:  "body, b",
							Usage: "body of the thread",
						},
						cli.StringFlag{
							Name:  "user-public-key, upk",
							Usage: "public key of the thread's creator",
						},
						cli.StringFlag{
							Name:  "user-secret-key, usk",
							Usage: "secret key of the thread's creator",
						},
					},
					Action: func(ctx *cli.Context) error {
						return call(rpc.NewThread(&object.NewThreadIO{
							BoardPubKeyStr: ctx.String("board-public-key"),
							Name:           ctx.String("name"),
							Body:           ctx.String("body"),
							UserPubKeyStr:  ctx.String("user-public-key"),
							UserSecKeyStr:  ctx.String("user-secret-key"),
						}))
					},
				},
				{
					Name:  "new_post",
					Usage: "submits a new post to specified board and thread",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "board-public-key, bpk",
							Usage: "public key of board in which to submit the post",
						},
						cli.StringFlag{
							Name:  "thread-hash, th",
							Usage: "hash of the thread in which to submit the post",
						},
						cli.StringFlag{
							Name:  "post-hash, ph",
							Usage: "(optional) hash of post in which this post is a reply to",
						},
						cli.StringFlag{
							Name:  "name, n",
							Usage: "name of the post",
						},
						cli.StringFlag{
							Name:  "body, b",
							Usage: "body of the post",
						},
						cli.StringFlag{
							Name:  "user-public-key, upk",
							Usage: "public key of the post's creator",
						},
						cli.StringFlag{
							Name:  "user-secret-key, usk",
							Usage: "secret key of the post's creator",
						},
					},
					Action: func(ctx *cli.Context) error {
						// TODO: Have images too.
						return call(rpc.NewPost(&object.NewPostIO{
							BoardPubKeyStr: ctx.String("board-public-key"),
							ThreadRefStr:   ctx.String("thread-hash"),
							PostRefStr:     ctx.String("post-hash"),
							Name:           ctx.String("name"),
							Body:           ctx.String("body"),
							UserPubKeyStr:  ctx.String("user-public-key"),
							UserSecKeyStr:  ctx.String("user-secret-key"),
						}))
					},
				},
				{
					Name:  "vote_thread",
					Usage: "submits a vote for a given thread",
					Flags: cli.FlagsByName{
						cli.StringFlag{
							Name:  "board-public-key, bpk",
							Usage: "public key of board in which to submit the vote",
						},
						cli.StringFlag{
							Name:  "thread-hash, th",
							Usage: "hash of the thread to vote",
						},
						cli.StringFlag{
							Name:  "value, v",
							Usage: "value of the vote (+1, 0, -1)",
						},
						cli.StringFlag{
							Name:  "tag, t",
							Usage: "the vote's tag",
						},
						cli.StringFlag{
							Name:  "user-public-key, upk",
							Usage: "public key of the vote's creator",
						},
						cli.StringFlag{
							Name:  "user-secret-key, usk",
							Usage: "secret key of the vote's creator",
						},
					},
					Action: func(ctx *cli.Context) error {
						return call(rpc.VoteThread(&object.ThreadVoteIO{
							BoardPubKeyStr: ctx.String("board-public-key"),
							ThreadRefStr:   ctx.String("thread-hash"),
							ModeStr:        ctx.String("value"),
							TagStr:         ctx.String("tag"),
							UserPubKeyStr:  ctx.String("user-public-key"),
							UserSecKeyStr:  ctx.String("user-secret-key"),
						}))
					},
				},
				//{
				//	Name: "vote_post",
				//	Usage: "submits a vote for a given post",
				//	Flags: cli.FlagsByName{
				//
				//	},
				//},
			},
		},
	}
	if e := app.Run(os.Args); e != nil {
		log.Println(e)
	}
}

func call(method string, in interface{}) error {
	log.Println(rpc.Send(Address)(method, in))
	return nil
}

func do(out interface{}, e error) error {
	log.Println(rpc.Do(out, e))
	return nil
}
