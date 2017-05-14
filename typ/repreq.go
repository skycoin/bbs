package typ

// ReqRep represents a json reply object.
type ReqRep struct {
	// Subscriptions.
	Subscriptions []string `json:"subscriptions,omitempty"`

	// Boards, Threads and Posts.
	Board   *Board     `json:"board,omitempty"`
	Boards  []*Board   `json:"boards,omitempty"`
	Thread  *Thread    `json:"thread,omitempty"`
	Threads []*Thread  `json:"threads,omitempty"`
	Post    *Post      `json:"post,omitempty"`
	Posts   []*Post    `json:"posts,omitempty"`
	Req     *SubReqRep `json:"request,omitempty"`

	// Additional request stuff.
	Seed string `json:"seed,omitempty"`
	CXO  *bool  `json:"cxo,omitempty"`
}

func NewRepReq() *ReqRep {
	return &ReqRep{}
}

func (ro *ReqRep) Prepare(e error, s interface{}) *ReqRep {
	if e == nil {
		ro.Req = &SubReqRep{true, nil, s}
	} else {
		ro.Req = &SubReqRep{false, e.Error(), nil}
	}
	return ro
}

// SubReqRep represents a sub-branch of ReqRep.
type SubReqRep struct {
	Okay    bool        `json:"okay"`
	Error   interface{} `json:"error,omitempty"`
	Message interface{} `json:"message,omitempty"`
}
