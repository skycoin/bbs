package rpc

import (
	"github.com/evanlinjin/bbs/cxo"
	"net"
	"net/rpc"
	"sync"
)

// Server represents a RPC server.
type Server struct {
	l   net.Listener
	rpc *rpc.Server
	cxo *cxo.Client

	// For closing.
	waiter sync.WaitGroup
}

// NewServer creates a new RPC server.
func NewServer(c *cxo.Client) *Server {
	return &Server{
		rpc: rpc.NewServer(),
		cxo: c,
	}
}

// Launch launches the rpc server.
func (s *Server) Launch(address string) error {
	var e error
	if e = s.rpc.RegisterName("bbs", s); e != nil {
		return e
	}
	if s.l, e = net.Listen("tcp", address); e != nil {
		return e
	}
	s.waiter.Add(1)
	go func(l net.Listener) {
		defer s.waiter.Done()
		s.rpc.Accept(l)
	}(s.l)
	return nil
}

// Shutdown shutdowns the rpc server.
func (s *Server) Shutdown() error {
	if s.l != nil {
		return s.l.Close()
	}
	return nil
}

/******************************************************
 * Exposed RPC Methods                                *
 *****************************************************/

// NewPost injects a new post to specified board and thread.
// TODO: Implement.
func (s *Server) NewPost(req *NewPostReq, ok *bool) error {
	*ok = true
	return nil
}

// NewThread injects a new thread to specified board.
// TODO: Implement.
func (s *Server) NewThread(req *NewThreadReq, ok *bool) error {
	*ok = true
	return nil
}
