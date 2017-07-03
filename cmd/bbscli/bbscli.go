package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/skycoin/bbs/src/gui"
	"gopkg.in/urfave/cli.v1"
	"os"
	"time"
)

var (
	Port        = 7410
	Timeout     = 10
	Indent      = true
	GlobalFlags = []cli.Flag{
		cli.IntFlag{
			Name:        "port, p",
			EnvVar:      "CLI_PORT",
			Usage:       "Http port of bbsnode where json api is served",
			Destination: &Port,
			Value:       Port,
		},
		cli.IntFlag{
			Name:        "timeout, t",
			EnvVar:      "CLI_TIMEOUT",
			Usage:       "Operation timeout in seconds, negative to disable",
			Destination: &Timeout,
			Value:       Timeout,
		},
		cli.BoolFlag{
			Name:        "indent, i",
			EnvVar:      "CLI_INDENT",
			Usage:       "Indent json output",
			Destination: &Indent,
		},
	}

	Board     = ""
	BoardFlag = cli.StringFlag{
		Name:        "board, b",
		Usage:       "Reference/public key of board",
		Destination: &Board,
	}
	BoardName     = ""
	BoardNameFlag = cli.StringFlag{
		Name:        "board_name, bn",
		Usage:       "Name of board",
		Destination: &BoardName,
	}
	BoardDescription     = ""
	BoardDescriptionFlag = cli.StringFlag{
		Name:        "board_description, bd",
		Usage:       "Description of board",
		Destination: &BoardDescription,
	}
	BoardSubmissionAddresses     = ""
	BoardSubmissionAddressesFlag = cli.StringFlag{
		Name:        "board_submission_addresses, bsa",
		Usage:       "Submission addresses of the board, separated with commas ','",
		Destination: &BoardSubmissionAddresses,
	}

	Thread     = ""
	ThreadFlag = cli.StringFlag{
		Name:        "thread, t",
		Usage:       "Reference/public key of thread",
		Destination: &Thread,
	}
	ThreadName     = ""
	ThreadNameFlag = cli.StringFlag{
		Name:        "thread_name, tn",
		Usage:       "Name of thread",
		Destination: &ThreadName,
	}
	ThreadDescription     = ""
	ThreadDescriptionFlag = cli.StringFlag{
		Name:        "thread_description, td",
		Usage:       "Description of thread",
		Destination: &ThreadDescription,
	}

	Post     = ""
	PostFlag = cli.StringFlag{
		Name:        "post, p",
		Usage:       "Reference/public key of post",
		Destination: &Post,
	}
	PostTitle     = ""
	PostTitleFlag = cli.StringFlag{
		Name:        "post_title, pt",
		Usage:       "Title of post",
		Destination: &PostTitle,
	}
	PostBody     = ""
	PostBodyFlag = cli.StringFlag{
		Name:        "post_body, pb",
		Usage:       "Body of post",
		Destination: &PostBody,
	}

	VoteMode     = ""
	VoteModeFlag = cli.StringFlag{
		Name:        "vote_mode, vm",
		Usage:       "Mode of vote, valid values are -1, 0, +1",
		Destination: &VoteMode,
	}
	VoteTag     = ""
	VoteTagFlag = cli.StringFlag{
		Name:        "vote_tag, vt",
		Usage:       "Additional tag information of vote",
		Destination: &VoteTag,
	}

	User     = ""
	UserFlag = cli.StringFlag{
		Name:        "user, u",
		Usage:       "Reference/public key of user",
		Destination: &User,
	}
	UserAlias     = ""
	UserAliasFlag = cli.StringFlag{
		Name:        "user_alias, ua",
		Usage:       "User alias of user",
		Destination: &UserAlias,
	}

	FromBoard     = ""
	FromBoardFlag = cli.StringFlag{
		Name:        "from_board, fb",
		Usage:       "Reference/public key of board to import thread from",
		Destination: &FromBoard,
	}
	ToBoard     = ""
	ToBoardFlag = cli.StringFlag{
		Name:        "to_board, tb",
		Usage:       "Reference/public key of board to import thread to",
		Destination: &ToBoard,
	}

	ConnectionAddress     = ""
	ConnectionAddressFlag = cli.StringFlag{
		Name:        "connection_address, ca",
		Usage:       "Connection address of external node",
		Destination: &ConnectionAddress,
	}
	SubmissionAddress     = ""
	SubmissionAddressFlag = cli.StringFlag{
		Name:        "submission_address, sa",
		Usage:       "Submission address of a board",
		Destination: &SubmissionAddress,
	}

	Seed     = ""
	SeedFlag = cli.StringFlag{
		Name:        "seed, s",
		Usage:       "Random collection of characters used to generate public/private key pair",
		Destination: &Seed,
	}

	ThreadCount     = ""
	ThreadCountFlag = cli.StringFlag{
		Name:        "thread_count, tc",
		Usage:       "Number of threads to generate",
		Destination: &ThreadCount,
	}
	PostCountMin     = ""
	PostCountMinFlag = cli.StringFlag{
		Name:        "post_count_min, pc_min",
		Usage:       "Minimum number of posts to generate for a thread",
		Destination: &PostCountMin,
	}
	PostCountMax     = ""
	PostCountMaxFlag = cli.StringFlag{
		Name:        "post_count_max, pc_max",
		Usage:       "Maximum number of posts to generate for a thread",
		Destination: &PostCountMax,
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "Command line interface for configuring bbsnode"
	app.EnableBashCompletion = true
	app.Flags = GlobalFlags
	app.Commands = []cli.Command{
		{
			Name:   "quit",
			Usage:  "Quits the node",
			Action: try(gui.Quit()),
		},
		{
			Name:  "stats",
			Usage: "Actions regarding node stats",
			Subcommands: []cli.Command{
				{
					Name:   "get",
					Usage:  "Obtains node stats",
					Action: try(gui.StatsGet()),
				},
			},
		},
		{
			Name:  "connections",
			Usage: "Actions regarding node connections",
			Subcommands: []cli.Command{
				{
					Name:   "get_all",
					Usage:  "Gets all addresses of nodes this node is connected to",
					Action: try(gui.ConnectionsGetAll()),
				},
				{
					Name:   "add",
					Usage:  "Connects node to provided address",
					Flags:  []cli.Flag{ConnectionAddressFlag},
					Action: try(gui.ConnectionsAdd(&ConnectionAddress)),
				},
				{
					Name:   "remove",
					Usage:  "Removes connection to node of provided address",
					Flags:  []cli.Flag{ConnectionAddressFlag},
					Action: try(gui.ConnectionsRemove(&ConnectionAddress)),
				},
			},
		},
		{
			Name:  "subscriptions",
			Usage: "Actions regarding board subscriptions of this node",
			Subcommands: []cli.Command{
				{
					Name:   "get_all",
					Usage:  "Gets a list of all board subscriptions",
					Action: try(gui.SubscriptionsGetAll()),
				},
				{
					Name:   "get",
					Usage:  "Gets subscription information of a provided board",
					Flags:  []cli.Flag{BoardFlag},
					Action: try(gui.SubscriptionsGet(&Board)),
				},
				{
					Name:   "add",
					Usage:  "Adds a board subscription",
					Flags:  []cli.Flag{BoardFlag, ConnectionAddressFlag},
					Action: try(gui.SubscriptionsAdd(&Board, &ConnectionAddress)),
				},
				{
					Name:   "remove",
					Usage:  "Removes a board subscription",
					Flags:  []cli.Flag{BoardFlag},
					Action: try(gui.SubscriptionsRemove(&Board)),
				},
			},
		},
		{
			Name:  "users",
			Usage: "Actions regarding user management of this node",
			Subcommands: []cli.Command{
				{
					Name:   "get_all",
					Usage:  "Gets all saved users on the node",
					Action: try(gui.UsersGetAll()),
				},
				{
					Name:   "add",
					Usage:  "Saves a non-master user to the node",
					Flags:  []cli.Flag{UserFlag, UserAliasFlag},
					Action: try(gui.UsersAdd(&User, &UserAlias)),
				},
				{
					Name:   "remove",
					Usage:  "Removes a user from the node",
					Flags:  []cli.Flag{UserFlag},
					Action: try(gui.UsersRemove(&User)),
				},
				{
					Name:  "masters",
					Usage: "Actions regarding master users saved on this node",
					Subcommands: []cli.Command{
						{
							Name:   "get_all",
							Usage:  "Gets all master users saved on this node",
							Action: try(gui.UsersMastersGetAll()),
						},
						{
							Name:   "add",
							Usage:  "Saves a new master user to this node",
							Flags:  []cli.Flag{UserAliasFlag, SeedFlag},
							Action: try(gui.UsersMastersAdd(&UserAlias, &Seed)),
						},
						{
							Name:  "current",
							Usage: "Actions regarding the currently active master user",
							Subcommands: []cli.Command{
								{
									Name:   "get",
									Usage:  "Gets the currently active master user",
									Action: try(gui.UsersMastersCurrentGet()),
								},
								{
									Name:   "set",
									Usage:  "Sets the currently active master user",
									Flags:  []cli.Flag{UserFlag},
									Action: try(gui.UsersMastersCurrentSet(&User)),
								},
							},
						},
					},
				},
			},
		},
		{
			Name:  "boards",
			Usage: "Actions regarding boards",
			Subcommands: []cli.Command{
				{
					Name:   "get_all",
					Usage:  "Gets all boards",
					Action: try(gui.BoardsGetAll()),
				},
				{
					Name:   "get",
					Usage:  "Gets a specified board",
					Flags:  []cli.Flag{BoardFlag},
					Action: try(gui.BoardsGet(&Board)),
				},
				{
					Name:   "add",
					Usage:  "Adds a board to host on this node",
					Flags:  []cli.Flag{BoardNameFlag, BoardDescriptionFlag, BoardSubmissionAddressesFlag, SeedFlag},
					Action: try(gui.BoardsAdd(&BoardName, &BoardDescription, &BoardSubmissionAddresses, &Seed)),
				},
				{
					Name:   "remove",
					Usage:  "Removes a specified board",
					Flags:  []cli.Flag{BoardFlag},
					Action: try(gui.BoardsRemove(&Board)),
				},
				{
					Name:  "meta",
					Usage: "Actions regarding board meta data",
					Subcommands: []cli.Command{
						{
							Name:   "get",
							Usage:  "Gets the entire board meta data object",
							Flags:  []cli.Flag{BoardFlag},
							Action: try(gui.BoardsMetaGet(&Board)),
						},
						{
							Name:  "submission_addresses",
							Usage: "Actions regarding the submission addresses of the board",
							Subcommands: []cli.Command{
								{
									Name:   "get_all",
									Usage:  "Gets all the submission addresses of the board",
									Flags:  []cli.Flag{BoardFlag},
									Action: try(gui.BoardsMetaSubmissionAddressesGetAll(&Board)),
								},
								{
									Name:   "add",
									Usage:  "Adds a submission address to the board",
									Flags:  []cli.Flag{BoardFlag, SubmissionAddressFlag},
									Action: try(gui.BoardsMetaSubmissionAddressesAdd(&Board, &SubmissionAddress)),
								},
								{
									Name:   "remove",
									Usage:  "Removes a submission address from the board",
									Flags:  []cli.Flag{BoardFlag, SubmissionAddressFlag},
									Action: try(gui.BoardsMetaSubmissionAddressesRemove(&Board, &SubmissionAddress)),
								},
							},
						},
					},
				},
				{
					Name:  "page",
					Usage: "Actions regarding the board represented as a page",
					Subcommands: []cli.Command{
						{
							Name:   "get",
							Usage:  "Gets the board represented as a page",
							Flags:  []cli.Flag{BoardFlag},
							Action: try(gui.BoardsPageGet(&Board)),
						},
					},
				},
			},
		},
		{
			Name:  "threads",
			Usage: "Actions regarding threads",
			Subcommands: []cli.Command{
				{
					Name:   "get_all",
					Usage:  "Gets all threads under the specified board",
					Flags:  []cli.Flag{BoardFlag},
					Action: try(gui.ThreadsGetAll(&Board)),
				},
				{
					Name:   "add",
					Usage:  "Adds a new thread to specified board",
					Flags:  []cli.Flag{BoardFlag, ThreadNameFlag, ThreadDescriptionFlag},
					Action: try(gui.ThreadsAdd(&Board, &ThreadName, &ThreadDescription)),
				},
				{
					Name:   "remove",
					Usage:  "Removes a thread from specified board",
					Flags:  []cli.Flag{BoardFlag, ThreadFlag},
					Action: try(gui.ThreadsRemove(&Board, &Thread)),
				},
				{
					Name:   "import",
					Usage:  "Imports a thread from one board to another",
					Flags:  []cli.Flag{FromBoardFlag, ThreadFlag, ToBoardFlag},
					Action: try(gui.ThreadsImport(&FromBoard, &Thread, &ToBoard)),
				},
				{
					Name:  "page",
					Usage: "Actions regarding the thread represented as a page",
					Subcommands: []cli.Command{
						{
							Name:   "get",
							Usage:  "Gets the thread represented as a page",
							Flags:  []cli.Flag{BoardFlag, ThreadFlag},
							Action: try(gui.ThreadsPageGet(&Board, &Thread)),
						},
					},
				},
				{
					Name:  "votes",
					Usage: "Actions regarding votes of a thread",
					Subcommands: []cli.Command{
						{
							Name:   "get",
							Usage:  "Gets the votes of the specified thread",
							Flags:  []cli.Flag{BoardFlag, ThreadFlag},
							Action: try(gui.ThreadsVotesGet(&Board, &Thread)),
						},
						{
							Name:   "add",
							Usage:  "Adds a vote to the specified thread",
							Flags:  []cli.Flag{BoardFlag, ThreadFlag, VoteModeFlag, VoteTagFlag},
							Action: try(gui.ThreadsVotesAdd(&Board, &Thread, &VoteMode, &VoteTag)),
						},
					},
				},
			},
		},
		{
			Name:  "posts",
			Usage: "Actions regarding posts",
			Subcommands: []cli.Command{
				{
					Name:   "get_all",
					Usage:  "Gets all posts under specified board and thread",
					Flags:  []cli.Flag{BoardFlag, ThreadFlag},
					Action: try(gui.PostsGetAll(&Board, &Thread)),
				},
				{
					Name:   "add",
					Usage:  "Adds a post to specified thread under specified board",
					Flags:  []cli.Flag{BoardFlag, ThreadFlag, PostTitleFlag, PostBodyFlag},
					Action: try(gui.PostsAdd(&Board, &Thread, &PostTitle, &PostBody)),
				},
				{
					Name:   "remove",
					Usage:  "Removes a specified post from a specified thread under a specified board",
					Flags:  []cli.Flag{BoardFlag, ThreadFlag, PostFlag},
					Action: try(gui.PostsRemove(&Board, &Thread, &Post)),
				},
				{
					Name:  "votes",
					Usage: "Actions regarding votes of a post",
					Subcommands: []cli.Command{
						{
							Name:   "get",
							Usage:  "Gets the votes of a specified post",
							Flags:  []cli.Flag{BoardFlag, PostFlag},
							Action: try(gui.PostsVotesGet(&Board, &Post)),
						},
						{
							Name:   "add",
							Usage:  "Adds a vote to the specified post",
							Flags:  []cli.Flag{BoardFlag, PostFlag, VoteModeFlag, VoteTagFlag},
							Action: try(gui.PostsVotesAdd(&Board, &Post, &VoteMode, &VoteTag)),
						},
					},
				},
			},
		},
		{
			Name:  "tests",
			Usage: "Actions regarding tests for the node",
			Subcommands: []cli.Command{
				{
					Name:   "add_filled_board",
					Usage:  "Adds a board, hosted on the node, filled with random threads and posts",
					Flags:  []cli.Flag{SeedFlag, ThreadCountFlag, PostCountMinFlag, PostCountMaxFlag},
					Action: try(gui.TestsAddFilledBoard(&Seed, &ThreadCount, &PostCountMin, &PostCountMax)),
				},
			},
		},
	}
	if e := app.Run(os.Args); e != nil {
		panic(e)
	}
}

func try(f func(ctx context.Context, port int) ([]byte, error)) cli.ActionFunc {
	ctx := context.Background()
	if Timeout >= 0 {
		ctx, _ = context.WithTimeout(
			ctx, time.Duration(Timeout)*time.Second)
	}
	return func(_ *cli.Context) error {
		data, e := f(ctx, Port)
		if e != nil {
			fmt.Println(e.Error())
			return e
		}
		if Indent {
			var prettyJSON bytes.Buffer
			json.Indent(&prettyJSON, data, "", "    ")
			data = prettyJSON.Bytes()
		}
		fmt.Println(string(data))
		return nil
	}
}
