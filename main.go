package main

import (
	"fmt"
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/extern"
	"github.com/evanlinjin/bbs/extern/gui"
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
	defer log.Println("Goodbye.")

	container, e := store.NewContainer(config)
	cmd.CatchError(e, "unable to create cxo container")
	defer container.Close()

	boardSaver, e := store.NewBoardSaver(config, container)
	cmd.CatchError(e, "unable to create board saver")

	userSaver, e := store.NewUserSaver(config, container)
	cmd.CatchError(e, "unable to create user saver")

	queueSaver, e := store.NewQueueSaver(config, container)
	cmd.CatchError(e, "unable to create queue saver")
	defer queueSaver.Close()

	var rpcServer *extern.RPCServer
	if config.Master() {
		rpcGateway := extern.NewRPCGateway(config, container, boardSaver, userSaver)
		rpcServer, e = extern.NewRPCServer(rpcGateway, config.RPCServerPort())
		cmd.CatchError(e, "unable to start rpc server")
		defer rpcServer.Close()
	}

	if config.WebGUIEnable() {
		host := fmt.Sprintf("%s:%d", LocalhostAddress, config.WebGUIPort())
		fullAddress := fmt.Sprintf("%s://%s", "http", host)

		gateway := extern.NewGateway(config, container, boardSaver, userSaver, queueSaver)
		e := gui.OpenWebInterface(host, gateway)
		cmd.CatchError(e, "unable to start web server")
		defer gui.Close()

		if config.WebGUIOpenBrowser() {
			go func() {
				time.Sleep(time.Millisecond * 100)
				log.Println("Opening web browser...")
				util.OpenBrowser(fullAddress)
			}()
		}
	}

	defer log.Println("Shutting down...")
	<-quit
}
