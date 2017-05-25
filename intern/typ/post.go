package typ

import (
	"errors"
	"github.com/evanlinjin/bbs/misc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"math"
	"strings"
	"time"
)

// Post represents a post as stored in cxo.
type Post struct {
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Author    string     `json:"author"`
	Created   int64      `json:"created"`
	Signature cipher.Sig `json:"-"`
	Ref       string     `json:"hash" enc:"-"`
}

func (p *Post) checkContent() error {
	title := strings.TrimSpace(p.Title)
	body := strings.TrimSpace(p.Body)
	if len(title) < 3 && len(body) < 3 {
		return errors.New("post content too short")
	}
	return nil
}

func (p *Post) checkCreated() error {
	dif := time.Now().UnixNano() - p.Created
	if dif < 0 {
		return errors.New("invalid timestamp")
	}
	if dif > int64(3*math.Pow(10, 11)) {
		return errors.New("invalid timestamp")
	}
	return nil
}

func (p *Post) checkAuthor() (cipher.PubKey, error) {
	if p.Author == (cipher.PubKey{}.Hex()) {
		return cipher.PubKey{}, errors.New("empty author public key")
	}
	return misc.GetPubKey(p.Author)
}

// Sign checks and signs the post.
func (p *Post) Sign(pk cipher.PubKey, sk cipher.SecKey) error {
	if e := p.checkContent(); e != nil {
		return e
	}
	p.Author = pk.Hex()
	p.Created = time.Now().UnixNano()
	p.Signature = cipher.Sig{}
	p.Signature = cipher.SignHash(cipher.SumSHA256(encoder.Serialize(*p)), sk)
	return nil
}

func (p Post) Verify() error {
	// Check title and body.
	if e := p.checkContent(); e != nil {
		return e
	}
	// Check author.
	authorPK, e := p.checkAuthor()
	if e != nil {
		return e
	}
	// Check creation time.
	if e := p.checkCreated(); e != nil {
		return e
	}
	// Check signature.
	sig := p.Signature
	p.Signature = cipher.Sig{}

	return cipher.VerifySignature(
		authorPK, sig,
		cipher.SumSHA256(encoder.Serialize(p)))
}

func (p *Post) Deserialize(data []byte) error {
	return encoder.DeserializeRaw(data, p)
}
