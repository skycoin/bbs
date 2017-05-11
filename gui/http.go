package gui

import (
	"net"
	"net/http"
	"github.com/skycoin/skycoin/src/util"
)

var (
	listener net.Listener
	quit     chan struct{}
)

const (
	guiDir = "./gui/static"
	resourceDir = "app/"
	devDir = "dev/"
	//indexPage   = "index.html"
)

func LaunchWebInterface(host string, g *Gateway) (e error) {
	quit = make(chan struct{})

	appLoc, e := util.DetermineResourcePath(guiDir, resourceDir, devDir)
	if e != nil {
		return
	}

	listener, e = net.Listen("tcp", host)
	if e != nil {
		return
	}
	serve(listener, NewServeMux(g, appLoc), quit)
	return
}

func serve(listener net.Listener, mux *http.ServeMux, q chan struct{}) {
	go func() {
		for {
			if e := http.Serve(listener, mux); e != nil {
				select {
				case <-q:
					return
				default:
				}
				continue
			}
		}
	}()
}

// Shutdown closes the http service.
func Shutdown() {
	if quit != nil {
		// must close quit first
		close(quit)
		listener.Close()
		listener = nil
	}
}

// NewServeMux creates a http.ServeMux with handlers registered.
func NewServeMux(g *Gateway, appLoc string) *http.ServeMux {
	// Register objects.
	jsonAPI := NewAPI(g)

	// Prepare mux.
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(appLoc)))

	mux.HandleFunc("/gui/boards", jsonAPI.BoardListHandler)
	mux.HandleFunc("/gui/boards/", jsonAPI.BoardHandler)
	return mux
}