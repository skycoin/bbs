package typ

import (
	"errors"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"strings"
)

// Thread represents a thread stored in cxo.
type Thread struct {
	Name string `json:"name"`
	Desc string `json:"description"`
	Hash string `json:"hash" enc:"-"`
}

func InitThread(tRef skyobject.Reference) *Thread {
	return &Thread{
		Hash: cipher.SHA256(tRef).Hex(),
	}
}

func (t *Thread) CheckAndPrep() error {
	if t == nil {
		return errors.New("nil thread")
	}
	t.Name = strings.TrimSpace(t.Name)
	t.Desc = strings.TrimSpace(t.Desc)
	if len(t.Name) == 0 {
		return errors.New("invalid thread name")
	}
	return nil
}
