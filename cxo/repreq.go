package cxo

import (
	"github.com/evanlinjin/bbs/typ"
)

// RepReq represents a json reply object.
type RepReq struct {
	Board   *typ.Board    `json:"board,omitempty"`
	Boards  []*typ.Board  `json:"boards,omitempty"`
	Thread  *typ.Thread   `json:"thread,omitempty"`
	Threads []*typ.Thread `json:"threads,omitempty"`
	Post    *typ.Post     `json:"post,omitempty"`
	Posts   []*typ.Post   `json:"posts,omitempty"`
	Req     *ReqObj       `json:"request,omitempty"`

	// Request stuff
	Seed string `json:"seed,omitempty"`
}

func NewRepReq() *RepReq {
	return &RepReq{}
}

func (ro *RepReq) Prepare(e error, s interface{}) *RepReq {
	if e == nil {
		ro.Req = &ReqObj{true, nil, s}
	} else {
		ro.Req = &ReqObj{false, e.Error(), nil}
	}
	return ro
}

// PutRequestObj represents a sub-branch of RepReq.
type ReqObj struct {
	Okay    bool        `json:"okay"`
	Error   interface{} `json:"error,omitempty"`
	Message interface{} `json:"message,omitempty"`
}
