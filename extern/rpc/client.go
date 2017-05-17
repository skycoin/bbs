package rpc

import (
	"github.com/evanlinjin/bbs/store/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"net/rpc"
)

// Client represents a RPC Client.
type Client struct {
	rpc *rpc.Client
}

// NewClient creates a new rpc client with given address.
func NewClient(address string) (*Client, error) {
	c, e := rpc.Dial("tcp", address)
	rpcc := Client{
		rpc: c,
	}
	return &rpcc, e
}

func (c *Client) NewPost(req *ReqNewPost) (ok *bool, e error) {
	ok = new(bool)
	e = c.rpc.Call("bbs.NewPost", req, ok)
	return
}

func (c *Client) NewPostOld(bpk cipher.PubKey, tRef skyobject.Reference, post *typ.Post) (
	ok *bool, e error,
) {
	ok = new(bool)
	req := &ReqNewPost{bpk, tRef, post}
	e = c.rpc.Call("bbs.NewPost", req, ok)
	return
}

func (c *Client) NewThread(req *ReqNewThread) (ok *bool, e error) {
	ok = new(bool)
	e = c.rpc.Call("bbs.NewThread", req, ok)
	return
}

func (c *Client) NewThreadOld(bpk, upk cipher.PubKey, usk cipher.SecKey, thread *typ.Thread) (
	ok *bool, e error,
) {
	ok = new(bool)
	req := &ReqNewThread{
		bpk, upk,
		thread.Sign(usk), thread,
	}
	e = c.rpc.Call("bbs.NewThread", req, ok)
	return
}
