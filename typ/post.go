package typ

import "github.com/skycoin/skycoin/src/cipher"

// Post represents a post as stored in cxo.
type Post struct {
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Creator   string     `json:"creator"`
	Created   int64      `json:"created"`
	Signature cipher.Sig `json:"-"`
}
