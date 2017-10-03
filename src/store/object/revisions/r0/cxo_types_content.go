package r0

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

type BasicBody struct {
	OfBoard string `json:"of_board"`
	Creator string `json:"creator"`
}

func newBasicBody(raw []byte) (*BasicBody, error) {
	out := new(BasicBody)
	if e := json.Unmarshal(raw, out); e != nil {
		return nil, e
	}
	return out, nil
}

func (c *BasicBody) GetOfBoard() (cipher.PubKey, error) {
	return keys.GetPubKey(c.OfBoard)
}

func (c *BasicBody) GetCreator() (cipher.PubKey, error) {
	return keys.GetPubKey(c.Creator)
}

type Content struct {
	Header []byte // Contains type, creator public key and signature.
	Body   []byte // Contains actual content.
}

func (c *Content) String() string {
	data := struct {
		Header string
		Body   string
	}{
		Header: string(c.Header),
		Body:   string(c.Body),
	}
	raw, _ := json.MarshalIndent(data, "", "    ")
	return string(raw)
}

func (c *Content) Verify() (*BasicBody, error) {
	b, e := newBasicBody(c.Body)
	if e != nil {
		return nil, e
	}

	h := c.GetHeader()
	h.Hash = cipher.SumSHA256(c.Body).Hex()
	h.PK = b.Creator

	if c.Header, e = json.Marshal(h); e != nil {
		return nil, e
	}

	if e := cipher.VerifySignature(h.GetPK(), h.GetSig(), h.GetHash()); e != nil {
		return nil, e
	}

	return b, nil
}

func GetContentFromRef(cRef *skyobject.Ref) (*Content, error) {
	cVal, e := cRef.Value()
	if e != nil {
		return nil, valueErr(e, cRef)
	}
	c, ok := cVal.(*Content)
	if !ok {
		return nil, extErr(cRef)
	}
	return c, nil
}

func GetContentFromElem(cElem *skyobject.RefsElem) (*Content, error) {
	cVal, e := cElem.Value()
	if e != nil {
		return nil, elemValueErr(e, cElem)
	}
	c, ok := cVal.(*Content)
	if !ok {
		return nil, elemExtErr(cElem)
	}
	return c, nil
}

func (c *Content) GetHeader() *ContentHeaderData {
	data := new(ContentHeaderData)
	jsonUnmarshal(c.Header, data)
	return data
}

func (c *Content) SetHeader(data *ContentHeaderData) {
	c.Header = jsonMarshal(data)
}

func (c *Content) GetBody(v interface{}) {
	jsonUnmarshal(c.Body, v)
}

func (c *Content) SetBody(v interface{}) {
	c.Body = jsonMarshal(v)
}

func (c *Content) SetBodyRaw(raw []byte) {
	c.Body = raw
}

func (c *Content) ToBoard() *Board {
	return &Board{Content: c}
}

func (c *Content) ToThread() *Thread {
	return &Thread{Content: c}
}

func (c *Content) ToPost() *Post {
	return &Post{Content: c}
}

func (c *Content) ToThreadVote() *ThreadVote {
	return &ThreadVote{Content: c}
}

func (c *Content) ToPostVote() *PostVote {
	return &PostVote{Content: c}
}

func (c *Content) ToUserVote() *UserVote {
	return &UserVote{Content: c}
}

type ContentRep struct {
	Header *ContentHeaderData `json:"header,omitempty"`
	Body   interface{}        `json:"body,omitempty"`
	Votes  interface{}        `json:"votes,omitempty"`
}

type ContentType string

const (
	V5BoardType      = ContentType("5,board")
	V5ThreadType     = ContentType("5,thread")
	V5PostType       = ContentType("5,post")
	V5ThreadVoteType = ContentType("5,thread_vote")
	V5PostVoteType   = ContentType("5,post_vote")
	V5UserVoteType   = ContentType("5,user_vote")
)

type ContentHeaderData struct {
	Type ContentType `json:"type"` // Content type and version.
	Hash string      `json:"hash"` // Hash of body.
	PK   string      `json:"pk"`   // Public key.
	Sig  string      `json:"sig"`  // Signature of body.
}

func (h *ContentHeaderData) GetHash() cipher.SHA256 {
	out, e := keys.GetHash(h.Hash)
	if e != nil {
		log.Println("failed to get 'hash' from header:", e)
	}
	return out
}

func (h *ContentHeaderData) GetPK() cipher.PubKey {
	out, e := keys.GetPubKey(h.PK)
	if e != nil {
		log.Println("failed to get 'pk' from header:", e)
	}
	return out
}

