package store

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// Board dependency.
type BoardDep struct {
	Board   string   `json:"board"`
	Threads []string `json:"threads"`
}

// BoardConfig represents the config of a board.
type BoardConfig struct {
	Master  bool        `json:"master"`
	Address string      `json:"address,omitempty"`
	PubKey  string      `json:"public_key"`
	SecKey  string      `json:"secret_key,omitempty"`
	Deps    []*BoardDep `json:"dependencies,omitempty"`
}

// Check checks the validity of the BoardConfig.
func (bc *BoardConfig) Check() (cipher.PubKey, error) {
	pk, e := misc.GetPubKey(bc.PubKey)
	if e != nil {
		return pk, e
	}
	if bc.Master {
		sk, e := misc.GetSecKey(bc.SecKey)
		if e != nil {
			return pk, e
		}
		if pk != cipher.PubKeyFromSecKey(sk) {
			return pk, errors.New("invalid public-secret pair")
		}
	}
	return pk, nil
}

func (bc *BoardConfig) GetAddress() string {
	return bc.Address
}

func (bc *BoardConfig) GetPK() cipher.PubKey {
	pk, _ := misc.GetPubKey(bc.PubKey)
	return pk
}

func (bc *BoardConfig) GetSK() cipher.SecKey {
	sk, _ := misc.GetSecKey(bc.SecKey)
	return sk
}

func (bc *BoardConfig) AddDep(bpk cipher.PubKey, tRef skyobject.Reference) error {
	if !bc.Master {
		return errors.New("bbs node is not master of board")
	}
	bpkStr, tRefStr := bpk.Hex(), tRef.String()
	for _, bd := range bc.Deps {
		if bd.Board != bpkStr {
			continue
		}
		for _, td := range bd.Threads {
			if td != tRefStr {
				continue
			}
			return nil
		}
		bd.Threads = append(bd.Threads, tRefStr)
		return nil
	}
	bc.Deps = append(bc.Deps, &BoardDep{
		Board:   bpkStr,
		Threads: []string{tRefStr},
	})
	return nil
}

func (bc *BoardConfig) RemoveDep(bpk cipher.PubKey, tRef skyobject.Reference) error {
	if !bc.Master {
		return errors.New("not master of board")
	}
	bpkStr, tRefStr := bpk.Hex(), tRef.String()
	for i := len(bc.Deps) - 1; i >= 0; i-- {
		bd := bc.Deps[i]
		if bd.Board != bpkStr {
			continue
		}
		for j := len(bd.Threads) - 1; j >= 0; j-- {
			if bd.Threads[j] != tRefStr {
				continue
			}
			bd.Threads = append(bd.Threads[:j], bd.Threads[j+1:]...)
		}
	}
	for i := len(bc.Deps) - 1; i >= 0; i-- {
		if len(bc.Deps[i].Threads) == 0 {
			bc.Deps = append(bc.Deps[:i], bc.Deps[i+1:]...)
			//return nil
		}
	}
	return nil
}

// GetBoardOfThreadDep returns the master board of thread dependency.
func (bc *BoardConfig) GetBoardOfThreadDep(tRef skyobject.Reference) (cipher.PubKey, bool) {
	for _, dep := range bc.Deps {
		for _, t := range dep.Threads {
			if t == tRef.String() {
				pk, _ := misc.GetPubKey(dep.Board)
				return pk, true
			}
		}
	}
	return cipher.PubKey{}, false
}

func (bc BoardConfig) String(indent bool) string {
	var data []byte
	switch indent {
	case true:
		data, _ = json.MarshalIndent(bc, "", "\t")
	case false:
		data, _ = json.Marshal(bc)
	}
	return string(data)
}
