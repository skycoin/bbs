package typ

import "github.com/skycoin/skycoin/src/cipher"

// Post represents a post as stored in cxo.
type Post struct {
	Title      string        `json:"title"`
	Body       string        `json:"body"`
	Creator    cipher.PubKey `json:"-"`
	CreatorStr string        `json:"creator" enc:"-"`
	Created    int64         `json:"created"`
	Signature  cipher.Sig    `json:"-"`
}
