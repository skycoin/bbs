package rpc

import (
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
	e = c.rpc.Call("Gateway.NewPost", req, ok)
	return
}

func (c *Client) NewThread(req *ReqNewThread) (ok *bool, e error) {
	ok = new(bool)
	e = c.rpc.Call("Gateway.NewThread", req, ok)
	return
}
