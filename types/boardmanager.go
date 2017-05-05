package types

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

// BoardManager manages board configurations.
type BoardManager struct {
	sync.Mutex
	Master  bool
	Configs map[cipher.PubKey]*BoardConfig
}

// NewBoardManager creates a new empty BoardManager.
func NewBoardManager(master bool) *BoardManager {
	bm := BoardManager{
		Master:  master,
		Configs: make(map[cipher.PubKey]*BoardConfig),
	}
	return &bm
}

// AddConfig adds a new BoardConfig.
func (bm *BoardManager) AddConfig(bc *BoardConfig) error {
	bm.Lock()
	defer bm.Unlock()

	pk, e := cipher.PubKeyFromHex(bc.PublicKeyStr)
	if e != nil {
		return e
	}
	if _, has := bm.Configs[pk]; has == true {
		return errors.New("config already exists")
	}
	bm.Configs[pk] = bc
	return nil
}

// RemoveConfig removes a BoardConfig.
func (bm *BoardManager) RemoveConfig(pk cipher.PubKey) {
	bm.Lock()
	defer bm.Unlock()

	delete(bm.Configs, pk)
}

// HasConfig checks whether we have specified BoardConfig.
func (bm *BoardManager) HasConfig(pk cipher.PubKey) bool {
	bm.Lock()
	defer bm.Unlock()

	_, has := bm.Configs[pk]
	return has
}

// GetList gets a list of BoardConfigs.
func (bm *BoardManager) GetList() []*BoardConfig {
	bm.Lock()
	defer bm.Unlock()

	list := []*BoardConfig{}
	for _, bc := range bm.Configs {
		list = append(list, bc)
	}
	return list
}

// GetConfig gets a BoardConfig from given Public Key.
func (bm *BoardManager) GetConfig(pk cipher.PubKey) (*BoardConfig, error) {
	bm.Lock()
	defer bm.Unlock()

	bc, has := bm.Configs[pk]
	var e error
	if has == false {
		e = errors.New("config does not exist")
	}
	return bc, e
}
