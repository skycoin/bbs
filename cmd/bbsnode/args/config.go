package args

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/skycoin/src/util"
	"os"
	"path/filepath"
	"strconv"
)

const configSubDir = ".skybbs"
const webSubDir = "src/github.com/skycoin/bbs/static/dist"

// Config represents commandline arguments.
type Config struct {

	// [TEST MODE] enforces the following behaviours:
	// - `cxoMemoryMode = false` (disables modification to cxo database, uses temp file instead).
	// - `saveConfig = false` (disables modification to config files).

	testMode            bool // Whether to enable test mode.
	testModeThreads     int  // Number of threads to use for test mode (will create them in test mode).
	testModeUsers       int  // Number of master users used for simulated activity.
	testModeMinInterval int  // Minimum interval between simulated activity (in seconds).
	testModeMaxInterval int  // Maximum interval between simulated activity (in seconds).
	testModeTimeOut     int  // Will stop simulated activity after this time (in seconds). Disabled if negative.
	testModePostCap     int  // Maximum number of posts allowed. Disabled if negative.

	master     bool   // Whether BBS node can host boards.
	saveConfig bool   // Whether to save and use BBS configuration files.
	configDir  string // Configuration directory.
	rpcPort    int    // RPC server port (master node only).
	rpcRemAdr  string // RPC remote address (master node only).

	cxoPort       int    // Port of CXO Daemon.
	cxoRPCPort    int    // Port of CXO Daemon's RPC.
	cxoMemoryMode bool   // Whether to use in-memory database for CXO.
	cxoDir        string // Directory to store cxo databases.

	webGUIEnable      bool   // Whether to enable web GUI.
	webGUIPort        int    // Port of web GUI.
	webGUIDir         string // Root directory that has the index.html file.
	webGUIOpenBrowser bool   // Whether to open browser on web GUI start.
}

// NewConfig makes Config with default values.
func NewConfig() *Config {
	return &Config{
		testMode:            false,
		testModeThreads:     3,
		testModeUsers:       1,
		testModeMinInterval: 1,
		testModeMaxInterval: 10,
		testModeTimeOut:     -1,
		testModePostCap:     -1,

		master:     false,
		saveConfig: true,
		configDir:  "",
		rpcPort:    6421,
		rpcRemAdr:  "",

		cxoPort:       8998,
		cxoRPCPort:    8997,
		cxoMemoryMode: false,
		cxoDir:        "",

		webGUIEnable:      true,
		webGUIPort:        7410,
		webGUIDir:         "",
		webGUIOpenBrowser: true,
	}
}

// Parse fills the Config with commandline argument values.
func (c *Config) Parse() *Config {
	/*
		<<< TEST FLAGS >>>
	*/

	flag.BoolVar(&c.testMode,
		"test-mode", c.testMode,
		"whether to enable test mode")

	flag.IntVar(&c.testModeThreads,
		"test-mode-threads", c.testModeThreads,
		"number of threads to use for test mode")

	flag.IntVar(&c.testModeUsers,
		"test-mode-users", c.testModeUsers,
		"number of users to use for test mode")

	flag.IntVar(&c.testModeMinInterval,
		"test-mode-min", c.testModeMinInterval,
		"minimum interval in seconds between simulated activity")

	flag.IntVar(&c.testModeMaxInterval,
		"test-mode-max", c.testModeMaxInterval,
		"maximum interval in seconds between simulated activity")

	flag.IntVar(&c.testModeTimeOut,
		"test-mode-timeout", c.testModeTimeOut,
		"time in seconds before simulated activity stops - disabled if negative")

	flag.IntVar(&c.testModePostCap,
		"test-mode-post-cap", c.testModePostCap,
		"maximum number of posts allowed to be created - disabled if negative")

	/*
		<<< BBS FLAGS >>>
	*/

	flag.BoolVar(&c.master,
		"master", c.master,
		"whether to enable bbs node to host boards")

	flag.BoolVar(&c.saveConfig,
		"save-config", c.saveConfig,
		"whether to save and use configuration files")

	flag.StringVar(&c.configDir,
		"config-dir", c.configDir,
		"configuration directory - set to $HOME/.skycoin/bbs if left empty")

	flag.IntVar(&c.rpcPort,
		"rpc-port", c.rpcPort,
		"port of rpc server for master node")

	flag.StringVar(&c.rpcRemAdr,
		"rpc-remote-address", c.rpcRemAdr,
		"remote address of rpc server for master node")

	/*
		<<< CXO FLAGS >>>
	*/

	flag.IntVar(&c.cxoPort,
		"cxo-port", c.cxoPort,
		"port of cxo daemon to connect to")

	flag.IntVar(&c.cxoRPCPort,
		"cxo-rpc-port", c.cxoRPCPort,
		"port of cxo daemon rpc to connect to")

	flag.BoolVar(&c.cxoMemoryMode,
		"cxo-memory-mode", c.cxoMemoryMode,
		"whether to use in-memory database")

	flag.StringVar(&c.cxoDir,
		"cxo-dir", c.cxoDir,
		"directory to store cxo databases - uses default if empty")

	/*
		<<< WEB GUI FLAGS >>>
	*/

	flag.BoolVar(&c.webGUIEnable,
		"web-gui-enable", c.webGUIEnable,
		"whether to enable the web gui")

	flag.IntVar(&c.webGUIPort,
		"web-gui-port", c.webGUIPort,
		"local port to serve web gui on")

	flag.StringVar(&c.webGUIDir,
		"web-gui-dir", c.webGUIDir,
		"root directory of index.html file")

	flag.BoolVar(&c.webGUIOpenBrowser,
		"web-gui-open-browser", c.webGUIOpenBrowser,
		"whether to open browser after web gui is ready")

	flag.Parse()
	return c
}

