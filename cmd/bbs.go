package cmd

import (
	"flag"
	"fmt"
	"github.com/evanlinjin/bbs/datastore"
	"github.com/evanlinjin/bbs/gui"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"os"
	"os/signal"
	"time"
)

const (
	LocalhostAddress = "127.0.0.1"
)

// CXOConfig represents commandline arguments.
type Config struct {
	// Master determines whether BBS Server is in Master Mode.
	Master bool

	// RPC.
	RPCPort int

	// Address of CXO Daemon.
	CXOAddress string

	// Localhost web interface.
	WebInterface     bool
	WebInterfacePort int

	// Launch System Default Browser after client startup
	LaunchBrowser bool
}

// MakeConfig makes a default configuration.
func MakeConfig() *Config {
	//pk, sc := cipher.GenerateKeyPair()
	return &Config{
		Master: false,
		//PublicKey:        "",
		//SecretKey:        "",
		CXOAddress:       "[::]:8998",
		WebInterface:     true,
		WebInterfacePort: 6420,
		LaunchBrowser:    true,
	}
}

func (c *Config) register() {
	// Master mode.
	flag.BoolVar(&c.Master, "master", c.Master, "whether node is started as master")
	// RPC Port (Only enabled if Master mode).
	flag.IntVar(&c.RPCPort, "rpc-port", c.RPCPort, "port number for RPC")
	// CXO Address.
	flag.StringVar(&c.CXOAddress, "cxo-address", c.CXOAddress, "address of cxo daemon to connect to")
	// Web Interface.
	flag.BoolVar(&c.WebInterface, "web-interface", c.WebInterface, "enable the web interface")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	// Launch Browser.
	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser, "launch system default webbrowser at client startup")
}

func (c *Config) postProcess() {
}

func (c *Config) Parse() {
	c.register()
	flag.Parse()
	c.postProcess()
}

func catchInterrupt(quit chan<- int) {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	signal.Stop(sigchan)
	quit <- 1
}

func configureDatastore(c *Config) *datastore.CXOConfig {
	dc := datastore.NewCXOConfig()
	dc.Master = c.Master
	dc.Address = c.CXOAddress
	return dc
}

func Run(c *Config) {
	host := fmt.Sprintf("%s:%d", LocalhostAddress, c.WebInterfacePort)
	fullAddress := fmt.Sprintf("%s://%s", "http", host)
	fmt.Println("[FULL ADDRESS]", fullAddress)

	// If the user Ctrl-C's, shutdown properly.
	quit := make(chan int)
	go catchInterrupt(quit)

	// Datastore.
	cxoClientConfig := configureDatastore(c)
	cxoClient, e := datastore.NewCXOClient(cxoClientConfig)
	panicIfError(e, "unable to create CXOClient")
	panicIfError(cxoClient.AddRandomIdentity(), "unable to create random indentity for CXOClient")
	panicIfError(cxoClient.Launch(), "unable to launch CXOClient")

	// Start web interface.
	if c.WebInterface {
		if e := gui.LaunchWebInterface(host, "", cxoClient); e != nil {
			fmt.Println("[FAILED START]", e)
			os.Exit(1)
		}
	}

	// Launch browser.
	if c.LaunchBrowser {
		// Wait a moment just to make sure the http interface is up
		time.Sleep(time.Millisecond * 100)

		fmt.Println("[BROWSER LAUNCH] Address:", fullAddress)
		if e := util.OpenBrowser(fullAddress); e != nil {
			fmt.Println("[BROWSER LAUNCH] Error:", e)
		}
	}

	// Wait for Ctrl-C signal.
	<-quit
	fmt.Println("Shutting down...")
	cxoClient.Shutdown()
	gui.Shutdown()
	fmt.Println("Goodbye.")
}

func panicIfError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}
