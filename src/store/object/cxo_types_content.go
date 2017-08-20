package object

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

type Content interface {
	GetRaw() []byte
	SetRaw(v []byte)
}

func ToContent(v interface{}) Content {
	return v.(Content)
}

func GetData(c Content) *ContentData {
	out := new(ContentData)
	if e := json.Unmarshal(c.GetRaw(), out); e != nil {
		log.Println("Error getting content data: ", e)
	}
	return out
}

func SetData(c Content, v *ContentData) {
	data, _ := json.Marshal(v)
	c.SetRaw(data)
}

type ContentData struct {
	Heading      string   `json:"heading"`
	Body         string   `json:"body"`
	SubAddresses []string `json:"submission_addresses,omitempty"`
}

type Board struct {
	R       cipher.PubKey `enc:"-"`
	Data    []byte
	Created int64
}

func (b *Board) GetRaw() []byte  { return b.Data }
func (b *Board) SetRaw(v []byte) { b.Data = v }

type Thread struct {
	R       cipher.SHA256 `enc:"-"`
	OfBoard cipher.PubKey
	Data    []byte
	Created int64         `verify:"time"`
	Creator cipher.PubKey `verify:"upk"`
	Sig     cipher.Sig    `verify:"sig"`
}

func (t Thread) Verify() error    { return tag.Verify(&t) }
func (t *Thread) GetRaw() []byte  { return t.Data }
func (t *Thread) SetRaw(v []byte) { t.Data = v }

type Post struct {
	R        cipher.SHA256 `enc:"-"`
	OfBoard  cipher.PubKey
	OfThread cipher.SHA256
	OfPost   cipher.SHA256 // Can be empty.
	Data     []byte
	Created  int64         `verify:"time"`
	Creator  cipher.PubKey `verify:"upk"`
	Sig      cipher.Sig    `verify:"sig"`
}

func (p Post) Verify() error    { return tag.Verify(&p) }
func (p *Post) GetRaw() []byte  { return p.Data }
func (p *Post) SetRaw(v []byte) { p.Data = v }

type Vote struct {
	OfUser   cipher.PubKey
	OfThread cipher.SHA256
	OfPost   cipher.SHA256

	Mode int8
	Tag  []byte

	Created int64         `verify:"time"`
	Creator cipher.PubKey `verify:"upk"`
	Sig     cipher.Sig    `verify:"sig"`
}

func ToVote(v interface{}) *Vote {
	if vote, ok := v.(*Vote); ok {
		return vote
	}
	return nil
}

func (v Vote) Verify() error { return tag.Verify(&v) }
