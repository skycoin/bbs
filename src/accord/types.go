package accord

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

// Type determines the message type of transported data.
type Type byte

// String obtains the string representation of the message type.
func (t Type) String() string {
	if int(t) >= len(MsgTypeStr) || int(t) < 0 {
		return "Unknown"
	} else {
		return MsgTypeStr[int(t)]
	}
}

// IsValid checks whether the message type is valid.
func (t Type) IsValid() bool {
	return int(t) < len(MsgTypeStr)
}

const (
	// Length of type byte in discovery message.
	tLen = 1
	// SubmissionType symbolises a message that is for submitting content to a board.
	SubmissionType Type = iota << 0
	// SubmissionResponseType symbolises a message that is a response for a "SubmissionType" message.
	SubmissionResponseType
)

// MsgTypeStr provides the string representations of the message types.
var MsgTypeStr = [...]string{
	SubmissionType:         "Submission",
	SubmissionResponseType: "Submission Response",
}

// Wrapper wraps data send via messenger.
type Wrapper struct {
	Raw []byte
}

func NewWrapper(raw []byte) (*Wrapper, error) {
	out := &Wrapper{Raw: raw}

	// Check if message is large enough to contain discovery header.
	if len(raw) < factory.SEND_MSG_META_END+tLen {
		return nil, boo.New(boo.InvalidRead, "invalid length")
	}
	// Check if OP is "OP_SEND".
	if raw[0] != factory.OP_SEND {
		return nil, boo.Newf(boo.InvalidRead, "op type not send")
	}
	// Check "from" public key.
	if e := out.GetFromPK().Verify(); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead, "invalid from public key")
	}
	// Check "to" public key.
	if e := out.GetToPK().Verify(); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead, "invalid to public key")
	}
	// Check type.
	if out.GetType().IsValid() == false {
		return nil, boo.New(boo.InvalidRead, "invalid type")
	}
	return out, nil
}

func (w *Wrapper) GetFromPK() cipher.PubKey {
	return cipher.NewPubKey(
		w.Raw[factory.SEND_MSG_PUBLIC_KEY_BEGIN:factory.SEND_MSG_PUBLIC_KEY_END])
}

func (w *Wrapper) GetToPK() cipher.PubKey {
	return cipher.NewPubKey(
		w.Raw[factory.SEND_MSG_TO_PUBLIC_KEY_BEGIN:factory.SEND_MSG_TO_PUBLIC_KEY_END])
}

func (w *Wrapper) GetType() Type {
	return Type(w.Raw[factory.SEND_MSG_META_END])
}

func (w *Wrapper) GetBody() []byte {
	return w.Raw[factory.SEND_MSG_META_END+tLen:]
}

func (w *Wrapper) ToSubmission() (*Submission, error) {
	out := new(Submission)
	if e := encoder.DeserializeRaw(w.GetBody(), out); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead, "submission corrupt")
	}
	return out, nil
}

func (w *Wrapper) ToSubmissionResponse() (*SubmissionResponse, error) {
	out := new(SubmissionResponse)
	if e := encoder.DeserializeRaw(w.GetBody(), out); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead, "submission response corrupt")
	}
	return out, nil
}

// Submission is the content of the message that is for submitting content to a board.
type Submission struct {
	Raw []byte
	Sig cipher.Sig
}

func (s *Submission) GetHash() cipher.SHA256 {
	return cipher.SumSHA256(s.Raw)
}

func (s *Submission) ToTransport() (*object.Transport, error) {
	return object.NewTransport(s.Raw, s.Sig)
}

// SubmissionResponse is the content of a message that acts as a response to a "Submission".
type SubmissionResponse struct {
	Hash   cipher.SHA256 // Hash of the submission body.
	Okay   bool          // Whether submission was successful.
	Seq    uint64        // (only on success) root sequence in which content is successfully submitted.
	ErrTyp int64         // (only on failure) error type.
	ErrMsg string        // (only on failure) error message.
}

func NewSubmissionResponse(hash cipher.SHA256, seq uint64, e error) *SubmissionResponse {
	if e == nil {
		return &SubmissionResponse{
			Hash: hash,
			Okay: true,
			Seq:  seq,
		}
	} else {
		return &SubmissionResponse{
			Hash:   hash,
			Okay:   false,
			ErrTyp: int64(boo.Type(e)),
			ErrMsg: e.Error(),
		}
	}
}

func (sr *SubmissionResponse) Serialize() []byte {
	return encoder.Serialize(sr)
}

func (sr *SubmissionResponse) Error() error {
	if sr.Okay {
		return nil
	} else {
		return boo.New(int(sr.ErrTyp), sr.ErrMsg)
	}
}
