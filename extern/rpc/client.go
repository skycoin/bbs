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

func (c *Client) NewPost(req *ReqNewPost) (ref string, e error) {
	e = c.rpc.Call("Gateway.NewPost", req, &ref)
	return
}

func (c *Client) NewThread(req *ReqNewThread) (ref string, e error) {
	e = c.rpc.Call("Gateway.NewThread", req, &ref)
	return
}

func (c *Client) VotePost(req *ReqVotePost) (ok bool, e error) {
	e = c.rpc.Call("Gateway.VotePost", req, &ok)
	return
}

func (c *Client) VoteThread(req *ReqVoteThread) (ok bool, e error) {
	e = c.rpc.Call("Gateway.VoteThread", req, &ok)
	return
}
