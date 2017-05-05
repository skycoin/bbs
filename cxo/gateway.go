package cxo

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/evanlinjin/bbs/types"
	"fmt"
)

// Reply represents a json reply.
type Reply struct {
	Okay   bool        `json:"okay"`
	Result interface{} `json:"result,omitempty"`
	Error  interface{} `json:"error,omitempty"`
}

// NewResultReply creates a new result Reply.
func NewResultReply(v interface{}) *Reply {
	return &Reply{Okay: true, Result: v}
}

// NewErrorReply creates a new error Reply.
func NewErrorReply(v interface{}) *Reply {
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
		return NewErrorReply(e)
	}
	// Subscribe to board.
	bc, e := g.c.SubscribeToBoard(pk)
	if e != nil {
		return NewErrorReply(e)
	}
	// Display result.
	bv, e := types.NewBoardView(bc, g.c.Client)
	if e != nil {
		return NewErrorReply(e)
	}
	return NewResultReply(bv)
}

// Unsubscribe unsubscribes from a board.
func (g *Gateway) Unsubscribe(pkStr string) *Reply {
	// Check public key.
	pk, e := cipher.PubKeyFromHex(pkStr)
	if e != nil {
		return NewErrorReply(e)
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
		return NewErrorReply(e)
	}
	// Get BoardConfig.
	bc, e := g.c.BoardManager.GetConfig(pk)
	if e != nil {
		return NewErrorReply(e)
	}
	// Display result.
	bv, e := types.NewBoardView(bc, g.c.Client)
	if e != nil {
		return NewErrorReply(e)
	}
	return NewResultReply(bv)
}

// ViewBoards lists all the boards we are subscribed to.
func (g *Gateway) ViewBoards() *Reply {
	// Get list of BoardConfigs.
	bcList := g.c.BoardManager.GetList()
	// Get list of BoardViews.
	bvList := make([]*types.BoardView, len(bcList))
	for i := 0; i < len(bcList); i++ {
		var e error
		bvList[i], e = types.NewBoardView(bcList[i], g.c.Client)
		if e != nil {
			return NewErrorReply(fmt.Sprintf("board %s: %v", bcList[i].PublicKey, e))
		}
	}
	return NewResultReply(bvList)
}

// NewBoard creates a new master board with a seed and name.
func (g *Gateway) NewBoard(name, seed string) *Reply {
	// Create master BoardConfig.
	bc := types.NewMasterBoardConfig(name, seed, "")
	// Inject board.
	if e := g.c.InjectBoard(bc); e != nil {
		return NewErrorReply(e)
	}
	// Display result.
	bv, e := types.NewBoardView(bc, g.c.Client)
	if e != nil {
		return NewErrorReply(e)
	}
	return NewResultReply(bv)
}