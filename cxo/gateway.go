package cxo

import (
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
func (g *Gateway) Subscribe(pkStr string) *Reply {
	// Check public key.
	pk, e := cipher.PubKeyFromHex(pkStr)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Subscribe to board.
	bc, e := g.c.Subscribe(pk)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Display result.
	return NewResultReply(bc.PubKey)
}

// Unsubscribe unsubscribes from a board.
func (g *Gateway) Unsubscribe(pkStr string) *Reply {
	// Check public key.
	pk, e := cipher.PubKeyFromHex(pkStr)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Unsubscribe from board.
	if g.c.Unsubscribe(pk) == false {
		return NewErrorReply("unable to unsubscribe")
	}
	return NewResultReply("successfully unsubscribed")
}

// ViewBoard views the specified board of public key.
func (g *Gateway) ViewBoard(pkStr string) *Reply {
	return nil
}

// ViewBoards lists all the boards we are subscribed to.
func (g *Gateway) ViewBoards() *Reply {

	return nil
}

// ViewThread views the specified thread of specified board and thread id.
// TODO: Implement.
func (g *Gateway) ViewThread(bpkStr, tidStr string) *Reply {
	return nil
}

// NewBoard creates a new master board with a name, description and seed.
func (g *Gateway) NewBoard(name, desc, seed string) *Reply {
	board := typ.NewBoard(name, desc, "")
	if e := g.c.InjectBoard(board, seed); e != nil {
		return NewErrorReply(e.Error())
	}
	return NewResultReply(board)
}

// NewThread adds a new thread to specified board.
func (g *Gateway) NewThread(bpkStr, name, desc string) *Reply {
	return nil
}

// NewPost adds a new post to specified board and thread.
func (g *Gateway) NewPost(bpkStr, tidStr, name, body string) *Reply {
	return nil
}
