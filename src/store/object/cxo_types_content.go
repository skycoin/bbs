package object

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/skycoin/src/cipher"
)

type Content struct {
	OfBoard   cipher.PubKey
	OfContent []cipher.SHA256 // THREAD:len(0), POST:len(1+).
	// If POST: [0](thread hash), [1](optional)(post hash).

	Data []byte

	Created uint64        `verify:"time"`
	Creator cipher.PubKey `verify:"upk"`
	Sig     cipher.Sig    `verify:"sig"`
}

func ToContent(v interface{}) *Content {
	if content, ok := v.(*Content); ok {
		return content
	}
	return nil
}

func (c Content) Verify() error { return tag.Verify(&c) }

func (c *Content) SetData(v interface{}) error {
	data, e := json.Marshal(v)
	if e != nil {
		return e
	}
	c.Data = data
	return nil
}

func (c *Content) GetData(v interface{}) error {
	return json.Unmarshal(c.Data, v)
}

func (c *Content) IsThread() bool {
	return len(c.OfContent) == 0
}

func (c *Content) IsPost() bool {
	return len(c.OfContent) != 0
}

func (c *Content) RefersPost() (cipher.SHA256, bool) {
	return cipher.SHA256{}, false
}

type ThreadData struct {
	Name string `json:"name"`
	Desc string `json:"description"`
}

type PostData struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}
