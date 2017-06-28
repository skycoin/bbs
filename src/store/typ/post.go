package typ

import (
	"errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
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
	Ref       string     `json:"ref" enc:"-"`
}

func (p *Post) checkContent() error {
	title := strings.TrimSpace(p.Title)
	body := strings.TrimSpace(p.Body)
	if len(title) < 3 && len(body) < 3 {
		return errors.New("post content too short")
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
	p.Created = 0
	p.Signature = cipher.Sig{}
	p.Signature = cipher.SignHash(cipher.SumSHA256(encoder.Serialize(*p)), sk)
	return nil
}

// Verify checks the legitimacy of the post.
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
	// Check signature.
	sig := p.Signature
	p.Signature = cipher.Sig{}
	p.Created = 0

	return cipher.VerifySignature(
		authorPK, sig,
		cipher.SumSHA256(encoder.Serialize(p)))
}

// Touch updates the timestamp of Post.
func (p *Post) Touch() {
	p.Created = time.Now().UnixNano()
}

func (p *Post) Deserialize(data []byte) error {
	return encoder.DeserializeRaw(data, p)
}
