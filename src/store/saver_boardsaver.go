package store

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
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
	Accepted      bool        `json:"accepted"`
	RejectedCount int         `json:"rejected_count"`
	Config        BoardConfig `json:"config"`
}

// BoardSaver manages boards.
type BoardSaver struct {
	sync.Mutex
	config *Config
	c      *CXO
	store  map[cipher.PubKey]*BoardInfo
	quit   chan struct{}
}

// NewBoardSaver creates a new BoardSaver.
func NewBoardSaver(config *Config, container *CXO) (*BoardSaver, error) {
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
	go bs.service()
	return &bs, nil
}

func (s *BoardSaver) Close() {
	for {
		select {
		case s.quit <- struct{}{}:
		default:
			return
		}
	}
}

func (s *BoardSaver) absConfigDir() string {
	return filepath.Join(s.config.ConfigDir, BoardSaverFileName)
}

// Helper function. Loads and checks boards' configuration file to memory.
func (s *BoardSaver) load() error {
	// Don't load if specified not to.
	if !s.config.MemoryMode {

		log.Println("[BOARDSAVER] Loading configuration file...")
		// Load boards from file.
		bcf := BoardSaverFile{}
		if e := file.LoadJSON(s.absConfigDir(), &bcf); e != nil {
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
			s.store[bpk] = &BoardInfo{Config: *bc}
			s.c.Subscribe(bc.Address, bpk)
			log.Println("\t\t loaded in memory")
		}
	}
	s.checkSynced()
	if s.config.Master {
		s.checkMasterURLs()
		s.checkMasterDeps()
	}
	return nil
}

// Helper function. Saves boards into configuration file.
func (s *BoardSaver) save() error {
	// Don't save if specified.
	if s.config.MemoryMode {
		return nil
	}
	// Load from memory.
	bcf := BoardSaverFile{}
	for _, bi := range s.store {
		bcf.Boards = append(bcf.Boards, &bi.Config)
	}
	return file.SaveJSON(s.absConfigDir(), bcf, os.FileMode(0700))
}

// Keeps imported threads synced.
func (s *BoardSaver) service() {
	log.Println("[BOARDSAVER] Sync service started.")
	msgs := s.c.GetUpdatesChan()
	for {
		select {
		case msg := <-msgs:
			switch msg.Mode() {
			case RootFilled:
				s.Lock()
				log.Printf("[BOARDSAVER] Checking dependencies for board '%s'", msg.PubKey().Hex())
				s.checkSingleDep(msg.PubKey())
				s.Unlock()
			case SubAccepted:
				s.Lock()
				if bi, got := s.store[msg.PubKey()]; got {
					bi.Accepted = true
				}
				s.Unlock()
			case SubRejected:
				s.Lock()
				if bi, got := s.store[msg.PubKey()]; got {
					bi.Accepted = false
					bi.RejectedCount += 1
				}
				s.Unlock()
			//case ConnCreated:
			//case ConnClosed:
				//s.Lock()
				//for bpk, bi := range s.store {
				//	addr := msg.Conn().Address()
				//	if bi.Config.Address == addr {
				//		s.c.Subscribe(addr, bpk)
				//	}
				//}
				//s.Unlock()
			}
		case <-s.quit:
			return
		}
	}
}

// Helper function. Check's the URL's of the boards which this node is master over.
func (s *BoardSaver) checkMasterURLs() {
	//for bpk, bi := range s.store {
	//if bi.Config.Master {
	//	b, e := s.c.GetBoard(bpk)
	//	if e != nil {
	//		continue
	//	}
	//if b.URL != s.config.RPCRemAdr() {
	//	s.c.ChangeBoardURL(bpk, bi.Config.GetSK(), s.config.RPCRemAdr())
	//}
	//}
	//}
}

// Helper function. Checks whether boards are synced.
func (s *BoardSaver) checkSynced() {
	for _, bi := range s.store {
		bi.Accepted = false
	}
	feeds := s.c.Feeds()
	for _, f := range feeds {
		bi, has := s.store[f]
		if has {
			bi.Accepted = true
		}
	}
	log.Printf("[BOARDSAVER] Accepted boards: (%d/%d)\n", len(feeds), len(s.store))
}

// Helper function. Checks whether single dependency is valid.
func (s *BoardSaver) checkSingleDep(bpkDep cipher.PubKey) {
	for _, bi := range s.store {
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
					return
				}
				// Sync.
				go func() {
					e = s.c.ImportThread(bpkDep, bi.Config.GetPK(), bi.Config.GetSK(), tRef)
					if e != nil {
						log.Println("[BOARDSAVER] sync failed for thread of reference:", t)
						log.Println("\t- cause:", e)
					}
				}()
			}
		}
	}
}

