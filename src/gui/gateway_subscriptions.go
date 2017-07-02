package gui

import (
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
)

// Subscriptions represents the subscriptions endpoint group.
type Subscriptions struct {
	*Gateway
}

// GetAll gets all subscriptions.
func (g *Subscriptions) GetAll(w http.ResponseWriter, r *http.Request) {
	send(w, g.getAll(), http.StatusOK)
}

func (g *Subscriptions) getAll() []store.BoardInfo {
	return g.boardSaver.List()
}

// Get gets a subscription.
func (g *Subscriptions) Get(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	bi, has := g.get(bpk)
	if has {
		send(w, bi, http.StatusOK)
	} else {
		send(w, nil, http.StatusNotFound)
	}
}

func (g *Subscriptions) get(bpk cipher.PubKey) (store.BoardInfo, bool) {
	return g.boardSaver.Get(bpk)
}

// Add subscribes to a board.
func (g *Subscriptions) Add(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	g.add(r.FormValue("address"), bpk)
	send(w, true, http.StatusOK)
}

func (g *Subscriptions) add(addr string, bpk cipher.PubKey) {
	g.boardSaver.Add(addr, bpk)
}

// Remove unsubscribes from a board.
func (g *Subscriptions) Remove(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	g.remove(bpk)
	send(w, true, http.StatusOK)
}

func (g *Subscriptions) remove(bpk cipher.PubKey) {
	g.boardSaver.Remove(bpk)
}
