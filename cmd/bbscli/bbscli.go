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
			Usage:       "http port of bbsnode where json api is served",
			Destination: &Port,
			Value:       Port,
		},
		cli.IntFlag{
			Name:        "timeout, t",
			EnvVar:      "CLI_TIMEOUT",
			Usage:       "operation timeout in seconds, negative to disable",
			Destination: &Timeout,
			Value:       Timeout,
		},
		cli.BoolFlag{
			Name:        "indent, i",
			EnvVar:      "CLI_INDENT",
			Usage:       "indent json output",
			Destination: &Indent,
		},
	}

	Board     = ""
	BoardFlag = cli.StringFlag{
		Name:        "board, b",
		Usage:       "reference/public key of board",
		Destination: &Board,
	}
	BoardName     = ""
	BoardNameFlag = cli.StringFlag{
		Name:        "board_name, bn",
		Usage:       "name of board",
		Destination: &BoardName,
	}
	BoardDescription     = ""
	BoardDescriptionFlag = cli.StringFlag{
		Name:        "board_description, bd",
		Usage:       "description of board",
		Destination: &BoardDescription,
	}
	BoardSubmissionAddresses     = ""
	BoardSubmissionAddressesFlag = cli.StringFlag{
		Name:        "board_submission_addresses, bsa",
		Usage:       "submission addresses of the board, separated with commas ','",
		Destination: &BoardSubmissionAddresses,
	}

	Thread     = ""
	ThreadFlag = cli.StringFlag{
		Name:        "thread, t",
		Usage:       "reference/public key of thread",
		Destination: &Thread,
	}
	ThreadName     = ""
	ThreadNameFlag = cli.StringFlag{
		Name:        "thread_name, tn",
		Usage:       "name of thread",
		Destination: &ThreadName,
	}
	ThreadDescription     = ""
	ThreadDescriptionFlag = cli.StringFlag{
		Name:        "thread_description, td",
		Usage:       "description of thread",
		Destination: &ThreadDescription,
	}

	Post     = ""
	PostFlag = cli.StringFlag{
		Name:        "post, p",
		Usage:       "reference/public key of post",
		Destination: &Post,
	}
	PostTitle     = ""
	PostTitleFlag = cli.StringFlag{
		Name:        "post_title, pt",
		Usage:       "title of post",
		Destination: &PostTitle,
	}
	PostBody     = ""
	PostBodyFlag = cli.StringFlag{
		Name:        "post_body, pb",
		Usage:       "body of post",
		Destination: &PostBody,
	}

	FromBoard     = ""
	FromBoardFlag = cli.StringFlag{
		Name:        "from_board, fb",
		Usage:       "reference/public key of board to import thread from",
		Destination: &FromBoard,
	}
	ToBoard     = ""
	ToBoardFlag = cli.StringFlag{
		Name:        "to_board, tb",
		Usage:       "reference/public key of board to import thread to",
		Destination: &ToBoard,
	}

	Seed     = ""
	SeedFlag = cli.StringFlag{
		Name:        "seed, s",
		Usage:       "random collection of charactors used to generate public/private key pair",
		Destination: &Seed,
	}
)

func main() {
	app := cli.NewApp()
	app.Usage = "command line interface for configuring bbsnode"
	app.EnableBashCompletion = true
	app.Flags = GlobalFlags
	app.Commands = []cli.Command{
		{
			Name:   "get_boards",
			Usage:  "lists boards in which the node is subscribed to",
			Action: try(gui.GetBoards()),
		},
		{
			Name:   "new_board",
			Usage:  "creates a new board that is hosted on the node",
			Flags:  []cli.Flag{BoardNameFlag, BoardDescriptionFlag, BoardSubmissionAddressesFlag, SeedFlag},
			Action: try(gui.NewBoard(&BoardName, &BoardDescription, &BoardSubmissionAddresses, &Seed)),
		},
		{
			Name:   "remove_board",
			Usage:  "removes a board that is hosted on the node",
			Flags:  []cli.Flag{BoardFlag},
			Action: try(gui.RemoveBoard(&Board)),
		},
		{
			Name:   "get_boardpage",
			Usage:  "obtains the board page - a page of board information and threads",
			Flags:  []cli.Flag{BoardFlag},
			Action: try(gui.GetBoardPage(&Board)),
		},
		{
			Name:   "get_threads",
			Usage:  "obtains threads of a specified board of public key",
			Flags:  []cli.Flag{BoardFlag},
			Action: try(gui.GetThreads(&Board)),
		},
		{
			Name:   "new_thread",
			Usage:  "creates a new thread on specified board of public key",
			Flags:  []cli.Flag{BoardFlag, ThreadNameFlag, ThreadDescriptionFlag},
			Action: try(gui.NewThread(&Board, &ThreadName, &ThreadDescription)),
		},
		{
			Name:   "remove_thread",
			Usage:  "removes a thread from specified board of public key",
			Flags:  []cli.Flag{BoardFlag, ThreadFlag},
			Action: try(gui.RemoveThread(&Board, &Thread)),
		},
		{
			Name:   "get_threadpage",
			Usage:  "obtains the thread page of specified board and thread",
			Flags:  []cli.Flag{BoardFlag, ThreadFlag},
			Action: try(gui.GetThreadPage(&Board, &Thread)),
		},
		{
			Name:   "get_posts",
			Usage:  "obtains the posts of a thread of specified board and thread",
			Flags:  []cli.Flag{BoardFlag, ThreadFlag},
			Action: try(gui.GetPosts(&Board, &Thread)),
		},
		{
			Name:   "new_post",
			Usage:  "creates a new post on specified board and thread",
			Flags:  []cli.Flag{BoardFlag, ThreadFlag, PostTitleFlag, PostBodyFlag},
			Action: try(gui.NewPost(&Board, &Thread, &PostTitle, &PostBody)),
		},
		{
			Name:   "remove_post",
			Usage:  "removes a post from specified board and thread",
			Flags:  []cli.Flag{BoardFlag, ThreadFlag, PostFlag},
			Action: try(gui.RemovePost(&Board, &Thread, &Post)),
		},
		{
			Name:   "import_thread",
			Usage:  "imports a thread from one board to another",
			Flags:  []cli.Flag{FromBoardFlag, ThreadFlag, ToBoardFlag},
			Action: try(gui.ImportThread(&FromBoard, &Thread, &ToBoard)),
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
			return nil
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
