package extern

import "github.com/evanlinjin/bbs/typ"

// ReqRes represents a json reply object.
type ReqRes struct {
	// Subscriptions.
	Subscriptions []string `json:"subscriptions,omitempty"`

	// Boards, Threads and Posts.
	Board   *typ.Board    `json:"board,omitempty"`
	Boards  []*typ.Board  `json:"boards,omitempty"`
	Thread  *typ.Thread   `json:"thread,omitempty"`
	Threads []*typ.Thread `json:"threads,omitempty"`
	Post    *typ.Post     `json:"post,omitempty"`
	Posts   []*typ.Post   `json:"posts,omitempty"`
	Req     *SubReqRep    `json:"request,omitempty"`

	// Additional request stuff.
	Seed string `json:"seed,omitempty"`
	CXO  *bool  `json:"cxo,omitempty"`
}

func NewRepRes() *ReqRes {
	return &ReqRes{}
}

func (rr *ReqRes) Prepare(e error, s interface{}) *ReqRes {
	if e == nil {
		rr.Req = &SubReqRep{true, nil, s}
	} else {
		rr.Req = &SubReqRep{false, e.Error(), nil}
	}
	return rr
}

// SubReqRep represents a sub-branch of ReqRep.
type SubReqRep struct {
	Okay    bool        `json:"okay"`
	Error   interface{} `json:"error,omitempty"`
	Message interface{} `json:"message,omitempty"`
}
