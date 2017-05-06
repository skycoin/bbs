package gui

import (
	"encoding/json"
	"github.com/evanlinjin/bbs/cxo"
	"net/http"
)

const (
	QueryBoard  = "board"
	QueryThread = "thread"
	QuerySeed   = "seed"
	QueryName   = "name"
	QueryDesc   = "desc"
)

func RegisterApiHandlers(mux *http.ServeMux, g *cxo.Gateway) {
	h := NewAPIHandler(g)

	mux.HandleFunc("/api/subscribe_to_board", h.SubscribeToBoard)
	mux.HandleFunc("/api/unsubscribe_from_board", h.UnSubscribeToBoard)

	mux.HandleFunc("/api/list_boards", h.ListBoards)
	mux.HandleFunc("/api/new_board", h.NewBoard)

	mux.HandleFunc("/api/list_threads", h.ListThreads)
	mux.HandleFunc("/api/new_thread", h.NewThread)

	mux.HandleFunc("/api/list_posts", h.ListPosts)
	mux.HandleFunc("/api/new_post", h.NewPost)
}

type APIHandler struct {
	g *cxo.Gateway
}

func NewAPIHandler(g *cxo.Gateway) *APIHandler {
	return &APIHandler{g: g}
}

// SubscribeToBoard handles subscription to board.
// Example usage: http://127.0.0.1:6420/api/subscribe_to_board?board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b
func (h *APIHandler) SubscribeToBoard(w http.ResponseWriter, r *http.Request) {
	board := r.URL.Query().Get(QueryBoard)
	reply := h.g.Subscribe(board)
	sendResponse(w, reply, http.StatusOK)
}

// UnSubscribeToBoard handles unsubscription to board.
// Example usage: http://127.0.0.1:6420/api/unsubscribe_from_board?board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b
func (h *APIHandler) UnSubscribeToBoard(w http.ResponseWriter, r *http.Request) {
	board := r.URL.Query().Get(QueryBoard)
	reply := h.g.Unsubscribe(board)
	sendResponse(w, reply, http.StatusOK)
}

// ListBoards lists all boards currently subscribed.
// Example usage: http://127.0.0.1:6420/api/list_boards
func (h *APIHandler) ListBoards(w http.ResponseWriter, r *http.Request) {
	reply := h.g.ViewBoards()
	sendResponse(w, reply, http.StatusOK)
}

// NewBoard creates a new board with a seed.
// Example usage: http://127.0.0.1:6420/api/new_board?seed=test
func (h *APIHandler) NewBoard(w http.ResponseWriter, r *http.Request) {
	seed := r.URL.Query().Get(QuerySeed)
	reply := h.g.NewBoard("", seed)
	sendResponse(w, reply, http.StatusOK)
}

// ListThreads lists all the threads of specified board.
// Example usage: http://127.0.0.1:6420/api/list_threads?board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b
func (h *APIHandler) ListThreads(w http.ResponseWriter, r *http.Request) {
	board := r.URL.Query().Get(QueryBoard)
	reply := h.g.ViewBoard(board)
	sendResponse(w, reply, http.StatusOK)
}

// NewThread creates a new thread under specified board.
// Example usage: http://127.0.0.1:6420/api/new_thread?board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b&name=Thread0&desc=hello
func (h *APIHandler) NewThread(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	reply := h.g.NewThread(q.Get(QueryBoard), q.Get(QueryName), q.Get(QueryDesc))
	sendResponse(w, reply, http.StatusOK)
}

// ListPosts lists all the posts of specified thread.
// Example usage: http://127.0.0.1:6420/api/list_posts?board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b
func (h *APIHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	reply := h.g.ViewThread(q.Get(QueryBoard), q.Get(QueryThread))
	sendResponse(w, reply, http.StatusOK)
}

func (h *APIHandler) NewPost(w http.ResponseWriter, r *http.Request) {
	sendTempResponse(w)
}

// Helper functions.

func setHeaderJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func sendTempResponse(w http.ResponseWriter) {
	setHeaderJSON(w)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func sendResponse(w http.ResponseWriter, v interface{}, httpStatus int) error {
	setHeaderJSON(w)
	respData, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.WriteHeader(httpStatus)
	w.Write(respData)
	return nil
}
