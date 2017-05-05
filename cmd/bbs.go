package cmd

import (
	"flag"
	"fmt"
	"github.com/evanlinjin/bbs/cxo"
	"github.com/evanlinjin/bbs/gui"
	"github.com/evanlinjin/bbs/rpc"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
)

const (
	LocalhostAddress = "127.0.0.1"
)

// CXOConfig represents commandline arguments.
type Config struct {
	// Master determines whether BBS Server is in Master Mode.
	Master bool
	// ConfigDir determines the directory where user and board configuration are to be stored.
	ConfigDir string
	// RPC.
	RPCPort int
	// Port of CXO Daemon.
	CXOPort int
	// Localhost web interface.
	WebInterface     bool
	WebInterfacePort int
	// Launch System Default Browser after client startup
	LaunchBrowser bool
}

// NewConfig makes a default configuration.
func NewConfig() *Config {
	//pk, sc := cipher.GenerateKeyPair()
	return &Config{
		Master:           false,
		ConfigDir:        ".",
		RPCPort:          6421,
		CXOPort:          8998,
		WebInterface:     true,
		WebInterfacePort: 6420,
		LaunchBrowser:    true,
	}
}

func (c *Config) register() {
	// Master mode.
	flag.BoolVar(&c.Master, "master", c.Master, "whether node is started as master")
	// Configuration directory.
	flag.StringVar(&c.ConfigDir, "config-dir", c.ConfigDir, "directory for configuration files")
	// RPC Port (Only enabled if Master mode).
	flag.IntVar(&c.RPCPort, "rpc-port", c.RPCPort, "port number for RPC")
	// CXO Address.
	flag.IntVar(&c.CXOPort, "cxo-port", c.CXOPort, "port of cxo daemon to connect to")
	// Web Interface.
	flag.BoolVar(&c.WebInterface, "web-interface", c.WebInterface, "enable the web interface")
	flag.IntVar(&c.WebInterfacePort, "web-interface-port", c.WebInterfacePort, "port to serve web interface on")
	// Launch Browser.
	flag.BoolVar(&c.LaunchBrowser, "launch-browser", c.LaunchBrowser, "launch system default webbrowser at client startup")
}

func (c *Config) postProcess() {
	//os.MkdirAll()
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

func configureCXO(c *Config) *cxo.CXOConfig {
	dc := cxo.NewCXOConfig()
	dc.Master = c.Master
	dc.Port = c.CXOPort
	return dc
}

func Run(c *Config) {
	fmt.Println("[MASTER MODE]", c.Master)

	host := fmt.Sprintf("%s:%d", LocalhostAddress, c.WebInterfacePort)
	fullAddress := fmt.Sprintf("%s://%s", "http", host)
	fmt.Println("[FULL ADDRESS]", fullAddress)

	// If the user Ctrl-C's, shutdown properly.
	quit := make(chan int)
	go catchInterrupt(quit)

	// Config files.
	fmt.Println("[CONFIG DIRECTORY]", util.InitDataDir(c.ConfigDir))

	// Datastore.
	cxoConfig := configureCXO(c)
	cxoClient, e := cxo.NewClient(cxoConfig)
	panicIfError(
		e,
		"unable to create CXOClient",
	)
	panicIfError(
		cxoClient.Launch(),
		"unable to launch CXOClient",
	)

	// RPC Server.
	var rpcServer *rpc.Server
	if c.Master {
		rpcServer = rpc.NewServer(cxoClient)
		panicIfError(
			rpcServer.Launch("[::]:"+strconv.Itoa(c.RPCPort)),
			"unable to start rpc server",
		)
		fmt.Println("[RPC SERVER] Address:", rpcServer.Address())
	}

	// Start web interface.
	if c.WebInterface {
		gateway := cxo.NewGateWay(cxoClient)
		if e := gui.LaunchWebInterface(host, "", gateway); e != nil {
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
	gui.Shutdown()
	if c.Master {
		rpcServer.Shutdown()
	}
	cxoClient.Shutdown()
	fmt.Println("Goodbye.")
}

func panicIfError(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Panicf(msg+": %v", append(args, err)...)
	}
}
