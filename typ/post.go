package typ

import (
	"errors"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"strings"
	"time"
)

// Post represents a post as stored in cxo.
type Post struct {
	Title      string        `json:"title"`
	Body       string        `json:"body"`
	Creator    cipher.PubKey `json:"-"`
	CreatorStr string        `json:"creator" enc:"-"`
	Created    int64         `json:"created"`
	Signature  cipher.Sig    `json:"-"`
}

// CheckContent checks the post's title and body.
func (p *Post) CheckContent() (e error) {
	if p == nil {
		return errors.New("nil post")
	}
	// Check title and body of post.
	if len(strings.TrimSpace(p.Title)) == 0 {
		return errors.New("invalid post title")
	}
	if len(strings.TrimSpace(p.Title)) == 0 {
		return errors.New("invalid post body")
	}
	return
}

// CheckCreator checks the posts's creator.
func (p *Post) CheckCreator() (e error) {
	// Check creator's public key.
	// Whether be it that it came from cxo or json data.
	if p.Creator == (cipher.PubKey{}) {
		p.Creator, e = GetPubKey(p.CreatorStr)
	} else {
		p.CreatorStr = p.Creator.Hex()
	}
	return
}

// Sign signs the post before putting in cxo.
func (p *Post) Sign(sk cipher.SecKey) cipher.Sig {
	hash := cipher.SumSHA256(encoder.Serialize(*p))
	p.Signature = cipher.SignHash(hash, sk)
	return p.Signature
}

// CheckSig checks the signature of the post.
func (p Post) CheckSig() error {
	// Extract signature.
	sig := p.Signature
	p.Signature = cipher.Sig{}

	// Obtain hash.
	hash := cipher.SumSHA256(encoder.Serialize(p))

	// Verify signature.
	return cipher.VerifySignature(p.Creator, sig, hash)
}

// Touch sets the created time of post to now.
func (p *Post) Touch() {
	p.Created = time.Now().UnixNano()
}
