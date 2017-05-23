package store

import (
	"errors"
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/intern/cxo"
	"github.com/evanlinjin/bbs/misc"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"sync"
	"time"
)

const BoardSaverFileName = "bbs_boards.json"

// BoardSaverFile represents the layout of a configuration file of boards.
type BoardSaverFile struct {
	Boards []*BoardConfig `json:"boards"`
}

// BoardInfo represents the board's information.
type BoardInfo struct {
	Synced bool        `json:"synced"`
	Config BoardConfig `json:"config"`
}

// BoardSaver manages boards.
type BoardSaver struct {
	sync.Mutex
	config *cmd.Config
	c      *cxo.Container
	store  map[cipher.PubKey]*BoardInfo
	quit   chan struct{}
}

// NewBoardSaver creates a new BoardSaver.
func NewBoardSaver(config *cmd.Config, container *cxo.Container) (*BoardSaver, error) {
	bs := BoardSaver{
		config: config,
		c:      container,
		store:  make(map[cipher.PubKey]*BoardInfo),
		quit:   make(chan struct{}),
	}
	bs.load()
	if e := bs.save(); e != nil {
		return nil, e
	}
	return &bs, nil
}

func (bs *BoardSaver) Close() {
	for {
		select {
		case bs.quit <- struct{}{}:
		default:
			return
		}
	}
}

// Helper function. Loads and checks boards' configuration file to memory.
func (bs *BoardSaver) load() error {
	log.Println("[BOARDSAVER] Loading configuration file...")
	// Load boards from file.
	bcf := BoardSaverFile{}
	if e := util.LoadJSON(BoardSaverFileName, &bcf); e != nil {
		log.Println("[BOARDSAVER]", e)
	}
	// Check loaded boards and intern in memory.
	for _, bc := range bcf.Boards {
		log.Printf("\t- %v (master: %v)", bc.PubKey, bc.Master)
		bpk, e := bc.Check()
		if e != nil {
			log.Println("\t\t config file check:", e)
			continue
		}
		bs.store[bpk] = &BoardInfo{Config: *bc}
		bs.c.Subscribe(bpk)
		log.Println("\t\t loaded in memory")
	}
	bs.checkSynced()
	if bs.config.Master() {
		bs.checkURLs()
		bs.checkDeps()
		go bs.service()
	}
	return nil
}

// Keeps imported threads synced.
func (bs *BoardSaver) service() {
	log.Println("[BOARDSAVER] Sync service started.")
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			bs.Lock()
			bs.checkDeps()
			bs.Unlock()

		case <-bs.quit:
			return
		}
	}
}

// Helper function. Check's the URL's of the boards which this node is master over.
func (bs *BoardSaver) checkURLs() {
	for bpk, bi := range bs.store {
		if bi.Config.Master {
			b, e := bs.c.GetBoard(bpk)
			if e != nil {
				continue
			}
			if b.URL != bs.config.RPCServerRemAdr() {
				bs.c.ChangeBoardURL(bpk, bi.Config.GetSK(), bs.config.RPCServerRemAdr())
			}
		}
	}
}

// Helper function. Checks whether boards are synced.
func (bs *BoardSaver) checkSynced() {
	for _, bi := range bs.store {
		bi.Synced = false
	}
	feeds := bs.c.Feeds()
	for _, f := range feeds {
		bi, has := bs.store[f]
		if has {
			bi.Synced = true
		}
	}
	log.Printf("[BOARDSAVER] Synced boards: (%d/%d)\n", len(feeds), len(bs.store))
}

// Helper function. Checks whether dependencies are valid.
// TODO: Implement.
func (bs *BoardSaver) checkDeps() {
	for _, bi := range bs.store {
		for j, dep := range bi.Config.Deps {
			fromBpk, e := misc.GetPubKey(dep.Board)
			if e != nil {
				log.Println("[BOARDSAVER] 'checkDeps()' error:", e)
				log.Println("[BOARDSAVER] removing all dependencies of board with public key:", dep.Board)
				bi.Config.Deps = append(bi.Config.Deps[:j], bi.Config.Deps[j+1:]...)
				return
			}
			// Subscribe internally.
			bs.c.Subscribe(fromBpk)

			for _, t := range dep.Threads {
				tRef, e := misc.GetReference(t)
				if e != nil {
					log.Println("[BOARDSAVER] 'checkDeps()' error:", e)
					log.Println("[BOARDSAVER] removing thread dependency of reference:", t)
					bi.Config.RemoveDep(fromBpk, tRef)
					return
				}
				// Sync.
				e = bs.c.ImportThread(fromBpk, bi.Config.GetPK(), bi.Config.GetSK(), tRef)
				if e != nil {
					log.Println("[BOARDSAVER] 'checkDeps()' error:", e)
					log.Println("[BOARDSAVER] sync failed for thread of reference:", t)
					log.Println("[BOARDSAVER] removing thread dependency of reference:", t)
					bi.Config.RemoveDep(fromBpk, tRef)
					return
				}
			}
		}
	}
}

// Helper function. Saves boards into configuration file.
func (bs *BoardSaver) save() error {
	// Load from memory.
	bcf := BoardSaverFile{}
	for _, bi := range bs.store {
		bcf.Boards = append(bcf.Boards, &bi.Config)
	}
	return util.SaveJSON(BoardSaverFileName, bcf, 0600)
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

// Add adds a board to configuration.
func (bs *BoardSaver) Add(bpk cipher.PubKey) {
	bc := BoardConfig{Master: false, PubKey: bpk.Hex()}
	bs.c.Subscribe(bpk)

	bs.Lock()
	defer bs.Unlock()
	bs.store[bpk] = &BoardInfo{Config: bc}
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
	bs.store[bpk] = &BoardInfo{Config: bc}
	bs.save()
	return nil
}

// AddBoardDep adds a dependency to a board.
func (bs *BoardSaver) AddBoardDep(bpk, depBpk cipher.PubKey, deptRef skyobject.Reference) error {
	bs.Lock()
	defer bs.Unlock()
	// Check if we are subscribed to board of `depBpk`.
	if _, has := bs.store[depBpk]; !has {
		return errors.New("not subscribed to board")
	}
	// Retrieve board info for board of `bpk`.
	var bi *BoardInfo
	var has bool
	if bi, has = bs.store[bpk]; !has {
		return errors.New("not subscribed of board")
	}
	// Add dependency.
	if e := bi.Config.AddDep(depBpk, deptRef); e != nil {
		return e
	}
	bs.save()
	return nil
}

// RemoveBoardDep removes a dependency from a board.
func (bs *BoardSaver) RemoveBoardDep(bpk, depBpk cipher.PubKey, deptRef skyobject.Reference) error {
	bs.Lock()
	defer bs.Unlock()
	// Retrieve board info for board of `bpk`.
	bi, has := bs.store[bpk]
	if !has {
		return errors.New("not subscribed of board")
	}
	// Remove dependency.
	if e := bi.Config.RemoveDep(depBpk, deptRef); e != nil {
		return e
	}
	bs.save()
	return nil
}

// Remove removes a board from configuration.
func (bs *BoardSaver) Remove(bpk cipher.PubKey) {
	bs.Lock()
	defer bs.Unlock()
	delete(bs.store, bpk)
	bs.save()
}
