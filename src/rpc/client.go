package rpc

import (
	"github.com/pkg/errors"
	"net/rpc"
	"time"
)

// Client represents a RPC Client.
type Client struct {
	rpc  *rpc.Client
	wait time.Duration
}

// NewClient creates a new rpc client with given address.
func NewClient(address string) (*Client, error) {
	c, e := rpc.Dial("tcp", address)
	rpcc := Client{
		rpc:  c,
		wait: 10 * time.Second,
	}
	return &rpcc, e
}

func (c *Client) PingPong() (ok bool, e error) {
	e = c.result(c.rpc.Go("Gateway.PingPong", nil, &ok, nil))
	return
}

func (c *Client) NewPost(req *ReqNewPost) (ref string, e error) {
	e = c.result(c.rpc.Go("Gateway.NewPost", req, &ref, nil))
	return
}

func (c *Client) NewThread(req *ReqNewThread) (ref string, e error) {
	e = c.result(c.rpc.Go("Gateway.NewThread", req, &ref, nil))
	return
}

func (c *Client) VotePost(req *ReqVotePost) (ok bool, e error) {
	e = c.result(c.rpc.Go("Gateway.VotePost", req, &ok, nil))
	return
}

func (c *Client) VoteThread(req *ReqVoteThread) (ok bool, e error) {
	e = c.result(c.rpc.Go("Gateway.VoteThread", req, &ok, nil))
	return
}

func (c *Client) result(call *rpc.Call) error {
	timer := time.NewTimer(c.wait)
	select {
	case <-call.Done:
		return call.Error
	case <-timer.C:
		return errors.New("timeout")
	}
}
