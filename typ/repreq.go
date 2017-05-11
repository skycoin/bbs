package typ

// RepReq represents a json reply object.
type RepReq struct {
	Board   *Board    `json:"board,omitempty"`
	Boards  []*Board  `json:"boards,omitempty"`
	Thread  *Thread   `json:"thread,omitempty"`
	Threads []*Thread `json:"threads,omitempty"`
	Post    *Post     `json:"post,omitempty"`
	Posts   []*Post   `json:"posts,omitempty"`
	Req     *ReqObj   `json:"request,omitempty"`

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
