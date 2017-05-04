package rpc

import "net/rpc"

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
func (c *Client) NewPost(req *NewPostReq) (ok bool, e error) {
	e = c.rpc.Call("rpc.NewPost", req, &ok)
	return
}

// NewThread injects a new thread to specified board.
func (c *Client) NewThread(req *NewThreadReq) (ok bool, e error) {
	e = c.rpc.Call("rpc.NewThread", req, &ok)
	return
}