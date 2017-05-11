package api

import (
	"errors"
	"github.com/evanlinjin/bbs/cxo"
	"github.com/evanlinjin/bbs/rpc"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
)

// Gateway is what's exposed to the GUI.
type Gateway struct {
	c *cxo.Client
}

// NewGateWay creates a new Gateway with specified Client.
func NewGateWay(c *cxo.Client) *Gateway {
	return &Gateway{c}
}

// Subscribe subscribes to a board.
func (g *Gateway) Subscribe(pkStr string) *typ.RepReq {
	var (
		reply = typ.NewRepReq()
		e     error
		pk    cipher.PubKey
	)
	// Check public key.
	pk, e = typ.GetPubKey(pkStr)
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
func (g *Gateway) Unsubscribe(pkStr string) *typ.RepReq {
	var (
		reply = typ.NewRepReq()
		e     error
	)
	// Check public key.
	pk, e := typ.GetPubKey(pkStr)
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

// ListBoards lists all the boards we are subscribed to.
func (g *Gateway) ListBoards() *typ.RepReq {
	var (
		reply = typ.NewRepReq()
		e     error
	)
	// Get list of boards.
	reply.Boards, _, e = g.c.ObtainAllBoards()
	return reply.Prepare(e, nil)
}

// ViewBoard views the specified board of public key.
func (g *Gateway) ViewBoard(pkStr string) *typ.RepReq {
	var (
		reply = typ.NewRepReq()
		e     error
	)
	// Check public key.
	pk, e := typ.GetPubKey(pkStr)
	if e != nil {
		goto ViewBoardResults
	}
	// Obtain board information.
	reply.Board, _, reply.Threads, e = g.c.ObtainThreads(pk)
ViewBoardResults:
	return reply.Prepare(e, nil)
}

// ViewThread views the specified thread of specified board and thread id.
func (g *Gateway) ViewThread(bpkStr, tHashStr string) *typ.RepReq {
	var (
		reply = typ.NewRepReq()
		e     error
		pk    cipher.PubKey
		hash  cipher.SHA256
	)
	// Check public key.
	pk, e = typ.GetPubKey(bpkStr)
	if e != nil {
		goto ViewThreadResult
	}
	// Check hash.
	hash, e = cipher.SHA256FromHex(tHashStr)
	if e != nil {
		goto ViewThreadResult
	}
	// Prepare results.
	reply.Board, _, reply.Thread, _, reply.Posts, _, e =
		g.c.ObtainPosts(pk, skyobject.Reference(hash))
ViewThreadResult:
	return reply.Prepare(e, nil)
}

// NewBoard creates a new master board with a name, description and seed.
func (g *Gateway) NewBoard(board *typ.Board, seed string) *typ.RepReq {
	var (
		reply = typ.NewRepReq()
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
	_, _, e = g.c.InjectBoard(typ.NewMasterBoardConfig(board, g.c.B.RPCAddr, seed))
NewBoardResults:
	// Get list of boards.
	reply.Boards, _, _ = g.c.ObtainAllBoards()
	return reply.Prepare(e, "new board successfully created")
}

// NewThread adds a new thread to specified board.
func (g *Gateway) NewThread(bpkStr string, thread *typ.Thread) *typ.RepReq {
	var (
		reply = typ.NewRepReq()
		e, e2 error
		pk    cipher.PubKey
		bc    *typ.BoardConfig
		rpcc  *rpc.Client
	)
	// Check public key.
	pk, e = typ.GetPubKey(bpkStr)
	if e != nil {
		goto NewThreadResult
	}
	// See if we own the board.
	bc, e = g.c.B.GetConfig(pk)
	if e != nil {
		e = errors.New("not subscribed to this board")
		goto NewThreadResult
	}
	if bc.Master == false {
		// We don't own the board, hence inject board and obtain info via rpc.
		reply.Board, _, _ = g.c.ObtainBoard(pk)
		rpcc, e2 = rpc.NewClient(reply.Board.URL)
		if e2 != nil {
			goto NewThreadResult
		}
		reply, e = rpcc.NewThread(pk, thread)
		goto NewThreadResult
	}
	// We own the board, inject to board directly and obtain info.
	reply.Board, _, reply.Threads, e = g.c.InjectThread(pk, thread)
NewThreadResult:
	if e2 != nil {
		reply.Board, _, reply.Threads, e = g.c.ObtainThreads(pk)
		e = e2
	}
	return reply.Prepare(e, "thread successfully created")
}

// NewPost adds a new post to specified board and thread.
func (g *Gateway) NewPost(bpkStr, tHashStr string, post *typ.Post) *typ.RepReq {
	var (
		reply = typ.NewRepReq()
		e, e2 error
		pk    cipher.PubKey
		hash  cipher.SHA256
		bc    *typ.BoardConfig
		u     *typ.User
		rpcc  *rpc.Client
	)
	// Check public key.
	pk, e = typ.GetPubKey(bpkStr)
	if e != nil {
		goto NewPostResult
	}
	// Check hash.
	hash, e = cipher.SHA256FromHex(tHashStr)
	if e != nil {
		goto NewPostResult
	}
	// See if we own the board.
	bc, e = g.c.B.GetConfig(pk)
	if e != nil {
		e = errors.New("not subscribed to this board")
		goto NewPostResult
	}
	if bc.Master == false {
		// We don't own the board, hence inject post to thread via rpc.
		// TODO: RPC part.
		reply.Board, _, reply.Thread, e = g.c.ObtainThread(pk, skyobject.Reference(hash))
		rpcc, e = rpc.NewClient(reply.Board.URL)
		if e != nil {
			goto NewPostResultCanDisplay
		}
		reply, e = rpcc.NewPost(pk, skyobject.Reference(hash), post)
		goto NewPostResult
	}
	// We own the board, inject to board directly and obtain info.
	// Authorise post.
	u, e2 = g.c.U.GetMaster()
	if e2 != nil {
		goto NewPostResultCanDisplay
	}
	post.Creator = u.PublicKey
	post.CreatorStr = u.PublicKeyStr
	post.Touch()
	post.Sign(u.SecretKey)
NewPostResultCanDisplay:
	if e2 != nil {
		// Because of authorisation error, unable to inject post.
		reply.Board, _, reply.Thread, _, reply.Posts, _, e =
			g.c.ObtainPosts(pk, skyobject.Reference(hash))
		if e != nil {
			return reply.Prepare(e, nil)
		}
		return reply.Prepare(e2, nil)
	}
	// Inject post.
	reply.Board, _, reply.Thread, _, reply.Posts, e =
		g.c.InjectPost(pk, skyobject.Reference(hash), post)
NewPostResult:
	return reply.Prepare(e, "post successfully created")
}
