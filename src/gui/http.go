package gui

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	listener net.Listener
	quit     chan struct{}
)

// HTTPConfig contains configurations for HTTP Server.
type HTTPConfig struct {
	RPCRemoteAddr string
	Port          int
	StaticDir     string
	EnableGUI     bool
}

func OpenWebInterface(config *HTTPConfig, g *Gateway) (string, error) {
	// Get host.
	host := fmt.Sprintf("127.0.0.1:%d", config.Port)

	quit = make(chan struct{})
	appLoc, e := filepath.Abs(config.StaticDir)
	if e != nil {
		return "", e
	}

	listener, e = net.Listen("tcp", host)
	if e != nil {
		return "", e
	}
	go serve(listener, NewServeMux(g, appLoc, config.EnableGUI), quit)
	return fmt.Sprintf("%s://%s", "http", host), nil
}

func serve(listener net.Listener, mux *http.ServeMux, q chan struct{}) {
	if e := http.Serve(listener, mux); e != nil {
		select {
		case <-q:
			return
		default:
			log.Panic(e)
		}
	}
}

// Allows serving Angular.JS content swiftly.
func fileServe(mux *http.ServeMux, appLoc string) error {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, e := ioutil.ReadFile(path.Join(appLoc, "index.html"))
		if e != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(e.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	return filepath.Walk(appLoc, func(path string, info os.FileInfo, err error) error {
		// Skip directories.
		if info.IsDir() {
			return nil
		}
		httpPath := strings.TrimPrefix(path, appLoc)
		log.Printf("[WEBGUI] Found path: '%s'.", httpPath)
		mux.HandleFunc(httpPath, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, path)
		})
		return nil
	})
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
func NewServeMux(api *Gateway, appLoc string, enableGUI bool) *http.ServeMux {
	// Prepare mux.
	mux := http.NewServeMux()

	if enableGUI {
		fileServe(mux, appLoc)
	}

	mux.HandleFunc("/api/quit", api.Quit)
	mux.HandleFunc("/api/ping_submission_address", api.PingSubmissionAddress)
	mux.HandleFunc("/api/generate_seed", api.GenerateSeed)

	mux.HandleFunc("/api/stats/get", api.Stats.Get)

	mux.HandleFunc("/api/connections/get_all", api.Connections.GetAll)
	mux.HandleFunc("/api/connections/add", api.Connections.Add)
	mux.HandleFunc("/api/connections/remove", api.Connections.Remove)

	mux.HandleFunc("/api/subscriptions/get_all", api.Subscriptions.GetAll)
	mux.HandleFunc("/api/subscriptions/get", api.Subscriptions.Get)
	mux.HandleFunc("/api/subscriptions/add", api.Subscriptions.Add)
	mux.HandleFunc("/api/subscriptions/remove", api.Subscriptions.Remove)

	mux.HandleFunc("/api/users/get_all", api.Users.GetAll)
	mux.HandleFunc("/api/users/add", api.Users.Add)
	mux.HandleFunc("/api/users/remove", api.Users.Remove)
	mux.HandleFunc("/api/users/masters/get_all", api.Users.Masters.GetAll)
	mux.HandleFunc("/api/users/masters/add", api.Users.Masters.Add)
	mux.HandleFunc("/api/users/masters/current/get", api.Users.Masters.Current.Get)
	mux.HandleFunc("/api/users/masters/current/set", api.Users.Masters.Current.Set)

	mux.HandleFunc("/api/boards/get_all", api.Boards.GetAll)
	mux.HandleFunc("/api/boards/get", api.Boards.Get)
	mux.HandleFunc("/api/boards/add", api.Boards.Add)
	mux.HandleFunc("/api/boards/remove", api.Boards.Remove)
	mux.HandleFunc("/api/boards/meta/get", api.Boards.Meta.Get)
	mux.HandleFunc("/api/boards/meta/submission_addresses/get_all", api.Boards.Meta.SubmissionAddresses.GetAll)
	mux.HandleFunc("/api/boards/meta/submission_addresses/add", api.Boards.Meta.SubmissionAddresses.Add)
	mux.HandleFunc("/api/boards/meta/submission_addresses/remove", api.Boards.Meta.SubmissionAddresses.Remove)
	mux.HandleFunc("/api/boards/page/get", api.Boards.Page.Get)

	mux.HandleFunc("/api/threads/get_all", api.Threads.GetAll)
	mux.HandleFunc("/api/threads/add", api.Threads.Add)
	mux.HandleFunc("/api/threads/remove", api.Threads.Remove)
	mux.HandleFunc("/api/threads/import", api.Threads.Import)
	mux.HandleFunc("/api/threads/page/get", api.Threads.Page.Get)
	mux.HandleFunc("/api/threads/votes/get", api.Threads.Votes.Get)
	mux.HandleFunc("/api/threads/votes/add", api.Threads.Votes.Add)

	mux.HandleFunc("/api/posts/get_all", api.Posts.Get)
	mux.HandleFunc("/api/posts/add", api.Posts.Add)
	mux.HandleFunc("/api/posts/remove", api.Posts.Remove)
	mux.HandleFunc("/api/posts/votes/get", api.Posts.Votes.Get)
	mux.HandleFunc("/api/posts/votes/add", api.Posts.Votes.Add)

	//mux.HandleFunc("/api/hex/get_thread_page", api.GetThreadPageAsHex)
	//mux.HandleFunc("/api/hex/get_thread_page/tp_ref", api.GetThreadPageWithTpRefAsHex)
	//mux.HandleFunc("/api/hex/add_thread", api.NewThreadWithHex)
	//mux.HandleFunc("/api/hex/add_post", api.NewPostWithHex)

	mux.HandleFunc("/api/tests/add_filled_board", api.Tests.AddFilledBoard)
	mux.HandleFunc("/api/tests/panic", api.Tests.Panic)

	return mux
}
