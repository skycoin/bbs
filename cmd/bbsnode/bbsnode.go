package main

import (
	"encoding/json"
	"fmt"
	"github.com/skycoin/bbs/src/http"
	"github.com/skycoin/bbs/src/rpc"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/medial"
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
	Version = "5.1"

	defaultConfigSubDir                    = ".skybbs"
	defaultStaticSubDir                    = "static/dist"
	defaultRPCPort                         = 8996
	defaultCXOPort                         = 8998
	defaultCXORPCPort                      = 8997
	defaultWebPort                         = 8080
	defaultMedialGarbageCollectionInterval = time.Minute
	defaultMedialItemTimeout               = time.Minute * 3
)

var (
	compilerInternal = 1
)

// Config represents configuration for node.
type Config struct {
	Public                     bool            `json:"public"`                       // Whether to expose node publicly.
	Memory                     bool            `json:"memory"`                       // Whether to run node in memory.
	ConfigDir                  string          `json:"config-dir"`                   // Full path for configuration directory.
	RPC                        bool            `json:"rpc"`                          // Enable RPC interface for admin control.
	RPCPort                    int             `json:"rpc-port"`                     // Listening port of node RPC.
	CXOPort                    int             `json:"cxo-port"`                     // Listening port of CXO.
	CXORPC                     bool            `json:"cxo-rpc"`                      // Whether to enable CXO RPC.
	CXORPCPort                 int             `json:"cxo-rpc-port,omitempty"`       // Listening RPC port of CXO.
	EnforcedMessengerAddresses cli.StringSlice `json:"enforced-messenger-addresses"` // Addresses of messenger servers to enforce.
	EnforcedSubscriptions      cli.StringSlice `json:"enforced-subscriptions"`       // Subscriptions to enforce.
	WebPort                    int             `json:"web-port"`                     // Port to serve HTTP API/GUI.
	WebGUI                     bool            `json:"web-gui"`                      // Whether to enable GUI.
	WebGUIDir                  string          `json:"web-gui-dir,omitempty"`        // Full path of GUI static files.
	WebTLS                     bool            `json:"web-tls"`                      // Whether to enable TLS.
	WebTLSCertFile             string          `json:"web-tls-cert-file"`            // Path for TLS Certificate file.
	WebTLSKeyFile              string          `json:"web-tls-key-file"`             // Path for TLS Key file.
	Browser                    bool            `json:"open-browser"`                 // Whether to open browser on GUI start.
}

// NewDefaultConfig returns a default configuration for BBS node.
func NewDefaultConfig() *Config {
	return &Config{
		Public:                     false,
		Memory:                     false, // Save to disk.
		ConfigDir:                  "",    // --> Action: set as '$HOME/.skybbs'
		RPC:                        true,
		RPCPort:                    defaultRPCPort,
		CXOPort:                    defaultCXOPort,
		CXORPC:                     false,
		CXORPCPort:                 defaultCXORPCPort,
		EnforcedMessengerAddresses: []string{},
		EnforcedSubscriptions:      []string{},
		WebPort:                    defaultWebPort,
		WebGUI:                     true,
		WebGUIDir:                  defaultStaticSubDir, // --> Action: set as '$HOME/.skybbs/static/dist'
		Browser:                    false,
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
				Port:        &c.WebPort,
				StaticDir:   &c.WebGUIDir,
				EnableGUI:   &c.WebGUI,
				EnableTLS:   &c.WebTLS,
				TLSCertFile: &c.WebTLSCertFile,
				TLSKeyFile:  &c.WebTLSKeyFile,
			},
			&http.Gateway{
				Access: &store.Access{
					CXO: cxo.NewManager(
						&cxo.ManagerConfig{
							Public: &c.Public,
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
					Medial: medial.NewServer(&medial.ServerConfig{
						GarbageCollectionInterval: defaultMedialGarbageCollectionInterval,
						ItemTimeoutInterval:       defaultMedialItemTimeout,
					}),
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

		if c.WebGUI && c.Browser {
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
			Name:        "public",
			Destination: &config.Public,
			Usage:       "whether this node is exposed publicly and shares it's subscribed boards",
		},
		cli.BoolFlag{
			Name:        "memory",
			Destination: &config.Memory,
			Usage:       "whether to avoid storing BBS data on disk and use memory instead",
		},
		cli.StringFlag{
			Name:        "config-dir",
			Destination: &config.ConfigDir,
			Usage:       "the name of the directory to store and access BBS configuration and associated cxo data (if left blank, $HOME/.skybbs will be used)",
		},
		cli.BoolTFlag{
			Name:        "rpc",
			Destination: &config.RPC,
			Usage:       "whether to enable RPC interface to interact with BBS node (used for bbscli) (default: true)",
		},
		cli.IntFlag{
			Name:        "rpc-port",
			Destination: &config.RPCPort,
			Value:       config.RPCPort,
			Usage:       "port to serve BBS RPC interface",
		},
		cli.IntFlag{
			Name:        "cxo-port",
			Destination: &config.CXOPort,
			Value:       config.CXOPort,
			Usage:       "port to listen for CXO connections",
		},
		cli.BoolTFlag{
			Name:        "cxo-rpc",
			Destination: &config.CXORPC,
			Usage:       "whether to enable RPC interface to interact with CXO (used for cxocli)",
		},
		cli.IntFlag{
			Name:        "cxo-rpc-port",
			Destination: &config.CXORPCPort,
			Value:       config.CXORPCPort,
			Usage:       "port to serve CXO RPC interface",
		},
		cli.StringSliceFlag{
			Name:  "enforced-messenger-addresses",
			Value: &config.EnforcedMessengerAddresses,
			Usage: "list of addresses to messenger servers to enforce connections with",
		},
		cli.StringSliceFlag{
			Name:  "enforced-subscriptions",
			Value: &config.EnforcedSubscriptions,
			Usage: "list of public keys of boards to enforce subscriptions with",
		},
		cli.IntFlag{
			Name:        "web-port",
			Destination: &config.WebPort,
			Value:       config.WebPort,
			Usage:       "port to serve http api",
		},
		cli.BoolTFlag{
			Name:        "web-gui",
			Destination: &config.WebGUI,
			Usage:       "whether to enable web interface thin client",
		},
		cli.StringFlag{
			Name:        "web-gui-dir",
			Destination: &config.WebGUIDir,
			Usage:       "directory where web interface static files are located",
		},
		cli.BoolFlag{
			Name:        "web-tls",
			Destination: &config.WebTLS,
			Usage:       "whether to enable https for web interface thin client and api",
		},
		cli.StringFlag{
			Name:        "web-tls-cert-file",
			Destination: &config.WebTLSCertFile,
			Value:       config.WebTLSCertFile,
			Usage:       "path of the tls certificate file",
		},
		cli.StringFlag{
			Name:        "web-tls-key-file",
			Destination: &config.WebTLSKeyFile,
			Value:       config.WebTLSKeyFile,
			Usage:       "path of the tls key file",
		},
		cli.BoolFlag{
			Name:        "open-browser",
			Destination: &config.Browser,
			Usage:       "whether to open a browser window",
		},
	}
	app := cli.NewApp()
	app.Name = "bbsnode"
	app.Version = Version
	app.Usage = "Runs a Skycoin BBS Node"
	app.Flags = flags
	app.Action = config.GenerateAction()
	if e := app.Run(os.Args); e != nil {
		panic(e)
	}
}
