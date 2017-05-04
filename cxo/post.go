package cxo

import (
	"github.com/skycoin/skycoin/src/cipher"
)

// Post represents a post.
type Post struct {
	ID        []byte
	ThreadID  []byte
	Signature []byte
	Publisher cipher.PubKey
	Title     string
	Body      string
}
