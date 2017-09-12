package msgs

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

const (
	tLen  = 1
	pkLen = 33
)

type Message []byte

func NewMessage(data []byte) (*Message, error) {
	if len(data) < tLen+pkLen {
		return nil, boo.New(boo.InvalidRead,
			"received data is of invalid length")
	}
	if e := cipher.NewPubKey(data[tLen : tLen+pkLen]).Verify(); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"received data has invalid public key")
	}
	msg := Message(data)
	return &msg, nil
}

func (in Message) GetMsgType() MsgType {
	return MsgType(in[0])
}

func (in Message) GetOutPK() cipher.PubKey {
	return cipher.NewPubKey(in[tLen : tLen+pkLen])
}

func (in Message) GetData() []byte {
	return in[tLen+pkLen:]
}

func (in Message) ExtractThread() (*r0.Thread, error) {
	v := new(r0.Thread)
	if e := encoder.DeserializeRaw(in.GetData(), v); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"failed to deserialize thread data")
	}
	if e := v.Verify(); e != nil {
		return nil, e
	}
	return v, nil
}

func (in Message) ExtractPost() (*r0.Post, error) {
	v := new(r0.Post)
	if e := encoder.DeserializeRaw(in.GetData(), v); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"failed to deserialize post data")
	}
	if e := v.Verify(); e != nil {
		return nil, e
	}
	return v, nil
}

func (in Message) ExtractVote() (*r0.Vote, error) {
	v := new(r0.Vote)
	if e := encoder.DeserializeRaw(in.GetData(), v); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"failed to deserialize vote data")
	}
	if e := v.Verify(); e != nil {
		return nil, e
	}
	return v, nil
}

func (in Message) ExtractResponse() (*Response, error) {
	v := new(Response)
	if e := encoder.DeserializeRaw(in.GetData(), v); e != nil {
		return nil, e
	}
	return v, nil
}

type Response struct {
	Hash   cipher.SHA256 // Hash of entire message.
	Type   MsgType       // Message type.
	Seq    uint64        // Root sequence that satisfies.
	Okay   bool          // Whether successful.
	ErrTyp int           // Type of error.
	ErrMsg string        // Message of error.
}

func NewResponse(msg *Message, goal uint64, e error) *Response {
	r := &Response{
		Hash: cipher.SumSHA256(msg.GetData()),
		Type: msg.GetMsgType(),
		Seq:  goal,
		Okay: true,
	}
	if e != nil {
		r.Okay = false
		r.ErrTyp = boo.Type(e)
		r.ErrMsg = e.Error()
	}
	return r
}
