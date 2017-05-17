package rpc

import (
	"github.com/evanlinjin/bbs/store"
	"github.com/evanlinjin/bbs/typ"
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

func (c *Client) NewPost(bpk cipher.PubKey, tRef skyobject.Reference, post *typ.Post) (
	ok *bool, e error,
) {
	ok = new(bool)
	req := &store.ReqNewPost{bpk, tRef, post}
	e = c.rpc.Call("bbs.NewPost", req, ok)
	return
}

func (c *Client) NewThread(bpk, upk cipher.PubKey, usk cipher.SecKey, thread *typ.Thread) (
	ok *bool, e error,
) {
	ok = new(bool)
	req := &store.ReqNewThread{
		bpk, upk,
		thread.Sign(usk), thread,
	}
	e = c.rpc.Call("bbs.NewThread", req, ok)
	return
}
