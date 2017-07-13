package typ

import (
	"encoding/hex"
	"errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strings"
	"time"
)

// Thread represents a thread stored in cxo.
type Thread struct {
	Name        string     `json:"name"`         //
	Desc        string     `json:"description"`  //
	Author      string     `json:"author"`       // TODO: Make use of.
	Created     int64      `json:"created"`      // TODO: Make use of.
	Sig         cipher.Sig `json:"-"`            // TODO: Make use of.
	Ref         string     `json:"ref" enc:"-"`  //
	Meta        []byte     `json:"-"`            // TODO: Make use of.
	MasterBoard string     `json:"master_board"` //
}

func (t *Thread) checkContent() error {
	name := strings.TrimSpace(t.Name)
	desc := strings.TrimSpace(t.Desc)
	if len(name) < 3 || len(desc) < 3 {
		return errors.New("post content too short")
	}
	return nil
}

func (t *Thread) checkAuthor() (cipher.PubKey, error) {
	if t.Author == (cipher.PubKey{}.Hex()) {
		return cipher.PubKey{}, errors.New("empty author public key")
	}
	return misc.GetPubKey(t.Author)
}

func (t *Thread) Sign(pk cipher.PubKey, sk cipher.SecKey) error {
	if e := t.checkContent(); e != nil {
		return e
	}
	t.Author = pk.Hex()
	t.Created = 0
	t.Ref = ""
	t.Sig = cipher.Sig{}
	t.Sig = cipher.SignHash(cipher.SumSHA256(encoder.Serialize(*t)), sk)
	return nil
}

func (t Thread) Verify() error {
	if e := t.checkContent(); e != nil {
		return e
	}
	authorPK, e := t.checkAuthor()
	if e != nil {
		return e
	}
	sig := t.Sig
	t.Sig = cipher.Sig{}
	t.Created = 0
	t.Ref = ""

	return cipher.VerifySignature(
		authorPK, sig,
		cipher.SumSHA256(encoder.Serialize(t)))
}

func (t Thread) GetRef() skyobject.Reference {
	ref, _ := misc.GetReference(t.Ref)
	return ref
}

func (t Thread) SerializeToHex() string {
	t.Ref = ""
	return hex.EncodeToString(encoder.Serialize(t))
}

func (t *Thread) Touch() {
	t.Created = time.Now().UnixNano()
}

func (t *Thread) Deserialize(data []byte) error {
	t.Ref = ""
	return encoder.DeserializeRaw(data, t)
}
