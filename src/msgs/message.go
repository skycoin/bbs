package msgs

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/bbs/src/store/object"
)

const (
	tLen = 1
)

type Message []byte

func (m Message) Check() error {
	if len(m) < factory.MSG_HEADER_END {
		return boo.New(boo.InvalidRead,
			"received data is of invalid length")
	}
	return nil
}

func (m Message) GetOP() uint {
	return uint(m[0])
}

func (m Message) ToSendMessage() (*SendMessage, error) {
	if m.GetOP() != factory.OP_SEND || len(m) < factory.SEND_MSG_META_END {
		return nil, boo.New(boo.InvalidRead, "not send message")
	}
	sm := SendMessage(m)
	if e := sm.GetFromPubKey().Verify(); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"send message has invalid from public key")
	}
	if e := sm.GetToPubKey().Verify(); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"send message has invalid to public key")
	}
	return &sm, nil
}

type SendMessage Message

func (m SendMessage) GetFromPubKey() cipher.PubKey {
	return cipher.NewPubKey(
		m[factory.SEND_MSG_PUBLIC_KEY_BEGIN:factory.SEND_MSG_PUBLIC_KEY_END])
}

func (m SendMessage) GetToPubKey() cipher.PubKey {
	return cipher.NewPubKey(
		m[factory.SEND_MSG_TO_PUBLIC_KEY_BEGIN:factory.SEND_MSG_TO_PUBLIC_KEY_END])
}

func (m SendMessage) GetBody() []byte {
	return m[factory.SEND_MSG_TO_PUBLIC_KEY_END:]
}

func (m SendMessage) ToBBSMessage() (*BBSMessage, error) {
	if len(m) < factory.SEND_MSG_META_END+tLen {
		return nil, boo.New(boo.InvalidRead, "not bbs message")
	}
	bm := BBSMessage{SendMessage: m}
	if int(bm.GetMsgType()) >= len(MsgTypeStr) {
		return nil, boo.New(boo.InvalidRead, "invalid bbs message type")
	}
	return &bm, nil
}

type BBSMessage struct {
	SendMessage
}

func (m *BBSMessage) GetMsgType() MsgType {
	return MsgType(m.SendMessage[factory.SEND_MSG_META_END])
}

func (m *BBSMessage) GetBody() []byte {
	return m.SendMessage[factory.SEND_MSG_META_END+tLen:]
}

func (m *BBSMessage) ExtractContent() (*object.Content, error) {
	v := new(object.Content)
	if e := encoder.DeserializeRaw(m.GetBody(), v); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead,
			"failed to deserialize content data")
	}
	// TODO: Verify data? maybe.
	return v, nil
}

func (m *BBSMessage) ExtractNewContentResponse() (*NewContentResponse, error) {
	v := new(NewContentResponse)
	if e := encoder.DeserializeRaw(m.GetBody(), v); e != nil {
		return nil, e
	}
	return v, nil
}

func (m *BBSMessage) ExtractDiscovererMsg() (*DiscovererMsg, error) {
	v := new(DiscovererMsg)
	if e := encoder.DeserializeRaw(m.GetBody(), v); e != nil {
		return nil, e
	}
	return v, nil
}

type NewContentResponse struct {
	Hash   cipher.SHA256 // Hash of entire message.
	Type   MsgType       // Message type.
	Seq    uint64        // Root sequence that satisfies.
	Okay   bool          // Whether successful.
	ErrTyp int64         // Type of error.
	ErrMsg string        // Message of error.
}

func GenerateNewContentResponse(msg *BBSMessage, goal uint64, e error) *NewContentResponse {
	r := &NewContentResponse{
		Hash: cipher.SumSHA256(msg.GetBody()),
		Type: msg.GetMsgType(),
		Seq:  goal,
		Okay: true,
	}
	if e != nil {
		r.Okay = false
		r.ErrTyp = int64(boo.Type(e))
		r.ErrMsg = e.Error()
	}
	return r
}
