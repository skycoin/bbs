package main

import (
	"github.com/skycoin/bbs/src/gui"
	"github.com/skycoin/bbs/src/rpc"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/msg"
	"github.com/skycoin/skycoin/src/util/browser"
	"github.com/skycoin/skycoin/src/util/file"
	"log"
	"os"
	"os/signal"
	"time"
	"flag"
	"path/filepath"
	"github.com/skycoin/bbs/src/misc"
	"strconv"
	"fmt"
	"github.com/pkg/errors"
)

const configSubDir = ".skybbs"
const webSubDir = "src/github.com/skycoin/bbs/static/dist"

// Config represents commandline arguments.
type Config struct {
	// [TEST MODE] enforces the following behaviours:
	// - `MemoryMode = false` (disables modification to cxo database, uses temp file instead).
	// - `SaveConfig = false` (disables modification to config files).

	TestMode            bool // Whether to enable test mode.
	TestModeThreads     int  // Number of threads to use for test mode (will create them in test mode).
	TestModeUsers       int  // Number of Master users used for simulated activity.
	TestModeMinInterval int  // Minimum interval between simulated activity (in seconds).
	TestModeMaxInterval int  // Maximum interval between simulated activity (in seconds).
	TestModeTimeOut     int  // Will stop simulated activity after this time (in seconds). Disabled if negative.
	TestModePostCap     int  // Maximum number of posts allowed. Disabled if negative.

	Master    bool   // Whether BBS node can host boards.
	ConfigDir string // Configuration directory.
	RPCPort   int    // RPC server port (Master node only).
	RPCRemAdr string // RPC remote address (Master node only).

	CXOPort    int  // Port of CXO Daemon.
	CXORPCPort int  // Port of CXO Daemon's RPC.
	MemoryMode bool // Whether to use in-memory database for CXO.

	WebGUIEnable      bool   // Whether to enable web GUI.
	WebGUIPort        int    // Port of web GUI.
	WebGUIDir         string // Root directory that has the index.html file.
	WebGUIOpenBrowser bool   // Whether to open browser on web GUI start.
}

// NewConfig makes Config with default values.
func NewConfig() *Config {
	return &Config{
		TestMode:            false,
		TestModeThreads:     3,
		TestModeUsers:       1,
		TestModeMinInterval: 1,
		TestModeMaxInterval: 10,
		TestModeTimeOut:     -1,
		TestModePostCap:     -1,

		Master:    false,
		ConfigDir: "",
		RPCPort:   6421,
		RPCRemAdr: "",

		CXOPort:    8998,
		CXORPCPort: 8997,

		MemoryMode: false,

		WebGUIEnable:      true,
		WebGUIPort:        7410,
		WebGUIDir:         "",
		WebGUIOpenBrowser: true,
	}
}

// Parse fills the Config with commandline argument values.
func (c *Config) Parse() *Config {
	/*
		<<< TEST FLAGS >>>
	*/

	flag.BoolVar(&c.TestMode,
		"test-mode", c.TestMode,
		"whether to enable test mode")

	flag.IntVar(&c.TestModeThreads,
		"test-mode-threads", c.TestModeThreads,
		"number of threads to use for test mode")

	flag.IntVar(&c.TestModeUsers,
		"test-mode-users", c.TestModeUsers,
		"number of users to use for test mode")

	flag.IntVar(&c.TestModeMinInterval,
		"test-mode-min", c.TestModeMinInterval,
		"minimum interval in seconds between simulated activity")

	flag.IntVar(&c.TestModeMaxInterval,
		"test-mode-max", c.TestModeMaxInterval,
		"maximum interval in seconds between simulated activity")

	flag.IntVar(&c.TestModeTimeOut,
		"test-mode-timeout", c.TestModeTimeOut,
		"time in seconds before simulated activity stops - disabled if negative")

	flag.IntVar(&c.TestModePostCap,
		"test-mode-post-cap", c.TestModePostCap,
		"maximum number of posts allowed to be created - disabled if negative")

	/*
		<<< BBS FLAGS >>>
	*/

	flag.BoolVar(&c.Master,
		"Master", c.Master,
		"whether to enable bbs node to host boards")

	flag.StringVar(&c.ConfigDir,
		"config-dir", c.ConfigDir,
		"configuration directory - set to $HOME/.skycoin/bbs if left empty")

	flag.IntVar(&c.RPCPort,
		"rpc-port", c.RPCPort,
		"port of rpc server for Master node")

	flag.StringVar(&c.RPCRemAdr,
		"rpc-remote-address", c.RPCRemAdr,
		"remote address of rpc server for Master node")

	/*
		<<< CXO FLAGS >>>
	*/

	flag.IntVar(&c.CXOPort,
		"cxo-port", c.CXOPort,
		"port of cxo daemon to connect to")

	flag.IntVar(&c.CXORPCPort,
		"cxo-rpc-port", c.CXORPCPort,
		"port of cxo daemon rpc to connect to")

	flag.BoolVar(&c.MemoryMode,
		"cxo-memory-mode", c.MemoryMode,
		"whether to use in-memory database")

	/*
		<<< WEB GUI FLAGS >>>
	*/

	flag.BoolVar(&c.WebGUIEnable,
		"web-gui-enable", c.WebGUIEnable,
		"whether to enable the web gui")

	flag.IntVar(&c.WebGUIPort,
		"web-gui-port", c.WebGUIPort,
		"local port to serve web gui on")

	flag.StringVar(&c.WebGUIDir,
		"web-gui-dir", c.WebGUIDir,
		"root directory of index.html file")

	flag.BoolVar(&c.WebGUIOpenBrowser,
		"web-gui-open-browser", c.WebGUIOpenBrowser,
		"whether to open browser after web gui is ready")

	flag.Parse()
	return c
}

