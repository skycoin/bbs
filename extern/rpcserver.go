package extern

import (
	"net"
	"net/rpc"
	"strconv"
	"sync"
)

type RPCServer struct {
	l      net.Listener
	rpc    *rpc.Server
	g      *RPCGateway
	waiter sync.WaitGroup
}

func NewRPCServer(g *RPCGateway, port int) (*RPCServer, error) {
	s := &RPCServer{
		rpc: rpc.NewServer(),
		g:   g,
	}
	if e := s.open("[::]:" + strconv.Itoa(port)); e != nil {
		return nil, e
	}
	return s, nil
}

func (s *RPCServer) open(address string) error {
	var e error
	if e = s.rpc.RegisterName("bbs", s.g); e != nil {
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

// Close closes the rpc server.
func (s *RPCServer) Close() error {
	if s == nil {
		return nil
	}
	var e error
	if s.l != nil {
		e = s.l.Close()
	}
	s.waiter.Wait()
	return e
}

// Address prints the rpc server's address.
func (s *RPCServer) Address() string {
	if s.l != nil {
		return s.l.Addr().String()
	}
	return ""
}
