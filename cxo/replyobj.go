package cxo

import (
	"encoding/hex"
	"errors"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/skycoin/src/cipher"
)

// JsonApiObj represents a json reply object.
type JsonApiObj struct {
	Board   *typ.Board    `json:"board,omitempty"`
	Boards  []*typ.Board  `json:"boards,omitempty"`
	Thread  *typ.Thread   `json:"thread,omitempty"`
	Threads []*typ.Thread `json:"threads,omitempty"`
	Posts   []*typ.Post   `json:"posts,omitempty"`
	Req     *ReqObj       `json:"request,omitempty"`

	// Request stuff
	Seed string `json:"seed,omitempty"`
}

func NewJsonApiObj() *JsonApiObj {
	return &JsonApiObj{}
}

func (ro *JsonApiObj) Prepare(e error, s interface{}) *JsonApiObj {
	if e == nil {
		ro.Req = &ReqObj{true, nil, s}
	} else {
		ro.Req = &ReqObj{false, e.Error(), nil}
	}
	return ro
}

// PutRequestObj represents a sub-branch of JsonApiObj.
type ReqObj struct {
	Okay    bool        `json:"okay"`
	Error   interface{} `json:"error,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

// GetPubKey obtains a public key from string, avoiding panics.
func GetPubKey(s string) (cipher.PubKey, error) {
	b, e := hex.DecodeString(s)
	if e != nil || len(b) != len(cipher.PubKey{}) {
		return cipher.PubKey{}, errors.New("invalid public key")
	}
	return cipher.NewPubKey(b), nil
}
