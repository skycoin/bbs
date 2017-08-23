package rpc

import (
	"github.com/skycoin/bbs/src/misc/inform"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
)

const (
	logPrefix = "RPC"
)

type ServerConfig struct {
	Port   *int
	Enable *bool
}

type Server struct {
	c   *ServerConfig
	l   *log.Logger
	lis net.Listener
	rpc *rpc.Server
	api *Gateway
	wg  sync.WaitGroup
}

func NewServer(config *ServerConfig, api *Gateway) (*Server, error) {
	server := &Server{
		c:   config,
		l:   inform.NewLogger(true, os.Stdout, logPrefix),
		rpc: rpc.NewServer(),
		api: api,
	}
	if e := server.open(":" + strconv.Itoa(*config.Port)); e != nil {
		return nil, e
	}
	return server, nil
}

func (s *Server) open(address string) error {
	var e error
	if e = s.rpc.Register(s.api); e != nil {
		return e
	}
	if s.lis, e = net.Listen("tcp", address); e != nil {
		return e
	}
	s.wg.Add(1)
	go func(l net.Listener) {
		defer s.wg.Done()
		s.rpc.Accept(l)
		s.l.Println("Closed.")
	}(s.lis)
	return nil
}

func (s *Server) Close() {
	if s != nil {
		if s.l != nil {
			if e := s.lis.Close(); e != nil {
				s.l.Println("Error on close:", e)
			}
		}
		s.wg.Wait()
	}
}

// Address prints the rpc server's address.
func (s *Server) Address() string {
	if s != nil && s.lis != nil {
		return s.lis.Addr().String()
	}
	return ""
}
