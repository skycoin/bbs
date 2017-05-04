package datastore

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
)

type BoardManager struct {
	Configs map[cipher.PubKey]*BoardConfig
}

func NewBoardManager() *BoardManager {
	bm := BoardManager{
		Configs: make(map[cipher.PubKey]*BoardConfig),
	}
	return &bm
}

func (bm *BoardManager) AddConfig(bc *BoardConfig) error {
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

func (bm *BoardManager) NewMasterConfigFromSeed(seed, URL string) (
	*BoardConfig, cipher.PubKey, cipher.SecKey, error,
) {
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

func (bm *BoardManager) RemoveConfig(pk cipher.PubKey) {
	delete(bm.Configs, pk)
}

func (bm *BoardManager) GetList() []*BoardConfig {
	list := []*BoardConfig{}
	for _, bc := range bm.Configs {
		list = append(list, bc)
	}
	return list
}

func (bm *BoardManager) GetConfig(pk cipher.PubKey) (*BoardConfig, error) {
	bc, has := bm.Configs[pk]
	var e error
	if has == false {
		e = errors.New("config does not exist")
	}
	return bc, e
}