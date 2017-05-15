package store

import (
	"errors"
	"github.com/evanlinjin/bbs/cmd"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"sync"
)

const BoardsConfigFileName = "bbs_boards.json"

// BoardsConfigFile represents the layout of a configuration file of boards.
type BoardsConfigFile struct {
	Boards []*BoardConfig `json:"boards"`
}

// BoardConfig represents the config of a board.
type BoardConfig struct {
	Master bool
	PubKey string `json:"public_key"`
	SecKey string `json:"secret_key"`
}

// Check checks the validity of the BoardConfig.
func (bc *BoardConfig) Check() (cipher.PubKey, error) {
	pk, e := cipher.PubKeyFromHex(bc.PubKey)
	if e != nil {
		return pk, e
	}
	if bc.Master {
		sk, e := cipher.SecKeyFromHex(bc.SecKey)
		if e != nil {
			return pk, e
		}
		if pk != cipher.PubKeyFromSecKey(sk) {
			return pk, errors.New("invalid public-secret pair")
		}
	}
	return pk, nil
}

// BoardInfo represents the board's information.
type BoardInfo struct {
	Synced      bool        `json:"synced"`
	BoardConfig BoardConfig `json:"config"`
}

// BoardSaver manages boards.
type BoardSaver struct {
	sync.Mutex
	config *cmd.Config
	c      *Container
	store  map[cipher.PubKey]*BoardInfo
}

// NewBoardSaver creates a new BoardSaver.
func NewBoardSaver(config *cmd.Config, container *Container) (*BoardSaver, error) {
	bs := BoardSaver{
		config: config,
		c:      container,
		store:  make(map[cipher.PubKey]*BoardInfo),
	}
	if e := bs.load(); e != nil {
		if e := bs.save(); e != nil {
			return nil, e
		}
	}
	return &bs, nil
}

// Helper function. Loads and checks boards' configuration file to memory.
func (bs *BoardSaver) load() error {
	log.Println("[BOARDSAVER] Loading configuration file...")
	// Load boards from file.
	bcf := BoardsConfigFile{}
	if e := util.LoadJSON(BoardsConfigFileName, &bcf); e != nil {
		return e
	}
	// Check loaded boards and store in memory.
	for _, bc := range bcf.Boards {
		log.Printf("\t%v", *bc)
		bpk, e := bc.Check()
		if e != nil {
			log.Println("\t\t config file check:", e)
			continue
		}
		bs.store[bpk] = &BoardInfo{BoardConfig: *bc}
		bs.c.Subscribe(bpk)
		log.Println("\t\t loaded in memory")
	}
	bs.checkSynced()
	return nil
}

// Helper function. Saves boards into configuration file.
func (bs *BoardSaver) save() error {
	// Load from memory.
	bcf := BoardsConfigFile{}
	for _, bi := range bs.store {
		bcf.Boards = append(bcf.Boards, &bi.BoardConfig)
	}
	return util.SaveJSON(BoardsConfigFileName, bcf, 0600)
}

// Helper function. Checks whether boards are synced.
func (bs *BoardSaver) checkSynced() {
	for _, bi := range bs.store {
		bi.Synced = false
	}
	feeds := bs.c.Feeds()
	for _, f := range feeds {
		bi, has := bs.store[f]
		if has == false {
			continue
		}
		bi.Synced = true
	}
	log.Printf("[BOARDSAVER] Synced boards: (%d/%d)\n", len(feeds), len(bs.store))
}

// List returns a list of boards that are in configuration.
func (bs *BoardSaver) List() []BoardInfo {
	bs.Lock()
	defer bs.Unlock()
	bs.checkSynced()
	list, i := make([]BoardInfo, len(bs.store)), 0
	for _, bi := range bs.store {
		list[i] = *bi
		i += 1
	}
	return list
}

// ListKeys returns list of public keys of boards we are subscribed to.
func (bs *BoardSaver) ListKeys() []cipher.PubKey {
	bs.Lock()
	defer bs.Unlock()
	keys, i := make([]cipher.PubKey, len(bs.store)), 0
	for k := range bs.store {
		keys[i] = k
		i += 1
	}
	return keys
}

// Remove removes a board from configuration.
func (bs *BoardSaver) Remove(bpk cipher.PubKey) {
	bs.Lock()
	defer bs.Unlock()
	delete(bs.store, bpk)
	bs.save()
}

// Add adds a board to configuration.
func (bs *BoardSaver) Add(bpk cipher.PubKey) {
	bc := BoardConfig{Master: false, PubKey: bpk.Hex()}
	bs.c.Subscribe(bpk)

	bs.Lock()
	defer bs.Unlock()
	bs.store[bpk] = &BoardInfo{BoardConfig: bc}
	bs.save()
}

// MasterAdd adds a board to configuration as master. Returns error if not master.
func (bs *BoardSaver) MasterAdd(bpk cipher.PubKey, bsk cipher.SecKey) error {
	if bs.config.Master() == false {
		return errors.New("bbs node is not in master mode")
	}
	bc := BoardConfig{Master: true, PubKey: bpk.Hex(), SecKey: bsk.Hex()}
	_, e := bc.Check()
	if e != nil {
		return e
	}
	bs.c.Subscribe(bpk)

	bs.Lock()
	defer bs.Unlock()
	bs.store[bpk] = &BoardInfo{BoardConfig: bc}
	bs.save()
	return nil
}

// Get gets a subscription of specified board.
func (bs *BoardSaver) Get(bpk cipher.PubKey) (BoardInfo, bool) {
	bs.Lock()
	defer bs.Unlock()
	bi, has := bs.store[bpk]
	if has == false {
		return BoardInfo{}, has
	}
	return *bi, has
}
