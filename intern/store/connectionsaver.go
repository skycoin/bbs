package store

import (
	"github.com/skycoin/bbs/cmd/bbsnode/args"
	"github.com/skycoin/bbs/intern/cxo"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"os"
	"path/filepath"
	"sync"
)

const ConnectionSaverFileName = "bbs_connections.json"

type ConnectionSaver struct {
	sync.Mutex
	config *args.Config
	c      *cxo.Container
	list   map[string]bool
}

func NewConnectionSaver(config *args.Config, container *cxo.Container) (*ConnectionSaver, error) {
	cs := &ConnectionSaver{
		config: config,
		c:      container,
		list:   make(map[string]bool),
	}
	cs.load()
	cs.checkConnections()
	if e := cs.save(); e != nil {
		return nil, e
	}
	return cs, nil
}

func (cs *ConnectionSaver) absConfigDir() string {
	return filepath.Join(cs.config.ConfigDir(), ConnectionSaverFileName)
}

func (cs *ConnectionSaver) load() {
	// Don't load if specified not to.
	if !cs.config.SaveConfig() {
		return
	}
	log.Println("[CONNECTIONSAVER] Loading configuration file...")
	// Load connections from file.
	if e := util.LoadJSON(cs.absConfigDir(), &cs.list); e != nil {
		log.Println("[CONNECTIONSAVER] Error:", e)
	}
}

func (cs *ConnectionSaver) checkConnections() error {
	for addr, got := range cs.list {
		log.Printf("\t- %s (%v)", addr, got)
		if e := cs.c.Connect(addr); e != nil {
			log.Printf("\t\t- error: %s", e.Error())
			continue
		}
	}
	return nil
}

func (cs *ConnectionSaver) List() []string {
	cs.Lock()
	defer cs.Unlock()

	out, i := make([]string, len(cs.list)), 0
	for addr := range cs.list {
		out[i] = addr
		i += 1
	}
	return out
}

func (cs *ConnectionSaver) Add(addr string) error {
	cs.Lock()
	defer cs.Unlock()

	if _, got := cs.list[addr]; got {
		return nil
	}
	if e := cs.c.Connect(addr); e != nil {
		return e
	}
	cs.list[addr] = true
	cs.save()
	return nil
}

func (cs *ConnectionSaver) Remove(addr string) error {
	cs.Lock()
	defer cs.Unlock()

	cs.c.Disconnect(addr)
	delete(cs.list, addr)
	cs.save()
	return nil
}

func (cs *ConnectionSaver) save() error {
	// Don't save if specified not to.
	if !cs.config.SaveConfig() {
		return nil
	}
	return util.SaveJSON(cs.absConfigDir(), cs.list, os.FileMode(0700))
}
