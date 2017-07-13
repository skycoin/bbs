package gui

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/skycoin/src/cipher"
	"net/http"
	"strconv"
)

// Users represents the users endpoint group.
type Users struct {
	*Gateway
	Masters MasterUsers
	Votes   UserVotes
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

// UserVotes represents the user votes endpoint group.
type UserVotes struct {
	*Gateway
}

// Get gets votes for user of specified board.
func (g *UserVotes) Get(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get user public key.
	upk, e := misc.GetPubKey(r.FormValue("user"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get posts.
	vv, e := g.get(bpk, upk)
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, vv, http.StatusOK)
}

func (g *UserVotes) get(bpk, upk cipher.PubKey) (*VotesView, error) {
	// Get current user.
	cu := g.userSaver.GetCurrent()
	cUPK := cu.GetPK()
	// Get votes.
	votes := g.container.GetVotesForUser(bpk, upk)
	vv := &VotesView{}
	for _, vote := range votes {
		switch vote.Mode {
		case +1:
			vv.UpVotes += 1
		case -1:
			vv.DownVotes += 1
		}
		if vote.User == cUPK {
			vv.CurrentUserVoted = true
			vv.CurrentUserVoteMode = int(vote.Mode)
		}
	}
	return vv, nil
}

// Add adds a vote for user of specified board.
func (g *UserVotes) Add(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get user public key.
	upk, e := misc.GetPubKey(r.FormValue("user"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get vote mode (up/down vote).
	mode, e := strconv.Atoi(r.FormValue("mode"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Prepare vote.
	vote := &typ.Vote{Mode: int8(mode), Tag: []byte(r.FormValue("tag"))}
	if e := g.add(bpk, upk, vote); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *UserVotes) add(bpk, upk cipher.PubKey, vote *typ.Vote) error {
	// Get current user.
	uc := g.userSaver.GetCurrent()
	// Check vote.
	if e := vote.Sign(uc.GetPK(), uc.GetSK()); e != nil {
		return errors.Wrap(e, "vote signing failed")
	}
	// Check board.
	bi, got := g.boardSaver.Get(bpk)
	if !got {
		return errors.Errorf("not subscribed to board '%s'", bpk.Hex())
	}
	// Check if this node owns the board.
	if bi.Config.Master {
		// Via CXO.
		switch vote.Mode {
		case 0:
			return g.container.RemoveVoteForUser(uc.GetPK(), bpk, upk, bi.Config.GetSK())
		case -1, +1:
			return g.container.AddVoteForUser(bpk, upk, bi.Config.GetSK(), vote)
		}
	} else {
		return g.queueSaver.AddVoteUserReq(bpk, upk, vote)
	}
	return nil
}
