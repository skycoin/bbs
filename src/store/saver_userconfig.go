package store

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/skycoin/src/cipher"
)

// UserConfig represents the config of a user.
type UserConfig struct {
	Alias  string `json:"alias,omitempty"`
	Master bool   `json:"master"`
	PubKey string `json:"public_key"`
	SecKey string `json:"secret_key,omitempty"`
}

func (uc *UserConfig) GetPK() cipher.PubKey {
	pk, _ := misc.GetPubKey(uc.PubKey)
	return pk
}

func (uc *UserConfig) GetSK() cipher.SecKey {
	sk, _ := misc.GetSecKey(uc.SecKey)
	return sk
}

// Check checks the validity of the UserConfig.
func (uc *UserConfig) Check() (cipher.PubKey, error) {
	pk, e := misc.GetPubKey(uc.PubKey)
	if e != nil {
		return pk, e
	}
	if uc.Master {
		sk, e := misc.GetSecKey(uc.SecKey)
		if e != nil {
			return pk, e
		}
		if pk != cipher.PubKeyFromSecKey(sk) {
			return pk, errors.New("invalid public-secret pair")
		}
	}
	return pk, nil
}