// PostProcess checks the validity and post processes the flags.
func (c *Config) PostProcess() (*Config, error) {
	// Action on test mode.
	if c.TestMode {
		// Check test mode settings.
		if c.TestModeThreads < 0 {
			return nil, errors.New("invalid number of test mode threads specified")
		}
		if c.TestModeUsers < 1 {
			return nil, errors.New("invalid number of test mode users specified")
		}
		if c.TestModeMinInterval < 0 {
			return nil, errors.New("invalid test mode minimum interval specified")
		}
		if c.TestModeMaxInterval < 0 {
			return nil, errors.New("invalid test mode maximum interval specified")
		}
		if c.TestModeMinInterval > c.TestModeMaxInterval {
			return nil, errors.New("test mode minimum interval > maximum interval")
		}
		// Enforce behaviour.
		c.Master = true
		c.WebGUIEnable = true
		c.MemoryMode = false
	}
	// Configure configuration directories if necessary.
	if !c.MemoryMode {
		// Action on BBS configuration files.
		if c.ConfigDir == "" {
			c.ConfigDir = filepath.Join(file.UserHome(), configSubDir)
		}
		// Ensure directories exist.
		if e := os.MkdirAll(c.ConfigDir, os.FileMode(0700)); e != nil {
			return nil, e
		}
	}
	// Master mode stuff.
	if c.Master && c.RPCRemAdr == "" {
		c.RPCRemAdr = misc.GetIP() + ":" + strconv.Itoa(c.RPCPort)
		fmt.Println("External Addr:", c.RPCRemAdr)
	}
	// Web interface.
	if c.WebGUIDir == "" {
		c.WebGUIDir = filepath.Join(os.Getenv("GOPATH"), webSubDir)
		fmt.Println("Web Dir:", c.WebGUIDir)
	}
	return c, nil
}

func main() {
	quit := CatchInterrupt()
	config, e := NewConfig().Parse().PostProcess()
	if e != nil {
		panic(e)
	}
	file.InitDataDir(config.ConfigDir)
	log.Println("[CONFIG] Master mode:", config.Master)
	defer log.Println("Goodbye.")

	log.Println("[CONFIG] Connecting to cxo on port", config.CXOPort)

	storeConfig := &store.Config{
		Master:     config.Master,
		TestMode:   config.TestMode,
		MemoryMode: config.MemoryMode,
		ConfigDir:  config.ConfigDir,
		CXOPort:    config.CXOPort,
		CXORPCPort: config.CXORPCPort,
	}
	container, e := store.NewCXO(storeConfig)
	CatchError(e, "unable to create cxo container")
	defer container.Close()

	boardSaver, e := store.NewBoardSaver(storeConfig, container)
	CatchError(e, "unable to create board saver")
	defer boardSaver.Close()

	userSaver, e := store.NewUserSaver(storeConfig, container)
	CatchError(e, "unable to create user saver")

	_, e = store.NewFirstRunSaver(storeConfig, boardSaver)
	CatchError(e, "unable to create first run saver")

	queueSaver, e := msg.NewQueueSaver(storeConfig, container)
	CatchError(e, "unable to create queue saver")
	defer queueSaver.Close()

	var rpcServer *rpc.Server
	if config.Master {
		rpcServer, e = rpc.NewServer(
			rpc.NewGateway(container, boardSaver, userSaver),
			config.RPCPort,
		)
		CatchError(e, "unable to start rpc server")
		defer rpcServer.Close()

		log.Println("[RPCSERVER] Serving on address:", rpcServer.Address())
	}

	httpConfig := &gui.HTTPConfig{
		RPCRemoteAddr: config.RPCRemAdr,
		Port:          config.WebGUIPort,
		StaticDir:     config.WebGUIDir,
		EnableGUI:     config.WebGUIEnable,
	}

	gateway := gui.NewGateway(
		httpConfig, container, boardSaver, userSaver, queueSaver, quit)

	serveAddr, e := gui.OpenWebInterface(httpConfig, gateway)
	CatchError(e, "unable to start web server")
	defer gui.Close()

	if config.TestMode {
		testConfig := &gui.TesterConfig{
			ThreadCount: config.TestModeThreads,
			UsersCount:  config.TestModeUsers,
			PostCap:     config.TestModePostCap,
			MinInterval: config.TestModeMinInterval,
			MaxInterval: config.TestModeMaxInterval,
			Timeout:     config.TestModeTimeOut,
		}

		tester, e := gui.NewTester(testConfig, gateway)
		CatchError(e, "unable to start tester")
		defer tester.Close()
	}

	log.Println("[WEBGUI] Serving on:", serveAddr)

	if config.WebGUIEnable && config.WebGUIOpenBrowser {
		go func() {
			time.Sleep(time.Millisecond * 100)
			log.Println("Opening web browser...")
			browser.Open(serveAddr)
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
