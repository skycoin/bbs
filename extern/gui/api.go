package gui

import (
	"encoding/json"
	"fmt"
	"github.com/evanlinjin/bbs/extern"
	"github.com/evanlinjin/bbs/store/misc"
	"net/http"
	//"log"
	"github.com/evanlinjin/bbs/store/typ"
)

// API wraps cxo.Gateway.
type API struct {
	g *extern.Gateway
}

// NewAPI creates a new API.
func NewAPI(g *extern.Gateway) *API {
	return &API{g}
}

/*
	<<< FOR SUBSCRIPTIONS >>>
*/

func (a *API) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetSubscriptions(), http.StatusOK)
}

func (a *API) GetSubscription(w http.ResponseWriter, r *http.Request) {
	bpkStr := r.FormValue("board")
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(w, fmt.Sprintln(e), http.StatusBadRequest)
		return
	}
	bi, has := a.g.GetSubscription(bpk)
	if has {
		sendResponse(w, bi, http.StatusOK)
	} else {
		sendResponse(w, nil, http.StatusNotFound)
	}
}

func (a *API) Subscribe(w http.ResponseWriter, r *http.Request) {
	bpkStr := r.FormValue("board")
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	a.g.Subscribe(bpk)
	sendResponse(w, true, http.StatusOK)
}

func (a *API) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	bpkStr := r.FormValue("board")
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	a.g.Unsubscribe(bpk)
	sendResponse(w, true, http.StatusOK)
}

/*
	<<< FOR USERS >>>
*/

func (a *API) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetCurrentUser(), http.StatusOK)
}

func (a *API) SetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user public key.
	upkStr := r.FormValue("user")
	upk, e := misc.GetPubKey(upkStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	// Set current user.
	if e := a.g.SetCurrentUser(upk); e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	sendResponse(w, a.g.GetCurrentUser(), http.StatusOK)
}

func (a *API) GetMasterUsers(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetMasterUsers(), http.StatusOK)
}

func (a *API) NewMasterUser(w http.ResponseWriter, r *http.Request) {
	// Get alias and seed.
	alias := r.FormValue("alias")
	seed := r.FormValue("seed")
	uc := a.g.NewMasterUser(alias, seed)
	sendResponse(w, uc, http.StatusOK)
}

func (a *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetUsers(), http.StatusOK)
}

func (a *API) NewUser(w http.ResponseWriter, r *http.Request) {
	// Get user public key.
	upkStr := r.FormValue("user")
	upk, e := misc.GetPubKey(upkStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	// Get alias.
	alias := r.FormValue("alias")
	uc := a.g.NewUser(alias, upk)
	sendResponse(w, uc, http.StatusOK)
}

func (a *API) RemoveUser(w http.ResponseWriter, r *http.Request) {
	// Get user public key.
	upkStr := r.FormValue("user")
	upk, e := misc.GetPubKey(upkStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	if e := a.g.RemoveUser(upk); e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

/*
	<<< FOR BOARDS, THREADS & POSTS >>>
*/

func (a *API) GetBoards(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetBoards(), http.StatusOK)
}

func (a *API) NewBoard(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	desc := r.FormValue("description")
	seed := r.FormValue("seed")
	board := &typ.Board{Name: name, Desc: desc}
	bi, e := a.g.NewBoard(board, seed)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
	} else {
		sendResponse(w, bi, http.StatusOK)
	}
}

func (a *API) GetThreads(w http.ResponseWriter, r *http.Request) {
	bpkStr := r.FormValue("board")
	if bpkStr == "" {
		sendResponse(w, a.g.GetThreads(), http.StatusOK)
		return
	}
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	sendResponse(w, a.g.GetThreads(bpk), http.StatusOK)
}

func (a *API) NewThread(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpkStr := r.FormValue("board")
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	// Get thread values.
	name := r.FormValue("name")
	desc := r.FormValue("description")
	thread := &typ.Thread{Name: name, Desc: desc, MasterBoard: bpk.Hex()}
	if e := a.g.NewThread(bpk, thread); e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	sendResponse(w, thread, http.StatusOK)
}

func (a *API) GetPosts(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpkStr := r.FormValue("board")
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRefStr := r.FormValue("thread")
	tRef, e := misc.GetReference(tRefStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	posts, e := a.g.GetPosts(bpk, tRef)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	sendResponse(w, posts, http.StatusOK)
}

func (a *API) NewPost(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpkStr := r.FormValue("board")
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRefStr := r.FormValue("thread")
	tRef, e := misc.GetReference(tRefStr)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	// Get post values.
	title := r.FormValue("title")
	body := r.FormValue("body")
	post := &typ.Post{Title: title, Body: body}
	if e := a.g.NewPost(bpk, tRef, post); e != nil {
		sendResponse(w, e, http.StatusBadRequest)
		return
	}
	sendResponse(w, post, http.StatusOK)
}

/*
	<<< HELPER FUNCTIONS >>>
*/

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
