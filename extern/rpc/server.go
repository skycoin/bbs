package rpc

import (
	"log"
	"net"
	"net/rpc"
	"strconv"
	"sync"
)

type Server struct {
	l      net.Listener
	rpc    *rpc.Server
	g      *Gateway
	waiter sync.WaitGroup
}

func NewServer(g *Gateway, port int) (*Server, error) {
	s := &Server{
		rpc: rpc.NewServer(),
		g:   g,
	}
	if e := s.open(":" + strconv.Itoa(port)); e != nil {
		return nil, e
	}
	return s, nil
}

func (s *Server) open(address string) error {
	var e error
	if e = s.rpc.Register(s.g); e != nil {
		return e
	}
	if s.l, e = net.Listen("tcp", address); e != nil {
		return e
	}
	s.waiter.Add(1)
	go func(l net.Listener) {
		defer s.waiter.Done()
		s.rpc.Accept(l)
		log.Println("[RPCSERVER] Closed.")
	}(s.l)
	return nil
}

// Close closes the rpc server.
func (s *Server) Close() error {
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
func (s *Server) Address() string {
	if s.l != nil {
		return s.l.Addr().String()
	}
	return ""
}
