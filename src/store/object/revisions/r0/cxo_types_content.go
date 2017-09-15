package r0

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"github.com/skycoin/bbs/src/store/object/transfer"
	"github.com/skycoin/bbs/src/misc/keys"
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
		log.Println("error getting content data:", e)
	}
	return out
}

func SetData(c Content, v *ContentData) {
	data, _ := json.Marshal(v)
	c.SetRaw(data)
}

type Board struct {
	R       cipher.PubKey `enc:"-" json:"-"`
	Data    []byte
	Created int64
}

func (b *Board) GetRaw() []byte  { return b.Data }
func (b *Board) SetRaw(v []byte) { b.Data = v }

func (b *Board) ToRep() (*transfer.BoardRep, error) {
	data := GetData(b)
	out := &transfer.BoardRep{
		Name: data.Name,
		Body: data.Body,
		Created: b.Created,
		Tags: nil,
	}
	return out, nil
}

func (b *Board) FromRep(bRep *transfer.BoardRep) error {
	SetData(b, &ContentData{
		Name: bRep.Name,
		Body: bRep.Body,
	})
	b.Created = bRep.Created
	return nil
}

type Thread struct {
	R       cipher.SHA256 `enc:"-" json:"-"`
	OfBoard cipher.PubKey `json:",string"`
	Data    []byte
	Created int64         `verify:"time"`
	Creator cipher.PubKey `verify:"upk"`
	Sig     cipher.Sig    `verify:"sig"`
}

func (t Thread) Verify() error    { return tag.Verify(&t) }
func (t *Thread) GetRaw() []byte  { return t.Data }
func (t *Thread) SetRaw(v []byte) { t.Data = v }

func (t *Thread) ToRep() (*transfer.ThreadRep, error) {
	data := GetData(t)
	out := &transfer.ThreadRep{
		Name: data.Name,
		Body: data.Body,
		Created: t.Created,
		Creator: t.Creator.Hex(),
	}
	return out, nil
}

func (t *Thread) FromRep(tRep *transfer.ThreadRep) error {
	SetData(t, &ContentData{
		Name: tRep.Name,
		Body: tRep.Body,
	})
	t.Created = tRep.Created
	var e error
	t.Creator, e = keys.GetPubKey(tRep.Creator)
	return e
}

type Post struct {
	R        cipher.SHA256 `enc:"-" json:"-"`
	OfBoard  cipher.PubKey `json:",string"`
	OfThread cipher.SHA256 `json:",string"`
	OfPost   cipher.SHA256 `json:",string"` // Can be empty.
	Data     []byte
	Created  int64         `verify:"time"`
	Creator  cipher.PubKey `verify:"upk"`
	Sig      cipher.Sig    `verify:"sig"`
}

func GetPost(pElem *skyobject.RefsElem) (*Post, error) {
	pVal, e := pElem.Value()
	if e != nil {
		return nil, elemValueErr(e, pElem)
	}
	p, ok := pVal.(*Post)
	if !ok {
		return nil, elemExtErr(pElem)
	}
	p.R = pElem.Hash
	return p, nil
}

func (p Post) Verify() error    { return tag.Verify(&p) }
func (p *Post) GetRaw() []byte  { return p.Data }
func (p *Post) SetRaw(v []byte) { p.Data = v }

func (p *Post) ToRep() (*transfer.PostRep, error) {
	data := GetData(p)
	out := &transfer.PostRep{
		OfPost: p.OfPost.Hex(),
		Name: data.Name,
		Body: data.Body,
		Created: p.Created,
		Creator: p.Creator.Hex(),
	}
	return out, nil
}

func (p *Post) FromRep(pRep *transfer.PostRep) error {
	SetData(p, &ContentData{
		Name: pRep.Name,
		Body: pRep.Body,
	})
	p.Created = pRep.Created
	var e error
	if p.OfPost, e = keys.GetHash(pRep.OfPost); e != nil {
		return e
	}
	if p.Creator, e = keys.GetPubKey(pRep.Creator); e != nil {
		return e
	}
	return nil
}

const (
	UserVote = iota
	ThreadVote
	PostVote
	UnknownVoteType
)

var VoteString = [...]string{
	UserVote:        "User Vote",
	ThreadVote:      "Thread Vote",
	PostVote:        "Post Vote",
	UnknownVoteType: "Unknown Vote Type",
}

type Vote struct {
	OfBoard  cipher.PubKey `json:",string"`
	OfUser   cipher.PubKey `json:",string"`
	OfThread cipher.SHA256 `json:",string"`
	OfPost   cipher.SHA256 `json:",string"`

	Mode int8
	Tag  []byte

	Created int64         `verify:"time"`
	Creator cipher.PubKey `verify:"upk"`
	Sig     cipher.Sig    `verify:"sig"`
}

func GetVote(vElem *skyobject.RefsElem) (*Vote, error) {
	vVal, e := vElem.Value()
	if e != nil {
		return nil, elemValueErr(e, vElem)
	}
	v, ok := vVal.(*Vote)
	if !ok {
		return nil, elemExtErr(vElem)
	}
	return v, nil
}

func (v Vote) Verify() error { return tag.Verify(&v) }

func (v *Vote) GetType() int {
	if v.OfUser != (cipher.PubKey{}) {
		return UserVote
	}
	if v.OfThread != (cipher.SHA256{}) {
		return ThreadVote
	}
	if v.OfPost != (cipher.SHA256{}) {
		return PostVote
	}
	return UnknownVoteType
}
