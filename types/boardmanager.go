package types

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

// BoardConfig represents a board's configuration as stored on a local file.
type BoardConfig struct {
	Master    bool   `json:"master"`
	PublicKey string `json:"public_key"`
	SecretKey string `json:"secret_key,omitempty"` // Empty if Master = false.
	URL       string `json:"url"`
}

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

	pk, e := cipher.PubKeyFromHex(bc.PublicKey)
	if e != nil {
		return e
	}
	if _, has := bm.Configs[pk]; has == true {
		return errors.New("config already exists")
	}
	bm.Configs[pk] = bc
	return nil
}

// NewMasterConfigFromSeed generates a new BoardConfig from a seed.
// This BoardConfig is one which is master.
func (bm *BoardManager) NewMasterConfigFromSeed(seed, URL string) (
	*BoardConfig, cipher.PubKey, cipher.SecKey, error,
) {
	if bm.Master == false {
		return nil, cipher.PubKey{}, cipher.SecKey{}, errors.New("not master")
	}
	pk, sk := cipher.GenerateDeterministicKeyPair([]byte(seed))
	bc := &BoardConfig{
		Master:    true,
		PublicKey: pk.Hex(),
		SecretKey: sk.Hex(),
		URL:       URL,
	}
	e := bm.AddConfig(bc)
	return bc, pk, sk, e
}

// RemoveConfig removes a BoardConfig.
func (bm *BoardManager) RemoveConfig(pk cipher.PubKey) {
	bm.Lock()
	defer bm.Unlock()

	delete(bm.Configs, pk)
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
