package gui

import (
	"encoding/json"
	"fmt"
	"github.com/evanlinjin/bbs/extern"
	"io/ioutil"
	"net/http"
	"github.com/evanlinjin/bbs/misc"
	//"log"
)

// API wraps cxo.Gateway.
type API struct {
	g *extern.Gateway
}

// NewAPI creates a new API.
func NewAPI(g *extern.Gateway) *API {
	return &API{g}
}

func (a *API) HelloWorld(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, "Hello, world!", http.StatusOK)
}

func (a *API) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetSubscriptions(), http.StatusOK)
}

func (a *API) GetSubscription(w http.ResponseWriter, r *http.Request) {
	bpkStr := r.FormValue("board")
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(
			w,
			fmt.Sprintln("invalid board public key provided:", e),
			http.StatusBadRequest,
		)
		return
	}
	bi, has := a.g.GetSubscription(bpk)
	if has {
		sendResponse(w, bi, http.StatusOK)
	} else {
		sendResponse(w, bi, http.StatusNotFound)
	}
}

func (a *API) Subscribe(w http.ResponseWriter, r *http.Request) {
	bpkStr := r.FormValue("board")
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(
			w,
			fmt.Sprintln("invalid board public key provided:", e),
			http.StatusBadRequest,
		)
		return
	}
	a.g.Subscribe(bpk)
	sendResponse(w, true, http.StatusOK)
}

func (a *API) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	bpkStr := r.FormValue("board")
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(
			w,
			fmt.Sprintln("invalid board public key provided:", e),
			http.StatusBadRequest,
		)
		return
	}
	a.g.Unsubscribe(bpk)
	sendResponse(w, true, http.StatusOK)
}

func (a *API) GetBoards(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetBoards(), http.StatusOK)
}

func (a *API) NewBoard(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	desc := r.FormValue("description")
	seed := r.FormValue("seed")
	bi, e := a.g.NewBoard(name, desc, seed)
	if e != nil {
		sendResponse(w, e, http.StatusBadRequest)
	} else {
		sendResponse(w, bi, http.StatusOK)
	}
}

/*
	Helper Functions.
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

func readRequestBody(r *http.Request) (*extern.ReqRes, error) {
	d, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, e
	}
	obj := extern.NewRepRes()
	e = json.Unmarshal(d, obj)
	return obj, e
}
