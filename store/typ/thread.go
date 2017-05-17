package typ

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strings"
)

// Thread represents a thread stored in cxo.
type Thread struct {
	Name        string `json:"name"`
	Desc        string `json:"description"`
	MasterBoard string `json:"master_board"`
	Hash        string `json:"hash" enc:"-"`
}

func (t *Thread) Check() error {
	if t == nil {
		return errors.New("thread is nil")
	}
	name := strings.TrimSpace(t.Name)
	desc := strings.TrimSpace(t.Desc)
	if len(name) < 3 && len(desc) < 3 {
		return errors.New("thread content too short")
	}
	return nil
}

func (t Thread) Sign(sk cipher.SecKey) cipher.Sig {
	return cipher.SignHash(
		cipher.SumSHA256(encoder.Serialize(t)), sk)
}

func (t Thread) Verify(pk cipher.PubKey, sig cipher.Sig) error {
	if e := t.Check(); e != nil {
		return e
	}
	return cipher.VerifySignature(
		pk, sig, cipher.SumSHA256(encoder.Serialize(t)))
}
