package r0

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/bbs/src/store/object/transfer"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

type Board struct {
	R   cipher.PubKey `enc:"-" json:"-"`
	Raw []byte
}

func (b *Board) GetData() *BoardData {
	data := new(BoardData)
	jsonUnmarshal(b.Raw, data)
	return data
}

func (b *Board) SetData(data *BoardData) {
	b.Raw = jsonMarshal(data)
}

func (b *Board) ToRep() (*transfer.BoardRep, error) {
	data := b.GetData()
	return &transfer.BoardRep{
		Name:    data.Name,
		Body:    data.Body,
		Created: data.Created,
		Tags:    data.Tags,
	}, nil
}

func (b *Board) FromRep(bRep *transfer.BoardRep) error {
	b.SetData(&BoardData{
		Name:    bRep.Name,
		Body:    bRep.Body,
		Created: bRep.Created,
		Tags:    bRep.Tags,
	})
	return nil
}

type Thread struct {
	R   cipher.SHA256 `enc:"-" json:"-"`
	Raw []byte
	Sig cipher.Sig `verify:"sig"`
}

func (t Thread) Verify(creator cipher.PubKey) error {
	return tag.Verify(&t, creator)
}

func (t *Thread) GetData() *ThreadData {
	data := new(ThreadData)
	jsonUnmarshal(t.Raw, data)
	return data
}

func (t *Thread) SetData(data *ThreadData) {
	t.Raw = jsonMarshal(data)
}

func (t *Thread) ToRep() (*transfer.ThreadRep, error) {
	data := t.GetData()
	return &transfer.ThreadRep{
		Name:    data.Name,
		Body:    data.Body,
		Created: data.Created,
		Creator: data.Creator,
	}, nil
}

func (t *Thread) FromRep(tRep *transfer.ThreadRep) error {
	t.SetData(&ThreadData{
		Name:    tRep.Name,
		Body:    tRep.Body,
		Creator: tRep.Creator,
		Created: tRep.Created,
	})
	return nil
}

type Post struct {
	R   cipher.SHA256 `enc:"-" json:"-"`
	Raw []byte
	Sig cipher.Sig `verify:"sig"`
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

func (p Post) Verify(creator cipher.PubKey) error {
	return tag.Verify(&p, creator)
}

func (p *Post) GetData() *PostData {
	data := new(PostData)
	jsonUnmarshal(p.Raw, data)
	return data
}

func (p *Post) SetData(data *PostData) {
	p.Raw = jsonMarshal(data)
}

func (p *Post) ToRep() (*transfer.PostRep, error) {
	data := p.GetData()
	return &transfer.PostRep{
		OfPost:  data.OfPost,
		Name:    data.Name,
		Body:    data.Body,
		Created: data.Created,
		Creator: data.Creator,
	}, nil
}

func (p *Post) FromRep(pRep *transfer.PostRep) error {
	p.SetData(&PostData{
		OfPost:  pRep.OfPost,
		Name:    pRep.Name,
		Body:    pRep.Body,
		Created: pRep.Created,
		Creator: pRep.Creator,
	})
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

/*
	<<< BOARD SUMMARY WRAP >>>
*/

type BoardSummaryWrap struct {
	PubKey cipher.PubKey `verify:"upk"`
	Raw    []byte
	Sig    cipher.Sig `verify:"sig"`
}

func (bsw *BoardSummaryWrap) Sign(pk cipher.PubKey, sk cipher.SecKey) {
	tag.Sign(bsw, pk, sk)
}

func (bsw BoardSummaryWrap) Verify() error {
	return tag.Verify(&bsw)
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func jsonUnmarshal(raw []byte, v interface{}) {
	if e := json.Unmarshal(raw, v); e != nil {
		log.Println("json unmarshal error:", e)
	}
}

func jsonMarshal(v interface{}) []byte {
	data, e := json.Marshal(v)
	if e != nil {
		log.Println("json marshal error:", e)
	}
	return data
}
