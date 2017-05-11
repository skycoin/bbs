package gui

import (
	"net"
	"net/http"
	"path/filepath"
	wh "github.com/skycoin/skycoin/src/util/http" //http,json helpers
	"fmt"
	"io/ioutil"
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
	indexPage   = "index.html"
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

	mux.HandleFunc("/", newIndexHandler(appLoc))

	fileInfos, _ := ioutil.ReadDir(appLoc)
	for _, fileInfo := range fileInfos {
		route := fmt.Sprintf("/%s", fileInfo.Name())
		if fileInfo.IsDir() {
			route = route + "/"
		}
		mux.Handle(route, http.FileServer(http.Dir(appLoc)))
	}

	mux.HandleFunc("/gui/boards", jsonAPI.BoardListHandler)
	mux.HandleFunc("/gui/boards/", jsonAPI.BoardHandler)
	return mux
}

// Returns a http.HandlerFunc for index.html, where index.html is in appLoc
func newIndexHandler(appLoc string) http.HandlerFunc {
	// Serves the main page
	return func(w http.ResponseWriter, r *http.Request) {
		page := filepath.Join(appLoc, indexPage)
		fmt.Printf("Serving index page: %s\n", page)
		if r.URL.Path == "/" {
			http.ServeFile(w, r, page)
		} else {
			wh.Error404(w)
		}
	}
}