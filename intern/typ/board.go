package typ

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
	"strings"
	"time"
)

var (
	ErrInvalidName = errors.New("invalid name provided")
)

// Board represents a board stored in cxo.
type Board struct {
	Name    string `json:"name"`
	Desc    string `json:"description"`
	PubKey  string `json:"public_key"`
	URL     string `json:"url"`
	Created int64  `json:"created"`
}

// Check checks the validity of the board, outputs it's public key and updates it's fields.
func (b *Board) Check() (cipher.PubKey, error) {
	bpk, e := cipher.PubKeyFromHex(b.PubKey)
	if e != nil {
		return bpk, e
	}
	b.Name = strings.TrimSpace(b.Name)
	b.Desc = strings.TrimSpace(b.Desc)
	if len(b.Name) < 3 {
		return bpk, ErrInvalidName
	}
	b.Created = time.Now().UnixNano()
	return bpk, nil
}

// TouchWithSeed updates the board's public key, returns the key pair, and updates the board's time.
func (b *Board) TouchWithSeed(seed []byte) (cipher.PubKey, cipher.SecKey) {
	bpk, bsk := cipher.GenerateDeterministicKeyPair(seed)
	b.PubKey = bpk.Hex()
	b.Created = time.Now().UnixNano()
	return bpk, bsk
}
