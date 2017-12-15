package main

import (
	"encoding/json"
	"fmt"
	"github.com/skycoin/bbs/src/http"
	"github.com/skycoin/bbs/src/msgs"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/session"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/skycoin/src/util/browser"
	"github.com/skycoin/skycoin/src/util/file"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

const (
	defaultConfigSubDir    = ".skybbs"
	defaultStaticSubDir    = "static/dist"
	defaultDevStaticSubDir = "src/github.com/skycoin/bbs/static/dist"
	defaultCXOPort         = 8998
	defaultCXORPCPort      = 8997
	defaultHTTPPort        = 7410
)

var (
	defaultMessengerAddresses = cli.StringSlice{
		"messenger.skycoin.net:9090",
	}
	defaultDevMessengerAddresses = cli.StringSlice{
		"127.0.0.1:8080",
	}
	devMode                    = false
	compilerInternal           = 1
	messengerReconnectInterval = time.Second * 3
)

// Config represents configuration for node.
type Config struct {
	Master    bool   `json:"master"`     // Whether to run node in master mode.
	Memory    bool   `json:"memory"`     // Whether to run node in memory.
	ConfigDir string `json:"config_dir"` // Full path for configuration directory.

	CXOPort    int  `json:"cxo_port"`               // Listening port of CXO.
	CXORPC     bool `json:"cxo_rpc"`                // Whether to enable CXO RPC.
	CXORPCPort int  `json:"cxo_rpc_port,omitempty"` // Listening RPC port of CXO.

	MessengerAddresses cli.StringSlice `json:"messenger_addresses"` // Addresses of messenger servers.

	HTTPPort   int    `json:"http_port"`              // Port to serve HTTP API/GUI.
	HTTPGUI    bool   `json:"http_gui"`               // Whether to enable GUI.
	HTTPGUIDir string `json:"http_gui_dir,omitempty"` // Full path of GUI static files.

	Browser  bool `json:"browser"`  // Whether to open browser on GUI start.
	Defaults bool `json:"defaults"` // Whether to have default connections/subscriptions.
}

// NewDefaultConfig returns a default configuration for BBS node.
func NewDefaultConfig() *Config {
	return &Config{
		Memory:             false, // Save to disk.
		ConfigDir:          "",    // --> Action: set as '$HOME/.skybbs'
		CXOPort:            defaultCXOPort,
		CXORPC:             false,
		CXORPCPort:         defaultCXORPCPort,
		MessengerAddresses: defaultMessengerAddresses,
		HTTPPort:           defaultHTTPPort,
		HTTPGUI:            true,
		HTTPGUIDir:         "", // --> Action: set as '$HOME/.skybbs/static'
		Browser:            true,
		Defaults:           true,
	}
}

func (c *Config) Print() {
	data, _ := json.MarshalIndent(*c, "", "    ")
	fmt.Println(string(data))
}

// PostProcess checks the flags and processes them.
func (c *Config) PostProcess() error {
	if !c.Memory {
		if c.ConfigDir == "" {
			c.ConfigDir = filepath.Join(file.UserHome(), defaultConfigSubDir)
		}
		if e := os.MkdirAll(c.ConfigDir, os.FileMode(0700)); e != nil {
			return e
		}
	}
	if devMode {
		c.MessengerAddresses = defaultDevMessengerAddresses
	}
	if c.HTTPGUI {
		if c.HTTPGUIDir == "" {
			if devMode {
				c.HTTPGUIDir = filepath.Join(os.Getenv("GOPATH"), defaultDevStaticSubDir)
			} else {
				c.HTTPGUIDir = defaultStaticSubDir
			}
		}
	} else {
		c.Browser = false
	}
	return nil
}

// GenerateAction generates a runnable action.
func (c *Config) GenerateAction() cli.ActionFunc {
	return func(_ *cli.Context) error {
		if e := c.PostProcess(); e != nil {
			return e
		}
		c.Print()

		quit := CatchInterrupt()
		defer log.Println("Goodbye.")

		httpServer, e := http.NewServer(
			&http.ServerConfig{
				Port:      &c.HTTPPort,
				StaticDir: &c.HTTPGUIDir,
				EnableGUI: &c.HTTPGUI,
			},
			&http.Gateway{
				Access: &store.Access{
					Session: session.NewManager(
						&session.ManagerConfig{
							MemoryMode: &c.Memory,
							ConfigDir:  &c.ConfigDir,
						},
					),
					CXO: cxo.NewManager(
						&cxo.ManagerConfig{
							Memory:             &c.Memory,
							Config:             &c.ConfigDir,
							MessengerAddresses: c.MessengerAddresses,
							CXOPort:            &c.CXOPort,
							CXORPCEnable:       &c.CXORPC,
							CXORPCPort:         &c.CXORPCPort,
							Defaults:           &c.Defaults,
						},
						&state.CompilerConfig{
							UpdateInterval: &compilerInternal,
						},
						&msgs.RelayConfig{
							Addresses:         c.MessengerAddresses,
							ReconnectInterval: &messengerReconnectInterval,
						},
					),
				},
				Quit: quit,
			},
		)
		CatchError(e, "failed to start HTTP Server")
		defer httpServer.Close()

		if c.Browser {
			address := fmt.Sprintf("http://127.0.0.1:%d", c.HTTPPort)
			log.Println("Opening browser at address:", address)
			if e := browser.Open(address); e != nil {
				log.Println("Error on browser open:", e)
			}
		}

		<-quit
		return nil
	}
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

func main() {
	config := NewDefaultConfig()
	flags := []cli.Flag{
		cli.BoolFlag{
			Name:        "dev",
			Destination: &devMode,
		},
		cli.BoolFlag{
			Name:        "master",
			Destination: &config.Master,
		},
		cli.BoolFlag{
			Name:        "memory",
			Destination: &config.Memory,
		},
		cli.StringFlag{
			Name:        "config-dir",
			Destination: &config.ConfigDir,
		},
		cli.IntFlag{
			Name:        "cxo-port",
			Destination: &config.CXOPort,
			Value:       config.CXOPort,
		},
		cli.BoolTFlag{
			Name:        "cxo-rpc",
			Destination: &config.CXORPC,
		},
		cli.IntFlag{
			Name:        "cxo-rpc-port",
			Destination: &config.CXORPCPort,
			Value:       config.CXORPCPort,
		},
		cli.StringSliceFlag{
			Name:  "messenger-addresses",
			Value: &config.MessengerAddresses,
		},
		cli.IntFlag{
			Name:        "http-port",
			Destination: &config.HTTPPort,
			Value:       config.HTTPPort,
		},
		cli.BoolTFlag{
			Name:        "http-gui",
			Destination: &config.HTTPGUI,
		},
		cli.StringFlag{
			Name:        "http-gui-dir",
			Destination: &config.HTTPGUIDir,
		},
		cli.BoolTFlag{
			Name:        "defaults",
			Destination: &config.Defaults,
		},
	}
	app := cli.NewApp()
	app.Name = "Skycoin BBS Node"
	app.Usage = "Runs a Skycoin BBS Node"
	app.Version = "4.4"
	app.Flags = flags
	app.Action = config.GenerateAction()
	if e := app.Run(os.Args); e != nil {
		panic(e)
	}
}
