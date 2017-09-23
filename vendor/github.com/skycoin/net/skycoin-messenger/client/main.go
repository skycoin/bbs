package main

import (
	"flag"
	"net/http"

	"fmt"
	"net"

	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/rpc"
	"github.com/skycoin/net/skycoin-messenger/websocket"
	"github.com/skycoin/skycoin/src/util/browser"
)

var (
	webDir           string
	rpcAddress       string
	webSocketAddress string
	openBrowser      bool
)

func parseFlags() {
	flag.StringVar(&webDir, "web-dir", "../web/dist", "directory of web files")
	flag.StringVar(&rpcAddress, "rpc-address", "localhost:8083", "rpc address to listen on")
	flag.StringVar(&webSocketAddress, "websocket-address", "localhost:8082", "websocket address to listen on")
	flag.BoolVar(&openBrowser, "open-browser", true, "whether to open browser")
	flag.Parse()
}

func main() {
	parseFlags()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, os.Kill)

	go func() {
		log.Debug("listening rpc")
		err := rpc.ServeRPC(rpcAddress)
		if err != nil {
			log.Error("rpc.ServeRPC: ", err)
			os.Exit(1)
		}
	}()

	log.Debug("listening web")
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(w, r)
	})
	ln, err := net.Listen("tcp", webSocketAddress)
	if err != nil {
		log.Error("net.Listen: ", err)
		os.Exit(1)
	}

	if openBrowser {
		go func() {
			browser.Open(fmt.Sprintf("http://%s", webSocketAddress))
		}()
	}
	go func() {
		err := http.Serve(ln, http.DefaultServeMux)
		if err != nil {
			log.Error("http.Serve: ", err)
			os.Exit(1)
		}
	}()

	select {
	case signal := <-osSignal:
		if signal == os.Interrupt {
			log.Debugln("exit by signal Interrupt")
		} else if signal == os.Kill {
			log.Debugln("exit by signal Kill")
		}
	}

}
