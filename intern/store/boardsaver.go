package store

import (
	"github.com/evanlinjin/bbs/cmd/bbsnode/args"
	"github.com/evanlinjin/bbs/intern/cxo"
	"github.com/evanlinjin/bbs/misc"
	"github.com/pkg/errors"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"os"
	"path/filepath"
	"sync"
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
	config *args.Config
	c      *cxo.Container
	store  map[cipher.PubKey]*BoardInfo
	quit   chan struct{}
}

// NewBoardSaver creates a new BoardSaver.
func NewBoardSaver(config *args.Config, container *cxo.Container) (*BoardSaver, error) {
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

func (bs *BoardSaver) absConfigDir() string {
	return filepath.Join(bs.config.ConfigDir(), BoardSaverFileName)
}

// Helper function. Loads and checks boards' configuration file to memory.
func (bs *BoardSaver) load() error {
	// Don't load if specified not to.
	if bs.config.SaveConfig() {

		log.Println("[BOARDSAVER] Loading configuration file...")
		// Load boards from file.
		bcf := BoardSaverFile{}
		if e := util.LoadJSON(bs.absConfigDir(), &bcf); e != nil {
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
	}
	bs.checkSynced()
	if bs.config.Master() {
		bs.checkURLs()
		bs.checkDeps()
		go bs.service()
	}
	return nil
}

// Helper function. Saves boards into configuration file.
func (bs *BoardSaver) save() error {
	// Don't save if specified.
	if !bs.config.SaveConfig() {
		return nil
	}
	// Load from memory.
	bcf := BoardSaverFile{}
	for _, bi := range bs.store {
		bcf.Boards = append(bcf.Boards, &bi.Config)
	}
	return util.SaveJSON(bs.absConfigDir(), bcf, os.FileMode(0700))
}

// Keeps imported threads synced.
func (bs *BoardSaver) service() {
	log.Println("[BOARDSAVER] Sync service started.")
	msgs := bs.c.GetUpdatesChan()
	for {
		select {
		case msg := <-msgs:
			switch msg.Mode() {
			case cxo.RootFilled:
				bs.Lock()
				log.Printf("[BOARDSAVER] Checking dependencies for board '%s'", msg.PubKey().Hex())
				bs.checkSingleDep(msg.PubKey())
				bs.Unlock()
			case cxo.FeedAdded:
				bs.Lock()
				if bi, got := bs.store[msg.PubKey()]; got {
					bi.Synced = true
				}
				bs.Unlock()
			case cxo.FeedDeleted:
				bs.Lock()
				if bi, got := bs.store[msg.PubKey()]; got {
					bi.Synced = false
				}
				bs.Unlock()
			}
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

// Helper function. Checks whether single dependency is valid.
func (bs *BoardSaver) checkSingleDep(bpkDep cipher.PubKey) {
	for _, bi := range bs.store {
		for _, dep := range bi.Config.Deps {
			if dep.Board != bpkDep.Hex() {
				continue
			}
			for _, t := range dep.Threads {
				tRef, e := misc.GetReference(t)
				if e != nil {
					log.Println("[BOARDSAVER] 'checkSingleDep()' error:", e)
					log.Println("[BOARDSAVER] removing thread dependency of reference:", t)
					bi.Config.RemoveDep(bpkDep, tRef)
					continue
				}
				// Sync.
				e = bs.c.ImportThread(bpkDep, bi.Config.GetPK(), bi.Config.GetSK(), tRef)
				if e != nil {
					log.Println("[BOARDSAVER] sync failed for thread of reference:", t)
					log.Println("\t- cause:", e)
					continue
				}
			}
			// List of dependencies are of unique dependencies. No duplicates so continue.
			continue
		}
	}
}

// Helper function. Checks whether dependencies are valid.
func (bs *BoardSaver) checkDeps() {
	for _, bi := range bs.store {
		for j, dep := range bi.Config.Deps {
			fromBpk, e := misc.GetPubKey(dep.Board)
			if e != nil {
				log.Println("[BOARDSAVER] 'checkDeps()' error:", e)
				log.Println("[BOARDSAVER] removing all dependencies of board with public key:", dep.Board)
				bi.Config.Deps = append(bi.Config.Deps[:j], bi.Config.Deps[j+1:]...)
				continue
			}
			// Subscribe internally.
			bs.c.Subscribe(fromBpk)

			for _, t := range dep.Threads {
				tRef, e := misc.GetReference(t)
				if e != nil {
					log.Println("[BOARDSAVER] 'checkDeps()' error:", e)
					log.Println("[BOARDSAVER] removing thread dependency of reference:", t)
					bi.Config.RemoveDep(fromBpk, tRef)
					continue
				}
				// Sync.
				e = bs.c.ImportThread(fromBpk, bi.Config.GetPK(), bi.Config.GetSK(), tRef)
				if e != nil {
					log.Println("[BOARDSAVER] sync failed for thread of reference:", t)
					log.Println("\t- cause:", e)
					return
				}
			}
		}
	}
}

// List returns a list of boards that are in configuration.
func (bs *BoardSaver) List() []BoardInfo {
	bs.Lock()
	defer bs.Unlock()
	//bs.checkSynced()
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
		return errors.New("failed to add board dependency: not subscribed to board " + depBpk.Hex())
	}
	// Retrieve board info for board of `bpk`.
	var bi *BoardInfo
	var has bool
	if bi, has = bs.store[bpk]; !has {
		return errors.New("failed to add board dependency: not subscribed to board " + bpk.Hex())
	}
	// Add dependency.
	if e := bi.Config.AddDep(depBpk, deptRef); e != nil {
		return errors.Wrap(e, "failed to add board dependency")
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
