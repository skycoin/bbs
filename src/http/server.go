package http

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/cxo"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

const (
	logPrefix     = "HTTP_SERVER"
	indexFileName = "index.html"
)

// ServerConfig represents a HTTP server configuration file.
type ServerConfig struct {
	Port      *int
	StaticDir *string
	EnableGUI *bool
	EnableTLS *bool
	TLSCertFile *string
	TLSKeyFile *string
}

// Server represents an HTTP Server that serves static files and JSON api.
type Server struct {
	c    *ServerConfig
	l    *log.Logger
	srv  *http.Server
	mux  *http.ServeMux
	api  *Gateway
	quit chan struct{}
}

// NewServer creates a new server.
func NewServer(config *ServerConfig, api *Gateway) (*Server, error) {
	s := &Server{
		c:    config,
		l:    inform.NewLogger(true, os.Stdout, logPrefix),
		mux:  http.NewServeMux(),
		api:  api,
		quit: make(chan struct{}),
	}
	var e error
	if *config.StaticDir, e = filepath.Abs(*config.StaticDir); e != nil {
		return nil, e
	}

	if e := s.prepareMux(); e != nil {
		return nil, e
	}
	go s.serve()
	return s, nil
}

func (s *Server) serve() {
	s.srv = &http.Server{
		Addr: fmt.Sprintf(":%d", *s.c.Port),
		Handler: s.mux,
	}
	if *s.c.EnableTLS {
		for {
			if e := s.srv.ListenAndServeTLS(*s.c.TLSCertFile, *s.c.TLSKeyFile); e != nil {
				s.l.Println("stopped with error:", e)
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}
	} else {
		for {
			if e := s.srv.ListenAndServe(); e != nil {
				s.l.Println("stopped with error:", e)
				time.Sleep(100 * time.Millisecond)
				continue
			} else {
				break
			}
		}
	}

	s.srv = nil
}

func (s *Server) prepareMux() error {
	if *s.c.EnableGUI {
		if e := s.prepareStatic(); e != nil {
			return e
		}
	}
	return s.api.host(s.mux)
}

func (s *Server) prepareStatic() error {
	appLoc := *s.c.StaticDir
	s.mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		page := path.Join(appLoc, indexFileName)
		http.ServeFile(w, r, page)
	})

	fInfos, _ := ioutil.ReadDir(appLoc)
	for _, fInfo := range fInfos {
		route := fmt.Sprintf("/%s", fInfo.Name())
		if fInfo.IsDir() {
			route += "/"
		}
		s.mux.Handle(route, http.FileServer(http.Dir(appLoc)))
	}
	return nil
}

// CXO obtains the CXO.
func (s *Server) CXO() *cxo.Manager {
	return s.api.Access.CXO
}

// Close quits the http server.
func (s *Server) Close() {
	if s.quit != nil {
		s.CXO().Close()
		close(s.quit)
		if s.srv != nil {
			s.srv.Close()
		}
	}
}
