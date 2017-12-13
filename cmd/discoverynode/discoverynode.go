package main

import (
	"flag"
	"os"
	"os/signal"

	"log"
	"github.com/skycoin/net/skycoin-messenger/factory"
)

var (
	address string
)

func parseFlags() {
	flag.StringVar(&address, "address", ":8080", "address to listen on")
	flag.Parse()
}

func main() {
	parseFlags()

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Interrupt, os.Kill)

	f := factory.NewMessengerFactory()
	f.SetLoggerLevel(factory.DebugLevel)

	log.Printf("listening on %s", address)

	if e := f.Listen(address); e != nil {
		log.Println(e)
		os.Exit(1)
	}

	select {
	case s := <-osSignal:
		switch s {
		case os.Interrupt:
			log.Printf("exit by signal Interrupt")
		case os.Kill:
			log.Printf("exit by signal Kill")
		}
	}
}
