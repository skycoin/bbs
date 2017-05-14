package gui

import (
	json "encoding/json"
	"github.com/evanlinjin/bbs/typ"
	"io/ioutil"
	"net/http"
	"strings"
)

// API wraps cxo.Gateway.
type API struct {
	g *Gateway
}

// NewAPI creates a new API.
func NewAPI(g *Gateway) *API {
	return &API{g}
}

// Stat handles the "stat" endpoint.
// It shows statistics of Skycoin BBS.
func (a *API) Stat(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		reply := a.g.Stat()
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// Subscriptions handles the "subscriptions" endpoint.
// Functions include: list subscriptions, subscribe to board, unsubscribe from board,
// check if subscribed to a board.
func (a *API) Subscriptions(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.EscapedPath(), "/")
	switch len(path) {
	case 3:
		// "/api/subscriptions" >>>
		switch r.Method {
		case "GET":
			// List subscriptions >>>
			reply := a.g.ListSubscriptions()
			sendResponse(w, reply, http.StatusOK)
			return
		}
	case 4:
		// "/api/subscriptions/:id" >>>
		pk := path[3]
		switch r.Method {
		case "GET":
			// Check subscription >>>
			switch a.g.CheckSubscription(pk) {
			case true:
				sendResponse(w, nil, http.StatusOK)
				return
			case false:
				sendResponse(w, nil, http.StatusNotFound)
				return
			}
		case "POST":
			// Subscribe >>>
			reply := a.g.Subscribe(pk)
			sendResponse(w, reply, http.StatusOK)
			return
		case "DELETE":
			// Unsubscribe >>>
			reply := a.g.Unsubscribe(pk)
			sendResponse(w, reply, http.StatusOK)
			return
		}
	}
	sendResponse(w, nil, http.StatusBadRequest)
	return
}

// Boards handles the "boards" endpoint.
// Functions include: list all boards, create a board, remove a board, list all threads of a board,
// create new a new thread on a board, remove a thread from a board.
func (a *API) Boards(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.EscapedPath(), "/")
	switch len(path) {
	case 3:
		// "/api/boards" >>>
		switch r.Method {
		case "GET":
			// List all boards >>>
			reply := a.g.ListBoards()
			sendResponse(w, reply, http.StatusOK)
			return
		case "POST":
			// Create a new board >>>
			req, e := readRequestBody(r)
			if e != nil || req.Board == nil {
				sendResponse(w, "invalid request body", http.StatusNotAcceptable)
				return
			}
			reply := a.g.NewBoard(req.Board, req.Seed)
			sendResponse(w, reply, http.StatusOK)
			return
		case "DELETE":
			// Delete a board >>>
			sendResponse(w, nil, http.StatusNotImplemented)
			return
		}
	case 4:
		// "/api/boards/:board_id" >>>
		//pk := path[3]
		switch r.Method {
		case "DELETE":
			// Remove board >>>
			sendResponse(w, nil, http.StatusNotImplemented)
			return
		}
	case 5:
		// "/api/boards/:board_id/[...]" >>>
		pk := path[3]
		switch path[4] {
		case "threads":
			// "/api/boards/:board_id/threads" >>>
			switch r.Method {
			case "GET":
				// List all threads of specified board >>>
				reply := a.g.ViewBoard(pk)
				sendResponse(w, reply, http.StatusOK)
				return
			case "POST":
				// Add a thread to specified board >>>
				req, e := readRequestBody(r)
				if e != nil || req.Thread == nil {
					sendResponse(w,
						"invalid request body, no thread specified",
						http.StatusNotAcceptable)
					return
				}
				reply := a.g.NewThread(pk, req.Thread)
				sendResponse(w, reply, http.StatusOK)
				return
			}
		}
	case 6:
		// "/api/boards/[...]/[...]/[...]" >>>
		switch path[4] {
		case "threads":
			// "/api/boards/:board_id/threads/:thread_id" >>>
			switch r.Method {
			case "DELETE":
				// Remove thread from board >>>
				sendResponse(w, nil, http.StatusNotImplemented)
				return
			}
		}
	}
	sendResponse(w, nil, http.StatusBadRequest)
	return
}

// Threads handles the "threads" endpoint.
// Functions include: listing posts of a thread, creating a new post in specified thread.
func (a *API) Threads(w http.ResponseWriter, r *http.Request) {
	path := strings.Split(r.URL.EscapedPath(), "/")
	switch len(path) {
	case 5:
		// "/api/threads/[...]/[...]" >>
		switch {
		case path[2] == "threads" && path[4] == "posts":
			// "/api/threads/[:thread_id]/posts" >>>
			switch r.Method {
			case "GET":
				// List all posts of specified thread >>>
				sendResponse(w, nil, http.StatusNotImplemented)
				return
			case "POST":
				// Create a new post on specified thread >>>
				sendResponse(w, nil, http.StatusNotImplemented)
				return
			}

		}
	}
	sendResponse(w, nil, http.StatusBadRequest)
	return
}

