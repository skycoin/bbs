package gui

import (
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
)

// Users represents the users endpoint group.
type Users struct {
	*Gateway
	Masters MasterUsers
}

// GetAll gets all users, master or not.
func (g *Users) GetAll(w http.ResponseWriter, r *http.Request) {
	send(w, g.getAll(), http.StatusOK)
}

func (g *Users) getAll() []store.UserConfig {
	return g.userSaver.List()
}

// Add adds a user configuration for user we are not master of.
func (g *Users) Add(w http.ResponseWriter, r *http.Request) {
	// Get user public key.
	upk, e := misc.GetPubKey(r.FormValue("user"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get alias.
	alias := r.FormValue("alias")
	uc := g.add(alias, upk)
	send(w, uc, http.StatusOK)
}

func (g *Users) add(alias string, upk cipher.PubKey) store.UserConfig {
	g.userSaver.Add(alias, upk)
	uc, _ := g.userSaver.Get(upk)
	return uc
}

// Remove removes a user configuration, master or not.
func (g *Users) Remove(w http.ResponseWriter, r *http.Request) {
	// Get user public key.
	upk, e := misc.GetPubKey(r.FormValue("user"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	if e := g.remove(upk); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *Users) remove(upk cipher.PubKey) error {
	return g.userSaver.Remove(upk)
}

// MasterUsers represents the master users endpoint group.
type MasterUsers struct {
	*Gateway
	Current CurrentMasterUser
}

// GetAll obtains all users in which we are master of.
func (g *MasterUsers) GetAll(w http.ResponseWriter, r *http.Request) {
	send(w, g.getAll(), http.StatusOK)
}

func (g *MasterUsers) getAll() []store.UserConfig {
	return g.userSaver.ListMasters()
}

// Add adds a new user configuration of a master user.
func (g *MasterUsers) Add(w http.ResponseWriter, r *http.Request) {
	// Get alias and seed.
	uc := g.add(
		r.FormValue("alias"),
		r.FormValue("seed"),
	)
	send(w, uc, http.StatusOK)
}

func (g *MasterUsers) add(alias, seed string) store.UserConfig {
	pk, sk := cipher.GenerateDeterministicKeyPair([]byte(seed))
	g.userSaver.MasterAdd(alias, pk, sk)
	uc, _ := g.userSaver.Get(pk)
	return uc
}

// CurrentMasterUser represents the current master user.
type CurrentMasterUser struct {
	*Gateway
}

// Get returns the currently active user.
func (g *CurrentMasterUser) Get(w http.ResponseWriter, r *http.Request) {
	send(w, g.get(), http.StatusOK)
}

func (g *CurrentMasterUser) get() store.UserConfig {
	return g.userSaver.GetCurrent()
}

// Set sets the currently active user.
func (g *CurrentMasterUser) Set(w http.ResponseWriter, r *http.Request) {
	// Get user public key.
	upk, e := misc.GetPubKey(r.FormValue("user"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Set current user.
	if e := g.set(upk); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, g.get(), http.StatusOK)
}

func (g *CurrentMasterUser) set(upk cipher.PubKey) error {
	return g.userSaver.SetCurrent(upk)
}
