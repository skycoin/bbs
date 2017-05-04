package rpc

import "github.com/evanlinjin/bbs/cxo"

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
func (c *CXO) NewPost(req *NewPostReq, ok *bool) error {
	*ok = true
	return nil
}

// NewThread injects a new thread to specified board.
// TODO: Implement.
func (c *CXO) NewThread(req *NewThreadReq, ok *bool) error {
	*ok = true
	return nil
}