// Subscribe handles the "subscribe" endpoint.
// It subscribes to a board.
func (a *API) Subscribe(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req, e := readRequestBody(r)
		if e != nil || req.Board == nil {
			sendResponse(w, "invalid request body", http.StatusNotAcceptable)
			return
		}
		reply := a.g.Subscribe(req.Board.PubKey)
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// Unsubscribe handles the "unsubscribe" endpoint.
// It unsubscribes to a board.
func (a *API) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req, e := readRequestBody(r)
		if e != nil || req.Board == nil {
			sendResponse(w, "invalid request body", http.StatusNotAcceptable)
			return
		}
		reply := a.g.Unsubscribe(req.Board.PubKey)
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// ListBoards handles the "list_boards" endpoint.
// It returns a list of boards which we are subscribed to.
func (a *API) ListBoards(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		reply := a.g.ListBoards()
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// NewBoard handles the "new_board" endpoint.
// It creates a new board if BBS server is 'master'.
func (a *API) NewBoard(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req, e := readRequestBody(r)
		if e != nil || req.Board == nil {
			sendResponse(w, "invalid request body", http.StatusNotAcceptable)
			return
		}
		reply := a.g.NewBoard(req.Board, req.Seed)
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// ListThreads handles the "list_threads" endpoint.
// It lists the threads of a specified board.
func (a *API) ListThreads(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req, e := readRequestBody(r)
		if e != nil || req.Board == nil {
			sendResponse(w, "invalid request body", http.StatusNotAcceptable)
			return
		}
		reply := a.g.ViewBoard(req.Board.PubKey)
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// NewThread handles the "new_thread" endpoint.
// It creates a new thread on a specified board.
func (a *API) NewThread(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req, e := readRequestBody(r)
		if e != nil || req.Board == nil || req.Thread == nil {
			sendResponse(w, "invalid request body", http.StatusNotAcceptable)
			return
		}
		reply := a.g.NewThread(req.Board.PubKey, req.Thread)
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// ListPosts handles the "list_posts" endpoint.
// It lists the posts of a specified board and thread.
func (a *API) ListPosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req, e := readRequestBody(r)
		if e != nil || req.Board == nil || req.Thread == nil {
			sendResponse(w, "invalid request body", http.StatusNotAcceptable)
			return
		}
		reply := a.g.ViewThread(req.Board.PubKey, req.Thread.Hash)
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// NewPost handles the "new_post" endpoint.
// It creates a new post on a specified board and thread.
func (a *API) NewPost(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		req, e := readRequestBody(r)
		if e != nil || req.Board == nil || req.Thread == nil || req.Post == nil {
			sendResponse(w, "invalid request body", http.StatusNotAcceptable)
			return
		}
		reply := a.g.NewPost(req.Board.PubKey, req.Thread.Hash, req.Post)
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// BoardListHandler for /gui/boards.
func (a *API) BoardListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		reply := a.g.ListBoards()
		sendResponse(w, reply, http.StatusOK)
		return
	case "PUT":
		req, e := readRequestBody(r)
		if e != nil || req.Board == nil {
			sendResponse(w, "invalid request body", http.StatusNotAcceptable)
			return
		}
		reply := a.g.NewBoard(req.Board, req.Seed)
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// BoardHandler for /gui/boards/BOARD_PUBLIC_KEY.
func (a *API) BoardHandler(w http.ResponseWriter, r *http.Request) {
	// Obtain path.
	path := strings.Split(r.URL.EscapedPath(), "/")
	// Obtain public key.
	pkStr := path[3]
	// If it's view board, or view thread.
	switch len(path) {
	case 4:
		// View Board.
		switch r.Method {
		case "GET":
			reply := a.g.ViewBoard(pkStr)
			sendResponse(w, reply, http.StatusOK)
			return
		case "PUT":
			req, e := readRequestBody(r)
			if e != nil || req.Thread == nil {
				sendResponse(w, "invalid request body", http.StatusNotAcceptable)
				return
			}
			reply := a.g.NewThread(pkStr, req.Thread)
			sendResponse(w, reply, http.StatusOK)
			return
		}
	case 5:
		// View Thread.
		tHashStr := path[4]
		switch r.Method {
		case "GET":
			reply := a.g.ViewThread(pkStr, tHashStr)
			sendResponse(w, reply, http.StatusOK)
			return
		case "PUT":
			req, e := readRequestBody(r)
			if e != nil || req.Post == nil {
				sendResponse(w, "invalid request body", http.StatusNotAcceptable)
				return
			}
			reply := a.g.NewPost(pkStr, tHashStr, req.Post)
			sendResponse(w, reply, http.StatusOK)
			return
		}
		//sendResponse(w, tHashStr, http.StatusNotImplemented)
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// Helper functions.
func sendResponse(w http.ResponseWriter, v interface{}, httpStatus int) error {
	w.Header().Set("Content-Type", "application/json")
	respData, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.WriteHeader(httpStatus)
	w.Write(respData)
	return nil
}

func readRequestBody(r *http.Request) (*typ.ReqRep, error) {
	d, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, e
	}
	obj := typ.NewRepReq()
	e = json.Unmarshal(d, obj)
	return obj, e
}
