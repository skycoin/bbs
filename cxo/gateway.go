package cxo

import (
	"errors"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/skycoin/src/cipher"
)

// Reply represents a json reply.
type Reply struct {
	Okay   bool        `json:"okay"`
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// NewResultReply creates a new result Reply.
func NewResultReply(v interface{}) *Reply {
	return &Reply{Okay: true, Result: v}
}

// NewErrorReply creates a new error Reply.
func NewErrorReply(v string) *Reply {
	return &Reply{Okay: false, Error: v}
}

// Gateway is what's exposed to the GUI.
type Gateway struct {
	c *Client
}

// NewGateWay creates a new Gateway with specified Client.
func NewGateWay(c *Client) *Gateway {
	return &Gateway{c}
}

// Subscribe subscribes to a board.
func (g *Gateway) Subscribe(pkStr string) *JsonApiObj {
	var (
		reply = NewJsonApiObj()
		e     error
		pk    cipher.PubKey
	)
	// Check public key.
	pk, e = GetPubKey(pkStr)
	if e != nil {
		goto SubscribeResult
	}
	// Subscribe to board.
	_, e = g.c.Subscribe(pk)
	if e != nil {
		goto SubscribeResult
	}
SubscribeResult:
	// Get list of boards.
	reply.Boards, _, _ = g.c.ObtainAllBoards()
	return reply.Prepare(e, "subscribed to "+pkStr)
}

// Unsubscribe unsubscribes from a board.
func (g *Gateway) Unsubscribe(pkStr string) *JsonApiObj {
	var (
		reply = NewJsonApiObj()
		e     error
	)
	// Check public key.
	pk, e := GetPubKey(pkStr)
	if e != nil {
		goto UnsubscribeResult
	}
	// Unsubscribe from board.
	if g.c.Unsubscribe(pk) == false {
		e = errors.New("unable to unsubscribe")
		goto UnsubscribeResult
	}
UnsubscribeResult:
	// Get list of boards.
	reply.Boards, _, _ = g.c.ObtainAllBoards()
	return reply.Prepare(e, "unsubscribed to "+pkStr)
}

// ViewBoard views the specified board of public key.
func (g *Gateway) ViewBoard(pkStr string) *JsonApiObj {
	var (
		reply = NewJsonApiObj()
		e     error
	)
	// Check public key.
	pk, e := GetPubKey(pkStr)
	if e != nil {
		goto ViewBoardResults
	}
ViewBoardResults:
	// Get list of boards.
	reply.Board, _, _ = g.c.ObtainBoard(pk)
	return reply.Prepare(e, "")
}

// ViewBoards lists all the boards we are subscribed to.
func (g *Gateway) ViewBoards() *JsonApiObj {
	var (
		reply = NewJsonApiObj()
		e     error
	)
	// Get list of boards.
	reply.Boards, _, e = g.c.ObtainAllBoards()
	return reply.Prepare(e, nil)
}

// ViewThread views the specified thread of specified board and thread id.
func (g *Gateway) ViewThread(bpkStr, tidStr string) *JsonApiObj {
	return NewJsonApiObj()
}

// NewBoard creates a new master board with a name, description and seed.
func (g *Gateway) NewBoard(board *typ.Board, seed string) *JsonApiObj {
	var (
		reply = NewJsonApiObj()
		e     error
	)
	// Check seed.
	if len(seed) == 0 {
		e = errors.New("invalid seed")
		goto NewBoardResults
	}
	// Check board.
	if e = board.CheckAndPrep(); e != nil {
		goto NewBoardResults
	}
	// Inject board.
	e = g.c.InjectBoard(typ.NewMasterBoardConfig(board, seed))
NewBoardResults:
	// Get list of boards.
	reply.Boards, _, _ = g.c.ObtainAllBoards()
	return reply.Prepare(e, "new board successfully created")
}

// NewThread adds a new thread to specified board.
func (g *Gateway) NewThread(bpkStr, name, desc string) *JsonApiObj {
	return nil
}

// NewPost adds a new post to specified board and thread.
func (g *Gateway) NewPost(bpkStr, tidStr, name, body string) *JsonApiObj {
	return nil
}
