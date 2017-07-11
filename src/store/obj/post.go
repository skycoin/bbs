package obj

import "github.com/skycoin/skycoin/src/cipher"

type Post struct {
	Title   string        `json:"title"`
	Body    string        `json:"body"`
	Created int64         `json:"created"`
	User    cipher.PubKey `json:"-"`
	Sig     cipher.Sig    `json:"-"`
	Meta    []byte        `json:"-"`
}
