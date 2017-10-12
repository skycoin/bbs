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
	c      *ServerConfig
	l      *log.Logger
	listen net.Listener
	rpc    *rpc.Server
	api    *Gateway
	wg     sync.WaitGroup
}

func NewServer(c *ServerConfig, g *Gateway) (*Server, error) {
	if !*c.Enable {
		return nil, nil
	}
	s := &Server{
		c:   c,
		l:   inform.NewLogger(true, os.Stdout, logPrefix),
		rpc: rpc.NewServer(),
		api: g,
	}
	if e := s.rpc.Register(g); e != nil {
		return nil, e
	}
	address := "[::]:" + strconv.Itoa(*c.Port)
	var e error
	s.listen, e = net.Listen("tcp", address)
	if e != nil {
		return nil, e
	}
	s.l.Printf("listening on: '%s'", address)
	s.wg.Add(1)
	go func(l net.Listener) {
		defer s.wg.Done()
		s.rpc.Accept(l)
		s.l.Println("Closed.")
	}(s.listen)
	return s, nil
}

func (s *Server) Close() {
	if s != nil {
		if s.l != nil {
			if e := s.listen.Close(); e != nil {
				s.l.Println("Error on close:", e)
			}
		}
		s.wg.Wait()
	}
}
