package main

import (
	"github.com/skycoin/bbs/cmd/bbsnode/args"
	"github.com/skycoin/bbs/src/dev"
	"github.com/skycoin/bbs/src/gui"
	"github.com/skycoin/bbs/src/rpc"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/msg"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	quit := CatchInterrupt()
	config, e := args.NewConfig().Parse().PostProcess()
	if e != nil {
		panic(e)
	}
	util.InitDataDir(config.ConfigDir())
	log.Println("[CONFIG] Master mode:", config.Master())
	defer log.Println("Goodbye.")

	log.Println("[CONFIG] Connecting to cxo on port", config.CXOPort())
	container, e := cxo.NewContainer(config)
	CatchError(e, "unable to create cxo container")
	defer container.Close()

	boardSaver, e := store.NewBoardSaver(config, container)
	CatchError(e, "unable to create board saver")
	defer boardSaver.Close()

	userSaver, e := store.NewUserSaver(config, container)
	CatchError(e, "unable to create user saver")

	queueSaver, e := msg.NewQueueSaver(config, container)
	CatchError(e, "unable to create queue saver")
	defer queueSaver.Close()

	var rpcServer *rpc.Server
	if config.Master() {
		rpcGateway := rpc.NewGateway(
			config, container, boardSaver, userSaver)

		rpcServer, e = rpc.NewServer(rpcGateway, config.RPCPort())
		CatchError(e, "unable to start rpc server")
		defer rpcServer.Close()

		log.Println("[RPCSERVER] Serving on address:", rpcServer.Address())
	}

	gateway := gui.NewGateway(
		config, container, boardSaver, userSaver, queueSaver, quit)

	serveAddr, e := gui.OpenWebInterface(gateway)
	CatchError(e, "unable to start web server")
	defer gui.Close()

	if config.TestMode() {
		tester, e := dev.NewTester(config, gateway)
		CatchError(e, "unable to start tester")
		defer tester.Close()
	}

	log.Println("[WEBGUI] Serving on:", serveAddr)

	if config.WebGUIEnable() && config.WebGUIOpenBrowser() {
		go func() {
			time.Sleep(time.Millisecond * 100)
			log.Println("Opening web browser...")
			util.OpenBrowser(serveAddr)
		}()
	}

	log.Println("!!! EVERYTHING UP AND RUNNING !!!")
	defer log.Println("Shutting down...")
	<-quit
	time.Sleep(time.Second)
}

// CatchInterrupt catches Ctrl+C behaviour.
func CatchInterrupt() chan int {
	quit := make(chan int)
	go func(q chan<- int) {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		signal.Stop(sigchan)
		q <- 1
	}(quit)
	return quit
}

// CatchError catches an error and panics.
func CatchError(e error, msg string, args ...interface{}) {
	if e != nil {
		log.Panicf(msg+": %v", append(args, e)...)
	}
}
