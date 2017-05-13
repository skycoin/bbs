package rpc

import (
	"github.com/evanlinjin/bbs/cxo"
	"github.com/evanlinjin/bbs/typ"
)

// CXO represents what's exposed over rpc.
type CXO struct {
	c *cxo.Client
}

// NewCXO creates a new CXO from given cxo client.
func NewCXO(c *cxo.Client) *CXO {
	return &CXO{
		c: c,
	}
}

// NewPost injects a new post to specified board and thread.
// TODO: Implement.
func (c *CXO) NewPost(req *NewPostReq, rep *typ.ReqRep) (e error) {
	rep.Board, _, rep.Thread, _, rep.Posts, e =
		c.c.InjectPost(req.PK, req.Hash, req.Post)
	return
}

// NewThread injects a new thread to specified board.
// TODO: Implement.
func (c *CXO) NewThread(req *NewThreadReq, rep *typ.ReqRep) (e error) {
	rep.Board, _, rep.Threads, e =
		c.c.InjectThread(req.PK, req.Thread)
	return
}
