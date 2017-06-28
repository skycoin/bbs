package typ

import (
	"encoding/json"
	"github.com/pkg/errors"
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
	Meta    []byte `json:"-"`
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

// GetMeta obtains the meta data of the board.
func (b *Board) GetMeta() (*BoardMeta, error) {
	bm := new(BoardMeta)
	if e := json.Unmarshal(b.Meta, bm); e != nil {
		return nil, errors.Wrap(e, "failed to obtain board meta")
	}
	return bm, nil
}

// SetMeta set's the meta data of the board.
func (b *Board) SetMeta(bm *BoardMeta) error {
	data, e := json.Marshal(*bm)
	if e != nil {
		return errors.Wrap(e, "failed to set board meta")
	}
	b.Meta = data
	return nil
}
