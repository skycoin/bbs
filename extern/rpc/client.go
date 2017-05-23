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

//TODO: needs refactoring

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

func (c *Client) RemoveBoard(req *ReqRemoveBoard) (ok *bool, e error) {
	ok = new(bool)
	e = c.rpc.Call("Gateway.RemoveBoard", req, ok)
	return
}

func (c *Client) RemoveThread(req *ReqRemoveThread) (ok *bool, e error) {
	ok = new(bool)
	e = c.rpc.Call("Gateway.RemoveThread", req, ok)
	return
}

func (c *Client) RemovePost(req *ReqRemovePost) (ok *bool, e error) {
	ok = new(bool)
	e = c.rpc.Call("Gateway.RemovePost", req, ok)
	return
}
