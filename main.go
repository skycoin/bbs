package main

import (
	"fmt"
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/extern"
	"github.com/evanlinjin/bbs/gui"
	"github.com/evanlinjin/bbs/store"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"time"
)

const LocalhostAddress = "127.0.0.1"

func main() {
	quit := cmd.CatchInterrupt()
	config := cmd.NewConfig().Parse()
	util.InitDataDir(config.ConfigDir())

	container, e := store.NewContainer(config)
	cmd.CatchError(e, "unable to create cxo container")

	boardSaver, e := store.NewBoardSaver(config, container)
	cmd.CatchError(e, "unable to create board saver")

	gateway := extern.NewGateway(config, container, boardSaver)

	if config.WebGUIEnable() {
		host := fmt.Sprintf("%s:%d", LocalhostAddress, config.WebGUIPort())
		fullAddress := fmt.Sprintf("%s://%s", "http", host)

		e := gui.LaunchWebInterface(host, gateway)
		cmd.CatchError(e, "unable to start web server")

		if config.WebGUIOpenBrowser() {
			go func() {
				time.Sleep(time.Millisecond * 100)
				log.Println("Opening web browser...")
				util.OpenBrowser(fullAddress)
			}()
		}
	}

	<-quit
	log.Println("Shutting down...")
	gui.Shutdown()
	container.Stop()
	log.Println("Goodbye.")
}