func (h *ContentHeaderData) GetSig() cipher.Sig {
	out, e := keys.GetSig(h.Sig)
	if e != nil {
		log.Println("failed to get 'sig' from header:", e)
	}
	return out
}

func (h *ContentHeaderData) Verify() error {
	if e := cipher.VerifySignature(h.GetPK(), h.GetSig(), h.GetHash()); e != nil {
		return genErrHeaderUnverified(e, string(h.Type))
	}
	return nil
}

type Board struct {
	*Content
}

func (b *Board) GetBody() *BoardData {
	data := new(BoardData)
	b.Content.GetBody(data)
	return data
}

func (b *Board) ToRep() *ContentRep {
	return &ContentRep{
		Header: b.GetHeader(),
		Body:   b.GetBody(),
	}
}

func (b *Board) Fill(bpk cipher.PubKey, data *BoardData) {
	b.Content = new(Content)
	b.SetBody(data)
	b.SetHeader(&ContentHeaderData{
		Type: V5BoardType,
		Hash: cipher.SumSHA256(b.Body).Hex(),
		PK:   bpk.Hex(),
	})
}

type Thread struct {
	*Content
}

func (t *Thread) GetBody() *ThreadData {
	data := new(ThreadData)
	t.Content.GetBody(data)
	return data
}

func (t *Thread) ToRep() *ContentRep {
	return &ContentRep{
		Header: t.GetHeader(),
		Body:   t.GetBody(),
	}
}

func (t *Thread) Fill(tt *ThreadTransport) {
	t.Content = new(Content)
	t.SetBodyRaw(tt.raw)
	t.SetHeader(tt.header)
}

type Post struct {
	*Content
}

func GetPost(pElem *skyobject.RefsElem) (*Post, error) {
	pVal, e := pElem.Value()
	if e != nil {
		return nil, elemValueErr(e, pElem)
	}
	p, ok := pVal.(*Content)
	if !ok {
		return nil, elemExtErr(pElem)
	}
	return p.ToPost(), nil
}

func (p *Post) GetBody() *PostData {
	data := new(PostData)
	p.Content.GetBody(data)
	return data
}

func (p *Post) ToRep() *ContentRep {
	return &ContentRep{
		Header: p.GetHeader(),
		Body:   p.GetBody(),
	}
}

func (p *Post) Fill(pt *PostTransport) {
	p.Content = new(Content)
	p.SetBodyRaw(pt.raw)
	p.SetHeader(pt.header)
}

type ThreadVote struct {
	*Content
}

func (tv *ThreadVote) GetBody() *ThreadVoteData {
	data := new(ThreadVoteData)
	tv.Content.GetBody(data)
	return data
}

func (tv *ThreadVote) ToRep() *ContentRep {
	return &ContentRep{
		Header: tv.GetHeader(),
		Body:   tv.GetBody(),
	}
}

func (tv *ThreadVote) Fill(tvt *ThreadVoteTransport) {
	tv.Content = new(Content)
	tv.SetBodyRaw(tvt.raw)
	tv.SetHeader(tvt.header)
}

type PostVote struct {
	*Content
}

func (pv *PostVote) GetBody() *PostVoteData {
	data := new(PostVoteData)
	pv.Content.GetBody(data)
	return data
}

func (pv *PostVote) ToRep() *ContentRep {
	return &ContentRep{
		Header: pv.GetHeader(),
		Body:   pv.GetBody(),
	}
}

func (pv *PostVote) Fill(pvt *PostVoteTransport) {
	pv.Content = new(Content)
	pv.SetBodyRaw(pvt.raw)
	pv.SetHeader(pvt.header)
}

type UserVote struct {
	*Content
}

func (uv *UserVote) GetBody() *UserVoteData {
	data := new(UserVoteData)
	uv.Content.GetBody(data)
	return data
}

func (uv *UserVote) ToRep() *ContentRep {
	return &ContentRep{
		Header: uv.GetHeader(),
		Body:   uv.GetBody(),
	}
}

func (uv *UserVote) Fill(uvt *UserVoteTransport) {
	uv.Content = new(Content)
	uv.SetBodyRaw(uvt.raw)
	uv.SetHeader(uvt.header)
}

//func GetVote(vElem *skyobject.RefsElem) (*Vote, error) {
//	vVal, e := vElem.Value()
//	if e != nil {
//		return nil, elemValueErr(e, vElem)
//	}
//	v, ok := vVal.(*Vote)
//	if !ok {
//		return nil, elemExtErr(vElem)
//	}
//	return v, nil
//}

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
