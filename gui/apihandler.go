package gui

import (
	"encoding/json"
	"github.com/evanlinjin/bbs-server/datastore"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
)

const (
	QueryBoard       = "board"
	QuerySeed        = "seed"
	QueryTitle       = "title"
	QueryDescription = "description"
)

func RegisterApiHandlers(mux *http.ServeMux, cc *datastore.CXOClient) {
	h := NewAPIHandler(cc)

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
	cc *datastore.CXOClient
}

func NewAPIHandler(cc *datastore.CXOClient) *APIHandler {
	return &APIHandler{cc: cc}
}

// SubscribeToBoard handles subscription to board.
// Example usage: http://127.0.0.1:6420/api/subscribe_to_board?board=03517b80b2889e4de80aae0fa2a4b2a408490f3178857df5b756e690b4524e1e61
func (h *APIHandler) SubscribeToBoard(w http.ResponseWriter, r *http.Request) {
	pk, e := cipher.PubKeyFromHex(r.URL.Query().Get(QueryBoard))
	if e != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, e.Error())
		return
	}
	bc, e := h.cc.SubscribeToBoard(pk)
	if e != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, e.Error())
		return
	}
	sendResponse(w, bc, http.StatusOK)
}

// UnSubscribeToBoard handles unsubscription to board.
// Example usage: http://127.0.0.1:6420/api/unsubscribe_from_board?board=03517b80b2889e4de80aae0fa2a4b2a408490f3178857df5b756e690b4524e1e61
func (h *APIHandler) UnSubscribeToBoard(w http.ResponseWriter, r *http.Request) {
	pk, e := cipher.PubKeyFromHex(r.URL.Query().Get(QueryBoard))
	if e != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, e.Error())
		return
	}
	okay := h.cc.UnSubscribeFromBoard(pk)
	if okay == false {
		sendErrorResponse(w, http.StatusNotAcceptable, "unable to unsubscribe")
		return
	}
	sendResponse(w, "unsubscribed", http.StatusOK)
}

// ListBoards lists all boards currently subscribed.
// Example usage: http://127.0.0.1:6420/api/list_boards
func (h *APIHandler) ListBoards(w http.ResponseWriter, r *http.Request) {
	boards := h.cc.ListBoards()
	sendResponse(w, boards, http.StatusOK)
}

// NewBoard creates a new board with a seed.
// Example usage: http://127.0.0.1:6420/api/new_board?seed=test
func (h *APIHandler) NewBoard(w http.ResponseWriter, r *http.Request) {
	seed := r.URL.Query().Get(QuerySeed)
	if seed == "" {
		sendErrorResponse(w, http.StatusNotAcceptable, "seed not specified")
		return
	}
	bc, e := h.cc.NewBoard("", seed)
	if e != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, e.Error())
		return
	}
	sendResponse(w, bc, http.StatusOK)
}

// ListThreads lists all the threads of specified board.
// Example usage: http://127.0.0.1:6420/api/list_threads?board=032ffee44b9554cd3350ee16760688b2fb9d0faae7f3534917ff07e971eb36fd6b
func (h *APIHandler) ListThreads(w http.ResponseWriter, r *http.Request) {
	pk, e := cipher.PubKeyFromHex(r.URL.Query().Get(QueryBoard))
	if e != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, e.Error())
		return
	}
	bp, e := h.cc.ListThreads(pk)
	if e != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, e.Error())
		return
	}
	sendResponse(w, bp, http.StatusOK)
}

// NewThread creates a new thread under specified board.
// Example usage:
func (h *APIHandler) NewThread(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()
	pkBoard := qs.Get(QueryBoard)
	tdTitle := qs.Get(QueryTitle)
	tdDesc := qs.Get(QueryDescription)
	if len(tdTitle) < 1 {
		sendErrorResponse(w, http.StatusNotAcceptable, "thread title too short")
		return
	}
	pk, e := cipher.PubKeyFromHex(pkBoard)
	if e != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, e.Error())
		return
	}
	thread, e := h.cc.NewThread(pk, tdTitle, tdDesc)
	if e != nil {
		sendErrorResponse(w, http.StatusNotAcceptable, e.Error())
		return
	}
	sendResponse(w, thread, http.StatusOK)
}

func (h *APIHandler) ListPosts(w http.ResponseWriter, r *http.Request) {
	sendTempResponse(w)
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

func sendErrorResponse(w http.ResponseWriter, status int, msg string) {
	setHeaderJSON(w)
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + msg + `"}`))
}
