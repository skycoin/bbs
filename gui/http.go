package gui

import (
	"github.com/skycoin/skycoin/src/util"
	"net"
	"net/http"
)

var (
	listener net.Listener
	quit     chan struct{}
)

const (
	guiDir      = "./gui/static"
	resourceDir = "app/"
	devDir      = "dev/"
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
	api := NewAPI(g)

	// Prepare mux.
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(appLoc)))

	// Current (v3).
	mux.HandleFunc("/api/subscriptions", api.Subscriptions)
	mux.HandleFunc("/api/subscriptions/", api.Subscriptions)
	mux.HandleFunc("/api/boards", api.Boards)
	mux.HandleFunc("/api/boards/", api.Boards)
	mux.HandleFunc("/api/threads", api.Threads)
	mux.HandleFunc("/api/threads/", api.Threads)

	// Deprecated (v2).
	mux.HandleFunc("/api_v2/stat", api.Stat)
	mux.HandleFunc("/api_v2/subscribe", api.Subscribe)
	mux.HandleFunc("/api_v2/unsubscribe", api.Unsubscribe)
	mux.HandleFunc("/api_v2/list_boards", api.ListBoards)
	mux.HandleFunc("/api_v2/new_board", api.NewBoard)
	mux.HandleFunc("/api_v2/list_threads", api.ListThreads)
	mux.HandleFunc("/api_v2/new_thread", api.NewThread)
	mux.HandleFunc("/api_v2/list_posts", api.ListPosts)
	mux.HandleFunc("/api_v2/new_post", api.NewPost)

	// Deprecated (v1).
	mux.HandleFunc("/api_v1/boards", api.BoardListHandler)
	mux.HandleFunc("/api_v1/boards/", api.BoardHandler)

	return mux
}
