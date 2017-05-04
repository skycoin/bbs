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
	cxo *CXO

	// For closing.
	waiter sync.WaitGroup
}

// NewServer creates a new RPC server.
func NewServer(c *cxo.Client) *Server {
	return &Server{
		rpc: rpc.NewServer(),
		cxo: NewCXO(c),
	}
}

// Launch launches the rpc server.
func (s *Server) Launch(address string) error {
	var e error
	if e = s.rpc.RegisterName("bbs", s.cxo); e != nil {
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
	var e error
	if s.l != nil {
		e = s.l.Close()
	}
	s.waiter.Wait()
	return e
}

// Address prints the rpc server's address.
func (s *Server) Address() string {
	if s.l != nil {
		return s.l.Addr().String()
	}
	return ""
}