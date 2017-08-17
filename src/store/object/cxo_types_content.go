package object

import (
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/skycoin/src/cipher"
)

type Content struct {
	OfBoard   cipher.PubKey
	OfContent []cipher.SHA256 // THREAD:len(0), POST:len(1+).
	// If POST: [0](thread hash), [1](optional)(post hash).

	Data []byte
	data interface{} `enc:"-"`

	Created uint64        `verify:"time"`
	Creator cipher.PubKey `verify:"upk"`
	Sig     cipher.Sig    `verify:"sig"`
}

func (c Content) Verify() error { return tag.Verify(&c) }

type ThreadData struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

type PostData struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
