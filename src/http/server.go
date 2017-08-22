package http

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc/inform"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	"github.com/skycoin/bbs/src/store/cxo"
)

const (
	httpLogPrefix     = "HTTPSERVER"
	httpLocalhost     = "127.0.0.1"
	httpIndexFileName = "index.html"
)

// ServerConfig represents a HTTP server configuration file.
type ServerConfig struct {
	Port      *int
	StaticDir *string
	EnableGUI *bool
}

// Server represents an HTTP Server that serves static files and JSON api.
type Server struct {
	c    *ServerConfig
	l    *log.Logger
	net  net.Listener
	mux  *http.ServeMux
	api  *Gateway
	quit chan struct{}
}

// NewServer creates a new server.
func NewServer(config *ServerConfig, api *Gateway) (*Server, error) {
	server := &Server{
		c:    config,
		l:    inform.NewLogger(true, os.Stdout, httpLogPrefix),
		mux:  http.NewServeMux(),
		api:  api,
		quit: make(chan struct{}),
	}
	var e error
	if *config.StaticDir, e = filepath.Abs(*config.StaticDir); e != nil {
		return nil, e
	}
	host := fmt.Sprintf("%s:%d", httpLocalhost, *config.Port)
	if server.net, e = net.Listen("tcp", host); e != nil {
		return nil, e
	}
	if e := server.prepareMux(); e != nil {
		return nil, e
	}
	go server.serve()
	return server, nil
}

func (s *Server) serve() {
	for {
		if e := http.Serve(s.net, s.mux); e != nil {
			select {
			case <-s.quit:
				return
			default:
				time.Sleep(100 * time.Millisecond)
				continue
			}
		}
	}
}

func (s *Server) prepareMux() error {
	if *s.c.EnableGUI {
		if e := s.prepareStatic(); e != nil {
			return e
		}
	}
	return s.api.prepare(s.mux)
}

func (s *Server) prepareStatic() error {
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, e := ioutil.ReadFile(path.Join(*s.c.StaticDir, httpIndexFileName))
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(e.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	return filepath.Walk(*s.c.StaticDir, func(path string, info os.FileInfo, e error) error {
		if info == nil || info.IsDir() {
			return nil
		}
		httpPath := strings.TrimPrefix(path, *s.c.StaticDir)
		s.mux.HandleFunc(httpPath, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path)
		})
		return nil
	})
}

func (s *Server) CXO() *cxo.Manager {
	return s.api.Access.CXO
}

func (s *Server) Close() {
	if s.quit != nil {
		close(s.quit)
		s.net.Close()
		s.net = nil
	}
}
