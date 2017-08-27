package rpc

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
)

// Gateway is the access point for remote clients to interact with.
type Gateway struct {
	CXO *cxo.Manager
}

// NewThread accepts a request to create a new thread, returning a sequence goal or error.
func (g *Gateway) NewThread(thread *object.Thread, goal *uint64) error {
	if thread == nil || goal == nil {
		return boo.New(boo.InvalidInput, "nil error")
	}
	if e := thread.Verify(); e != nil {
		return e
	}
	bi, e := g.CXO.GetBoardInstance(thread.OfBoard)
	if e != nil {
		return e
	}
	if !bi.IsMaster() {
		return notMasterErr(thread.OfBoard)
	}
	return send(goal)(bi.NewThread(thread))
}

// NewPost accepts a request to create a new post, returning a sequence goal or error.
func (g *Gateway) NewPost(post *object.Post, goal *uint64) error {
	if post == nil || goal == nil {
		return boo.New(boo.InvalidInput, "nil error")
	}
	if e := post.Verify(); e != nil {
		return e
	}
	bi, e := g.CXO.GetBoardInstance(post.OfBoard)
	if e != nil {
		return e
	}
	if !bi.IsMaster() {
		return notMasterErr(post.OfBoard)
	}
	return send(goal)(bi.NewPost(post))
}

// NewVote accepts a request to create a new vote, returning a sequence goal or error.
func (g *Gateway) NewVote(vote *object.Vote, goal *uint64) error {
	if vote == nil || goal == nil {
		return boo.New(boo.InvalidInput, "nil error")
	}
	if e := vote.Verify(); e != nil {
		return e
	}
	bi, e := g.CXO.GetBoardInstance(vote.OfBoard)
	if e != nil {
		return e
	}
	if !bi.IsMaster() {
		return notMasterErr(vote.OfBoard)
	}
	return send(goal)(bi.NewVote(vote))
}

func notMasterErr(bpk cipher.PubKey) error {
	return boo.Newf(boo.NotAllowed,
		"node is not owner of board of public key '%s'", bpk.Hex())
}

func send(goal *uint64) func(g uint64, e error) error {
	return func(g uint64, e error) error {
		*goal = g
		return e
	}
}
