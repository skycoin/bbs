package typ

import (
	"encoding/hex"
	"errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strings"
)

// Thread represents a thread stored in cxo.
type Thread struct {
	Name        string `json:"name"`
	Desc        string `json:"description"`
	MasterBoard string `json:"master_board"`
	Ref         string `json:"ref" enc:"-"`
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
	t.Ref = ""
	return cipher.SignHash(
		cipher.SumSHA256(encoder.Serialize(t)), sk)
}

func (t Thread) Verify(pk cipher.PubKey, sig cipher.Sig) error {
	if e := t.Check(); e != nil {
		return e
	}
	t.Ref = ""
	return cipher.VerifySignature(
		pk, sig, cipher.SumSHA256(encoder.Serialize(t)))
}

func (t Thread) GetRef() skyobject.Reference {
	ref, _ := misc.GetReference(t.Ref)
	return ref
}

func (t Thread) SerializeToHex() string {
	t.Ref = ""
	return hex.EncodeToString(encoder.Serialize(t))
}

func (t *Thread) Deserialize(data []byte) error {
	t.Ref = ""
	return encoder.DeserializeRaw(data, t)
}