// Helper function. Checks whether dependencies are valid.
func (s *BoardSaver) checkMasterDeps() {
	for _, bi := range s.store {
		for j, dep := range bi.Config.Deps {
			fromBpk, e := misc.GetPubKey(dep.Board)
			if e != nil {
				log.Println("[BOARDSAVER] 'checkMasterDeps()' error:", e)
				log.Println("[BOARDSAVER] removing all dependencies of board with public key:", dep.Board)
				bi.Config.Deps = append(bi.Config.Deps[:j], bi.Config.Deps[j+1:]...)
				continue
			}
			// Subscribe internally.
			s.c.Subscribe("", fromBpk)

			for _, t := range dep.Threads {
				tRef, e := misc.GetReference(t)
				if e != nil {
					log.Println("[BOARDSAVER] 'checkMasterDeps()' error:", e)
					log.Println("[BOARDSAVER] removing thread dependency of reference:", t)
					bi.Config.RemoveDep(fromBpk, tRef)
					continue
				}
				// Sync.
				e = s.c.ImportThread(fromBpk, bi.Config.GetPK(), bi.Config.GetSK(), tRef)
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
func (s *BoardSaver) List() []BoardInfo {
	s.Lock()
	defer s.Unlock()
	//s.checkSynced()
	list, i := make([]BoardInfo, len(s.store)), 0
	for _, bi := range s.store {
		list[i] = *bi
		i += 1
	}
	return list
}

// ListKeys returns list of public keys of boards we are subscribed to.
func (s *BoardSaver) ListKeys() []cipher.PubKey {
	s.Lock()
	defer s.Unlock()
	keys, i := make([]cipher.PubKey, len(s.store)), 0
	for k := range s.store {
		keys[i] = k
		i += 1
	}
	return keys
}

func (s *BoardSaver) Check(bpk cipher.PubKey) error {
	s.Lock()
	defer s.Unlock()
	bi, has := s.store[bpk]
	if !has {
		return errors.Errorf("not subscribed to board '%s'", bpk.Hex())
	}
	addr := bi.Config.GetAddress()
	if !s.c.HasConnection(addr) {
		return s.c.Subscribe(addr, bpk)
	}
	return nil
}

// Get gets a subscription of specified board.
func (s *BoardSaver) Get(bpk cipher.PubKey) (BoardInfo, bool) {
	s.Lock()
	defer s.Unlock()
	bi, has := s.store[bpk]
	if has == false {
		return BoardInfo{}, has
	}
	addr := bi.Config.GetAddress()
	if !s.c.HasConnection(addr) {
		s.c.Subscribe(addr, bpk)
	}
	return *bi, has
}

// GetOfAddress obtains board information of boards of specified address.
func (s *BoardSaver) GetOfAddress(addr string) []BoardInfo {
	boards := []BoardInfo{}
	for _, bi := range s.store {
		if bi.Config.Address == addr {
			boards = append(boards, *bi)
		}
	}
	return boards
}

// Add adds a board to configuration.
func (s *BoardSaver) Add(addr string, bpk cipher.PubKey) {
	s.Lock()
	defer s.Unlock()

	if _, has := s.store[bpk]; has {
		return
	}

	bc := BoardConfig{Master: false, Address: addr, PubKey: bpk.Hex()}
	s.c.Subscribe(addr, bpk)

	s.store[bpk] = &BoardInfo{Config: bc}
	s.save()
}

// Remove removes a board from configuration.
func (s *BoardSaver) Remove(bpk cipher.PubKey) {
	s.Lock()
	defer s.Unlock()
	s.c.Unsubscribe("", bpk)
	delete(s.store, bpk)
	s.save()
}

// MasterAdd adds a board to configuration as master. Returns error if not master.
func (s *BoardSaver) MasterAdd(bpk cipher.PubKey, bsk cipher.SecKey) error {
	if s.config.Master == false {
		return errors.New("bbs node is not in master mode")
	}
	bc := BoardConfig{Master: true, PubKey: bpk.Hex(), SecKey: bsk.Hex()}
	_, e := bc.Check()
	if e != nil {
		return e
	}
	s.c.Subscribe("", bpk)

	s.Lock()
	defer s.Unlock()
	s.store[bpk] = &BoardInfo{Accepted: true, Config: bc}
	s.save()
	return nil
}

// AddBoardDep adds a dependency to a board.
func (s *BoardSaver) AddBoardDep(bpk, depBpk cipher.PubKey, deptRef skyobject.Reference) error {
	s.Lock()
	defer s.Unlock()
	// Check if we are subscribed to board of `depBpk`.
	if _, has := s.store[depBpk]; !has {
		return errors.New("failed to add board dependency: not subscribed to board " + depBpk.Hex())
	}
	// Retrieve board info for board of `bpk`.
	var bi *BoardInfo
	var has bool
	if bi, has = s.store[bpk]; !has {
		return errors.New("failed to add board dependency: not subscribed to board " + bpk.Hex())
	}
	// Add dependency.
	if e := bi.Config.AddDep(depBpk, deptRef); e != nil {
		return errors.Wrap(e, "failed to add board dependency")
	}
	s.save()
	return nil
}

// RemoveBoardDep removes a dependency from a board.
func (s *BoardSaver) RemoveBoardDep(bpk, depBpk cipher.PubKey, deptRef skyobject.Reference) error {
	s.Lock()
	defer s.Unlock()
	// Retrieve board info for board of `bpk`.
	bi, has := s.store[bpk]
	if !has {
		return errors.New("not subscribed of board")
	}
	// Remove dependency.
	if e := bi.Config.RemoveDep(depBpk, deptRef); e != nil {
		return e
	}
	s.save()
	return nil
}
