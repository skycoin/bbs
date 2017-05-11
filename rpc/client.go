package rpc

import (
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

// NewPost injects a new post to specified board and thread.
func (c *Client) NewPost(bpk cipher.PubKey, tHash skyobject.Reference, post *typ.Post) (
	rep *typ.RepReq, e error,
) {
	rep = typ.NewRepReq()
	req := &NewPostReq{bpk, tHash, post}
	e = c.rpc.Call("rpc.NewPost", req, &rep)
	return
}

// NewThread injects a new thread to specified board.
func (c *Client) NewThread(bpk cipher.PubKey, thread *typ.Thread) (
	rep *typ.RepReq, e error,
) {
	rep = typ.NewRepReq()
	req := NewThreadReq{bpk, thread}
	e = c.rpc.Call("rpc.NewThread", req, &rep)
	return
}

// Close closes the rpc client.
func (c *Client) Close() error {
	return c.rpc.Close()
}
