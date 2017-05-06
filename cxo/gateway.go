package cxo

import (
	"encoding/base64"
	"fmt"
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
	bc, e := g.c.SubscribeToBoard(pk)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Display result.
	bv, e := typ.NewBoardView(bc, g.c.Client, false)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	return NewResultReply(bv)
}

// Unsubscribe unsubscribes from a board.
func (g *Gateway) Unsubscribe(pkStr string) *Reply {
	// Check public key.
	pk, e := cipher.PubKeyFromHex(pkStr)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Unsubscribe from board.
	if g.c.UnSubscribeFromBoard(pk) == false {
		return NewErrorReply("unable to unsubscribe")
	}
	return NewResultReply("successfully unsubscribed")
}

// ViewBoard views the specified board of public key.
func (g *Gateway) ViewBoard(pkStr string) *Reply {
	// Check public key.
	pk, e := cipher.PubKeyFromHex(pkStr)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Get BoardConfig.
	bc, e := g.c.BoardManager.GetConfig(pk)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Display result.
	bv, e := typ.NewBoardView(bc, g.c.Client, true)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	return NewResultReply(bv)
}

// ViewBoards lists all the boards we are subscribed to.
func (g *Gateway) ViewBoards() *Reply {
	// Get list of BoardConfigs.
	bcList := g.c.BoardManager.GetList()
	// Get list of BoardViews.
	bvList := make([]*typ.BoardView, len(bcList))
	for i := 0; i < len(bcList); i++ {
		var e error
		bvList[i], e = typ.NewBoardView(bcList[i], g.c.Client, false)
		if e != nil {
			return NewErrorReply(fmt.Sprintf("board %s: %v", bcList[i].PublicKey, e))
		}
	}
	return NewResultReply(bvList)
}

// ViewThread views the specified thread of specified board and thread id.
// TODO: Implement.
func (g *Gateway) ViewThread(bpkStr, tidStr string) *Reply {
	// Check public key.
	bpk, e := cipher.PubKeyFromHex(bpkStr)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Check thread id.
	tid, e := base64.StdEncoding.DecodeString(tidStr)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Get BoardConfig.
	bc, e := g.c.BoardManager.GetConfig(bpk)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	return NewResultReply(string(tid) + bc.Name)
}

// NewBoard creates a new master board with a seed and name.
func (g *Gateway) NewBoard(name, seed string) *Reply {
	// Create master BoardConfig.
	bc := typ.NewMasterBoardConfig(name, seed, "")
	// Inject board.
	if e := g.c.InjectBoard(bc); e != nil {
		return NewErrorReply(e.Error())
	}
	// Display result.
	bv, e := typ.NewBoardView(bc, g.c.Client, false)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	return NewResultReply(bv)
}

// NewThread adds a new thread to specified board.
func (g *Gateway) NewThread(bpkStr, name, desc string) *Reply {
	// Check public key.
	bpk, e := cipher.PubKeyFromHex(bpkStr)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Create new Thread.
	thread := typ.NewThread(name, desc)
	// Inject thread to specified board in cxo.
	if e := g.c.InjectThread(bpk, thread); e != nil {
		return NewErrorReply(e.Error())
	}
	// Display result.
	tv, e := typ.NewThreadView(bpk, thread, g.c.Client, false)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	return NewResultReply(tv)
}

// NewPost adds a new post to specified board and thread.
func (g *Gateway) NewPost(bpkStr, tidStr, name, body string) *Reply {
	// Check public key.
	bpk, e := cipher.PubKeyFromHex(bpkStr)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Check thread id.
	tid, e := cipher.PubKeyFromHex(tidStr)
	if e != nil {
		return NewErrorReply(e.Error())
	}
	// Make new post.
	post := typ.NewPost(name, body, cipher.PubKey{})
	// Inject post.
	if e := g.c.InjectPost(bpk, tid, post); e != nil {
		return NewErrorReply(e.Error())
	}
	return NewResultReply(post)
}
