package object

import (
	"encoding/json"
	"fmt"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
)

func errGetFromBody(e error, what string) error {
	return boo.WrapTypef(e, boo.InvalidRead,
		"failed to get '%s' from body", what)
}

const (
	TrustTag = "trust"
	SpamTag  = "spam"
	BlockTag = "block"
)

type ImageData struct {
	Name   string       `json:"name"`
	Hash   string       `json:"hash"`
	URL    string       `json:"url,omitempty"`
	Size   int          `json:"size,omitempty"`
	Height int          `json:"height,omitempty"`
	Width  int          `json:"width,omitempty"`
	Thumbs []*ImageData `json:"thumbs,omitempty"`
}

type Body struct {
	Type     ContentType       `json:"type"`                      // ALL
	TS       int64             `json:"ts"`                        // ALL
	OfBoard  string            `json:"of_board,omitempty"`        // thread, post, thread_vote, post_vote, user_vote
	OfThread string            `json:"of_thread,omitempty"`       // post, thread_vote
	OfPost   string            `json:"of_post,omitempty"`         // post (optional), post_vote
	OfUser   string            `json:"of_user,omitempty"`         // vote
	Name     string            `json:"name,omitempty"`            // board, thread, post
	Body     string            `json:"body,omitempty"`            // board, thread, post
	Images   []*ImageData      `json:"images,omitempty"`          // post (optional)
	Value    int               `json:"value,omitempty"`           // thread_vote, post_vote, user_vote
	Tags     []string          `json:"tags,omitempty"`            // board, thread_vote, post_vote, user_vote
	SubKeys  []MessengerSubKey `json:"submission_keys,omitempty"` // board
	Creator  string            `json:"creator,omitempty"`         // thread, post, thread_vote, post_vote, user_vote
}

func NewBody(raw []byte) (*Body, error) {
	out := new(Body)
	if e := json.Unmarshal(raw, out); e != nil {
		return nil, e
	}
	return out, nil
}

func (c *Body) GetOfBoard() (cipher.PubKey, error) {
	if pk, e := tag.GetPubKey(c.OfBoard); e != nil {
		return pk, errGetFromBody(e, "of_board")
	} else {
		return pk, nil
	}
}

func (c *Body) GetOfThread() (cipher.SHA256, error) {
	if hash, e := tag.GetHash(c.OfThread); e != nil {
		return hash, errGetFromBody(e, "of_thread")
	} else {
		return hash, nil
	}
}

func (c *Body) GetOfPost() (cipher.SHA256, error) {
	if hash, e := tag.GetHash(c.OfPost); e != nil {
		return hash, errGetFromBody(e, "of_post")
	} else {
		return hash, nil
	}
}

func (c *Body) GetOfUser() (cipher.PubKey, error) {
	if pk, e := tag.GetPubKey(c.OfUser); e != nil {
		return pk, errGetFromBody(e, "of_user")
	} else {
		return pk, nil
	}
}

func (c *Body) GetSubKeys() []*MessengerSubKeyTransport {
	out := make([]*MessengerSubKeyTransport, len(c.SubKeys))
	for i, subKey := range c.SubKeys {
		var e error
		if out[i], e = subKey.ToTransport(); e != nil {
			log.Printf("failed to obtain 'submission_keys[%d]' with error: %v", i, e)
		}
	}
	return out
}

func (c *Body) SetSubKeys(subKeys []*MessengerSubKeyTransport) {
	c.SubKeys = make([]MessengerSubKey, len(subKeys))
	for i, subKey := range subKeys {
		c.SubKeys[i] = subKey.ToMessengerSubKey()
	}
}

func (c *Body) GetCreator() (cipher.PubKey, error) {
	if pk, e := tag.GetPubKey(c.Creator); e != nil {
		return pk, errGetFromBody(e, "creator")
	} else {
		return pk, nil
	}
}

func (c *Body) HasTag(tag string) bool {
	for _, v := range c.Tags {
		if v == tag {
			return true
		}
	}
	return false
}

func (c *Body) HasValue(v int) bool {
	return c.Value == v
}

