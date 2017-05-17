package gui

import (
	"github.com/evanlinjin/bbs/extern"
	"github.com/skycoin/skycoin/src/util"
	"net"
	"net/http"
)

var (
	listener net.Listener
	quit     chan struct{}
)

const (
	guiDir      = "./extern/gui/static"
	resourceDir = "app/"
	devDir      = "dev/"
)

func OpenWebInterface(host string, g *extern.Gateway) (e error) {
	quit = make(chan struct{})

	appLoc, e := util.DetermineResourcePath(guiDir, resourceDir, devDir)
	if e != nil {
		return
	}

	listener, e = net.Listen("tcp", host)
	if e != nil {
		return
	}
	go serve(listener, NewServeMux(g, appLoc), quit)
	return
}

func serve(listener net.Listener, mux *http.ServeMux, q chan struct{}) {
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
}

// Close closes the http service.
func Close() {
	if quit != nil {
		// must close quit first
		close(quit)
		listener.Close()
		listener = nil
	}
}

// NewServeMux creates a http.ServeMux with handlers registered.
func NewServeMux(g *extern.Gateway, appLoc string) *http.ServeMux {
	// Register objects.
	api := NewAPI(g)

	// Prepare mux.
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(appLoc)))

	mux.HandleFunc("/api/get_subscription", api.GetSubscription)
	mux.HandleFunc("/api/get_subscriptions", api.GetSubscriptions)
	mux.HandleFunc("/api/subscribe", api.Subscribe)
	mux.HandleFunc("/api/unsubscribe", api.Unsubscribe)

	mux.HandleFunc("/api/users/get_current", api.GetCurrentUser)
	mux.HandleFunc("/api/users/set_current", api.SetCurrentUser)
	mux.HandleFunc("/api/users/get_masters", api.GetMasterUsers)
	mux.HandleFunc("/api/users/new_master", api.NewMasterUser)
	mux.HandleFunc("/api/users/get_all", api.GetUsers)
	mux.HandleFunc("/api/users/new", api.NewUser)
	mux.HandleFunc("/api/users/remove", api.RemoveUser)

	mux.HandleFunc("/api/get_boards", api.GetBoards)
	mux.HandleFunc("/api/new_board", api.NewBoard)
	mux.HandleFunc("/api/get_threads", api.GetThreads)
	mux.HandleFunc("/api/new_thread", api.NewThread)
	mux.HandleFunc("/api/get_posts", api.GetPosts)
	mux.HandleFunc("/api/new_post", api.NewPost)

	mux.HandleFunc("/api/tests/new_filled_board", api.TestNewFilledBoard)

	return mux
}
