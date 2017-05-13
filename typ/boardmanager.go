package typ

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
	"sync"
)

const BoardsConfigFileName = "bbs_boards.json"

// BoardsConfig represents a configuration file for boards.
type BoardsConfig []*BoardConfig

// BoardManager manages board configurations.
type BoardManager struct {
	sync.Mutex
	Master  bool
	Masters []cipher.PubKey
	Boards  map[cipher.PubKey]*BoardConfig
	Loaded  map[cipher.PubKey]bool
	RPCAddr string
}

// NewBoardManager creates a new empty BoardManager.
func NewBoardManager(master bool, rpcAddr string) *BoardManager {
	bm := BoardManager{
		Master:  master,
		Boards:  make(map[cipher.PubKey]*BoardConfig),
		Loaded:  make(map[cipher.PubKey]bool),
		RPCAddr: rpcAddr,
	}
	return &bm
}

// Load loads data from config file.
func (m *BoardManager) Load() error {
	// Load from file.
	bc := BoardsConfig{}
	if e := util.LoadJSON(BoardsConfigFileName, &bc); e != nil {
		return e
	}
	// Check loaded boards. Append to m.Boards appropriately.
	for _, b := range bc {
		if e := b.PrepAndCheck(); e != nil {
			return e
		}
		m.Boards[b.PubKey] = b
		if m.Master && b.Master {
			b.URL = m.RPCAddr
			m.Masters = append(m.Masters, b.PubKey)
		}
	}
	if m.Master && len(m.Masters) == 0 {
		return errors.New("no masters")
	}
	return nil
}

// Save saves current in-memory config to file.
func (m *BoardManager) Save() error {
	bc := make(BoardsConfig, len(m.Boards))
	i := 0
	for _, b := range m.Boards {
		bc[i] = b
		i++
	}
	if e := util.SaveJSON(BoardsConfigFileName, bc, 0600); e != nil {
		return e
	}
	return nil
}

// Clear clears the in-memory config.
func (m *BoardManager) Clear() {
	m.Masters = []cipher.PubKey{}
	m.Boards = make(map[cipher.PubKey]*BoardConfig)
}

// AddConfig adds a new BoardConfig.
func (m *BoardManager) AddConfig(bc *BoardConfig) error {
	m.Lock()
	defer m.Unlock()

	pk, e := cipher.PubKeyFromHex(bc.PubKeyStr)
	if e != nil {
		return e
	}
	if _, has := m.Boards[pk]; has == true {
		return errors.New("config already exists")
	}
	m.Boards[pk] = bc
	m.Save()
	return nil
}

// RemoveConfig removes a BoardConfig.
func (m *BoardManager) RemoveConfig(pk cipher.PubKey) {
	m.Lock()
	defer m.Unlock()
	delete(m.Boards, pk)
	m.Save()
}

// HasConfig checks whether we have specified BoardConfig.
func (m *BoardManager) HasConfig(pk cipher.PubKey) bool {
	m.Lock()
	defer m.Unlock()

	_, has := m.Boards[pk]
	return has
}

// GetConfig gets a BoardConfig from given Public Key.
func (m *BoardManager) GetConfig(pk cipher.PubKey) (*BoardConfig, error) {
	m.Lock()
	defer m.Unlock()

	bc, has := m.Boards[pk]
	var e error
	if has == false {
		e = errors.New("config does not exist")
	}
	return bc, e
}

// GetList gets a list of BoardConfigs.
func (m *BoardManager) GetList() []*BoardConfig {
	m.Lock()
	defer m.Unlock()

	list := []*BoardConfig{}
	for _, bc := range m.Boards {
		list = append(list, bc)
	}
	return list
}

// GetCount gets the number of Boards.
func (m *BoardManager) GetCount() int {
	m.Lock()
	defer m.Unlock()

	return len(m.Boards)
}