// PostProcess checks the validity and post processes the flags.
func (c *Config) PostProcess() (*Config, error) {
	// Action on test mode.
	if c.testMode {
		// Check test mode settings.
		if c.testModeThreads < 0 {
			return nil, errors.New("invalid number of test mode threads specified")
		}
		if c.testModeUsers < 1 {
			return nil, errors.New("invalid number of test mode users specified")
		}
		if c.testModeMinInterval < 0 {
			return nil, errors.New("invalid test mode minimum interval specified")
		}
		if c.testModeMaxInterval < 0 {
			return nil, errors.New("invalid test mode maximum interval specified")
		}
		if c.testModeMinInterval > c.testModeMaxInterval {
			return nil, errors.New("test mode minimum interval > maximum interval")
		}
		// Enforce behaviour.
		c.master = true
		c.webGUIEnable = true
		c.cxoMemoryMode = false
		c.cxoDir = ""
		c.saveConfig = false
	}
	// Configure configuration directories if necessary.
	if c.saveConfig {
		// Action on BBS configuration files.
		if c.configDir == "" {
			c.configDir = filepath.Join(util.UserHome(), configSubDir)
		}
		// Action on CXO configuration files.
		if c.cxoDir == "" {
			c.cxoDir = filepath.Join(c.configDir, "cxo")
		}
		// Ensure directories exist.
		if e := os.MkdirAll(c.configDir, os.FileMode(0700)); e != nil {
			return nil, e
		}
		if e := os.MkdirAll(c.cxoDir, os.FileMode(0700)); e != nil {
			return nil, e
		}
	}
	// Master mode stuff.
	if c.Master() && c.rpcRemAdr == "" {
		c.rpcRemAdr = misc.GetIP() + ":" + strconv.Itoa(c.rpcPort)
		fmt.Println("External Addr:", c.rpcRemAdr)
	}
	// Web interface.
	if c.webGUIDir == "" {
		c.webGUIDir = filepath.Join(os.Getenv("GOPATH"), webSubDir)
		fmt.Println("Web Dir:", c.webGUIDir)
	}
	return c, nil
}

/*
	These functions ensure that configuration values are not accidentally modified.
*/

func (c *Config) TestMode() bool           { return c.testMode }
func (c *Config) TestModeThreads() int     { return c.testModeThreads }
func (c *Config) TestModeUsers() int       { return c.testModeUsers }
func (c *Config) TestModeMinInterval() int { return c.testModeMinInterval }
func (c *Config) TestModeMaxInterval() int { return c.testModeMaxInterval }
func (c *Config) TestModeTimeOut() int     { return c.testModeTimeOut }
func (c *Config) TestModePostCap() int     { return c.testModePostCap }

func (c *Config) Master() bool      { return c.master }
func (c *Config) SaveConfig() bool  { return c.saveConfig }
func (c *Config) ConfigDir() string { return c.configDir }
func (c *Config) RPCPort() int      { return c.rpcPort }
func (c *Config) RPCRemAdr() string { return c.rpcRemAdr }

func (c *Config) CXOPort() int       { return c.cxoPort }
func (c *Config) CXORPCPort() int    { return c.cxoRPCPort }
func (c *Config) CXOUseMemory() bool { return c.cxoMemoryMode }
func (c *Config) CXODir() string     { return c.cxoDir }

func (c *Config) WebGUIEnable() bool      { return c.webGUIEnable }
func (c *Config) WebGUIPort() int         { return c.webGUIPort }
func (c *Config) WebGUIDir() string       { return c.webGUIDir }
func (c *Config) WebGUIOpenBrowser() bool { return c.webGUIOpenBrowser }
