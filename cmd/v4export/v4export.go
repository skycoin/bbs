package main

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/tag"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"encoding/json"
)

const (
	Version = "5.0"
)

var (
	FlagPort        = 7410
	FlagBoardPubKey = ""
	FlagFileName    = "exported_board"
)

func main() {
	app := cli.NewApp()
	app.Name = "v4export"
	app.Version = Version
	app.Usage = "used for exporting boards from bbsnode v4.*"
	app.Flags = cli.FlagsByName{
		cli.IntFlag{
			Name:        "http-port, p",
			Usage:       "http port of bbsnode v4.*",
			Value:       FlagPort,
			Destination: &FlagPort,
		},
		cli.StringFlag{
			Name:        "board-public-key, bpk",
			Usage:       "public key of the board to export",
			Value:       FlagBoardPubKey,
			Destination: &FlagBoardPubKey,
		},
		cli.StringFlag{
			Name:        "file-name, fn",
			Usage:       "file name to export the board to",
			Value:       FlagFileName,
			Destination: &FlagFileName,
		},
	}
	app.Action = action
	if e := app.Run(os.Args); e != nil {
		log.Println(e)
		log.Println("run 'v4export -h' to see usage")
	}
}

func action(ctx *cli.Context) error {

	// Check inputs.
	if e := tag.CheckPort(FlagPort); e != nil {
		return boo.Wrap(e, "invalid 'http-port' provided")
	}
	if _, e := tag.GetPubKey(FlagBoardPubKey); e != nil {
		return boo.Wrap(e, "invalid 'board-public-key' provided")
	}
	if e := tag.CheckPath(FlagFileName); e != nil {
		return boo.Wrap(e, "invalid 'file-name' provided")
	}

	v := make(url.Values)
	v.Add("board_public_key", FlagBoardPubKey)
	v.Add("file_name", FlagFileName)

	address := "http://127.0.0.1:" + strconv.Itoa(FlagPort) + "/api/admin/board/export"

	res, e := http.PostForm(address, v)
	if e != nil {
		return boo.Wrap(e, "failed to connect to v4 bbsnode")
	}

	out, e := ioutil.ReadAll(res.Body)
	if e != nil {
		return boo.Wrap(e, "unable to read reply")
	}

	outRep := make(map[string]interface{})
	if e := json.Unmarshal(out, &outRep); e != nil {
		return boo.Wrap(e, "unable to parse reply")
	}

	outRepRep, _ := json.MarshalIndent(outRep, "", "    ")

	log.Println(string(outRepRep))
	return nil
}