type Content struct {
	Header []byte `json:"header"` // Contains type, creator public key and signature.
	Body   []byte `json:"body"`   // Contains actual content.
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

func (c *Content) Verify() (*Body, error) {
	b, e := NewBody(c.Body)
	if e != nil {
		return nil, e
	}

	h := c.GetHeader()
	h.Hash = cipher.SumSHA256(c.Body).Hex()

	creator, e := b.GetCreator()
	if e != nil {
		return nil, e
	}

	if c.Header, e = json.Marshal(h); e != nil {
		return nil, e
	}

	if e := cipher.VerifySignature(creator, h.GetSig(), h.GetHash()); e != nil {
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

func (c *Content) GetBody() *Body {
	body := new(Body)
	jsonUnmarshal(c.Body, body)
	return body
}

func (c *Content) SetBody(body *Body) {
	c.Body = jsonMarshal(body)
}

func (c *Content) SetBodyRaw(raw []byte) {
	c.Body = raw
}

func (c *Content) ToRep() *ContentRep {
	return &ContentRep{
		Header: c.GetHeader(),
		Body:   c.GetBody(),
	}
}

type ContentRep struct {
	PubKey string             `json:"public_key,omitempty"`
	Header *ContentHeaderData `json:"header,omitempty"`
	Body   interface{}        `json:"body,omitempty"`
	Votes  interface{}        `json:"votes,omitempty"`
}

type ContentType string

func (t *ContentType) IsValid() bool {
	switch *t {
	case
		V5BoardType,
		V5ThreadType,
		V5PostType,
		V5ThreadVoteType,
		V5PostVoteType,
		V5UserVoteType:
		return true
	}
	return false
}

const (
	V5BoardType      = ContentType("5,board")
	V5ThreadType     = ContentType("5,thread")
	V5PostType       = ContentType("5,post")
	V5ThreadVoteType = ContentType("5,thread_vote")
	V5PostVoteType   = ContentType("5,post_vote")
	V5UserVoteType   = ContentType("5,user_vote")
)

type ContentHeaderData struct {
	Hash string `json:"hash,omitempty"` // Hash of body.
	Sig  string `json:"sig,omitempty"`  // Signature of body.
}

func (h *ContentHeaderData) GetHash() cipher.SHA256 {
	out, e := tag.GetHash(h.Hash)
	if e != nil {
		log.Println("failed to get 'hash' from header:", e)
	}
	return out
}

func (h *ContentHeaderData) GetSig() cipher.Sig {
	out, e := tag.GetSig(h.Sig)
	if e != nil {
		log.Println("failed to get 'sig' from header:", e)
	}
	return out
}

func (h *ContentHeaderData) Verify(upk cipher.PubKey) error {
	if e := cipher.VerifySignature(upk, h.GetSig(), h.GetHash()); e != nil {
		return genErrHeaderUnverified(e, h.Hash)
	}
	return nil
}

/*
	<<< TRANSPORT >>>
*/

type Transport struct {
	Header  *ContentHeaderData
	Body    *Body
	Content *Content
}

func NewTransport(rawBody []byte, sig cipher.Sig) (*Transport, error) {
	out := new(Transport)

	var e error
	if out.Body, e = NewBody(rawBody); e != nil {
		return nil, e
	}

	creator, e := out.Body.GetCreator()
	if e != nil {
		return nil, e
	}

	out.Header = &ContentHeaderData{
		Hash: cipher.SumSHA256(rawBody).Hex(),
		Sig:  sig.Hex(),
	}
	if e := out.Header.Verify(creator); e != nil {
		return nil, e
	}

	out.Content = new(Content)
	if out.Content.Header, e = json.Marshal(out.Header); e != nil {
		return nil, e
	}
	out.Content.Body = rawBody

	return out, nil
}

func (t *Transport) GetOfBoard() cipher.PubKey {
	pk, _ := t.Body.GetOfBoard()
	return pk
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

func genErrInvalidJSON(e error, what string) error {
	return boo.WrapType(e, boo.InvalidInput,
		fmt.Sprintf("failed to read '%s' data", what))
}

func genErrHeaderUnverified(e error, hash string) error {
	return boo.WrapType(e, boo.NotAuthorised,
		fmt.Sprintf("failed to verify content of hash '%s'", hash))
}
