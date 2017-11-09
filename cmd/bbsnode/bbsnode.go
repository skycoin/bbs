package main

import (
	"encoding/json"
	"fmt"
	"github.com/skycoin/bbs/src/http"
	"github.com/skycoin/bbs/src/rpc"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/skycoin/src/util/browser"
	"github.com/skycoin/skycoin/src/util/file"
	"gopkg.in/urfave/cli.v1"
	"log"
	"os"
	"os/signal"
	"path/filepath"
)

const (
	defaultConfigSubDir    = ".skybbs"
	defaultStaticSubDir    = "static/dist"
	defaultDevStaticSubDir = "src/github.com/skycoin/bbs/static/dist"
	defaultRPCPort         = 8996
	defaultCXOPort         = 8998
	defaultCXORPCPort      = 8997
	defaultWebPort         = 8080
)

var (
	//defaultSubscriptions = []string{
	//	"03588a2c8085e37ece47aec50e1e856e70f893f7f802cb4f92d52c81c4c3212742",
	//}
	defaultMessengerAddresses = cli.StringSlice{
		"messenger.skycoin.net:8080",
	}
	devMode          = false
	compilerInternal = 1
)

// Config represents configuration for node.
type Config struct {
	Memory    bool   `json:"memory"`     // Whether to run node in memory.
	ConfigDir string `json:"config_dir"` // Full path for configuration directory.

	RPC     bool `json:"rpc"`      // Enable RPC interface for admin control.
	RPCPort int  `json:"rpc_port"` // Listening port of node RPC.

	CXOPort    int  `json:"cxo_port"`               // Listening port of CXO.
	CXORPC     bool `json:"cxo_rpc"`                // Whether to enable CXO RPC.
	CXORPCPort int  `json:"cxo_rpc_port,omitempty"` // Listening RPC port of CXO.

	EnforcedMessengerAddresses cli.StringSlice `json:"ensured_messenger_addresses"` // Addresses of messenger servers to enforce.
	EnforcedSubscriptions      cli.StringSlice `json:"ensured_subscriptions"`       // Subscriptions to enforce.

	WebPort   int    `json:"web_port"`              // Port to serve HTTP API/GUI.
	WebGUI    bool   `json:"web_gui"`               // Whether to enable GUI.
	WebGUIDir string `json:"web_gui_dir,omitempty"` // Full path of GUI static files.
	WebTLS    bool   `json:"web_tls"` // Whether to enable TLS.
	WebTLSCertFile string `json:"web_tls_cert_file"` // Path for TLS Certificate file.
	WebTLSKeyFile string `json:"web_tls_key_file"` // Path for TLS Key file.


	Browser bool `json:"browser"` // Whether to open browser on GUI start.
}

// NewDefaultConfig returns a default configuration for BBS node.
func NewDefaultConfig() *Config {
	return &Config{
		Memory:                     false, // Save to disk.
		ConfigDir:                  "",    // --> Action: set as '$HOME/.skybbs'
		RPC:                        true,
		RPCPort:                    defaultRPCPort,
		CXOPort:                    defaultCXOPort,
		CXORPC:                     false,
		CXORPCPort:                 defaultCXORPCPort,
		EnforcedMessengerAddresses: defaultMessengerAddresses,
		EnforcedSubscriptions:      []string{},
		WebPort:                    defaultWebPort,
		WebGUI:                     true,
		WebGUIDir:                  "", // --> Action: set as '$HOME/.skybbs/static'
		Browser:                    true,
	}
}

func (c *Config) Print() {
	data, _ := json.MarshalIndent(*c, "", "    ")
	log.Println(string(data))
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
	if c.WebGUI {
		if c.WebGUIDir == "" {
			if devMode {
				c.WebGUIDir = filepath.Join(os.Getenv("GOPATH"), defaultDevStaticSubDir)
			} else {
				c.WebGUIDir = defaultStaticSubDir
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
				Port:      &c.WebPort,
				StaticDir: &c.WebGUIDir,
				EnableGUI: &c.WebGUI,
				EnableTLS: &c.WebTLS,
				TLSCertFile: &c.WebTLSCertFile,
				TLSKeyFile: &c.WebTLSKeyFile,
			},
			&http.Gateway{
				Access: &store.Access{
					CXO: cxo.NewManager(
						&cxo.ManagerConfig{
							Memory: &c.Memory,
							Config: &c.ConfigDir,
							EnforcedMessengerAddresses: c.EnforcedMessengerAddresses,
							EnforcedSubscriptions:      c.EnforcedSubscriptions,
							CXOPort:                    &c.CXOPort,
							CXORPCEnable:               &c.CXORPC,
							CXORPCPort:                 &c.CXORPCPort,
						},
						&state.CompilerConfig{
							UpdateInterval: &compilerInternal,
						},
					),
				},
				QuitChan: quit,
			},
		)
		CatchError(e, "failed to start HTTP server")
		defer httpServer.Close()

		rpcServer, e := rpc.NewServer(
			&rpc.ServerConfig{
				Enable: &c.RPC,
				Port:   &c.RPCPort,
			},
			&rpc.Gateway{
				Access: &store.Access{
					CXO: httpServer.CXO(),
				},
				QuitChan: quit,
			},
		)
		CatchError(e, "failed to start RPC server")
		defer rpcServer.Close()

		if c.Browser {
			address := fmt.Sprintf("http://127.0.0.1:%d", c.WebPort)
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
			Name:        "memory",
			Destination: &config.Memory,
		},
		cli.StringFlag{
			Name:        "config-dir",
			Destination: &config.ConfigDir,
		},
		cli.BoolTFlag{
			Name:        "rpc",
			Destination: &config.RPC,
		},
		cli.IntFlag{
			Name:        "rpc-port",
			Destination: &config.RPCPort,
			Value:       config.RPCPort,
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
			Name:  "enforced-messenger-addresses",
			Value: &config.EnforcedMessengerAddresses,
		},
		cli.StringSliceFlag{
			Name:  "enforced-subscriptions",
			Value: &config.EnforcedSubscriptions,
		},
		cli.IntFlag{
			Name:        "web-port",
			Destination: &config.WebPort,
			Value:       config.WebPort,
		},
		cli.BoolTFlag{
			Name:        "web-gui",
			Destination: &config.WebGUI,
		},
		cli.StringFlag{
			Name:        "web-gui-dir",
			Destination: &config.WebGUIDir,
		},
		cli.BoolFlag{
			Name:        "web-tls",
			Destination: &config.WebTLS,
		},
		cli.StringFlag{
			Name: "web-tls-cert-file",
			Destination: &config.WebTLSCertFile,
			Value: config.WebTLSCertFile,
		},
		cli.StringFlag{
			Name: "web-tls-key-file",
			Destination: &config.WebTLSKeyFile,
			Value: config.WebTLSKeyFile,
		},
	}
	app := cli.NewApp()
	app.Name = "bbsnode"
	app.Usage = "Runs a Skycoin BBS Node"
	app.Flags = flags
	app.Action = config.GenerateAction()
	if e := app.Run(os.Args); e != nil {
		panic(e)
	}
}
