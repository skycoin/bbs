package gui

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store/typ"
	"net/http"
	"strconv"
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

/*
	<<< MISC >>>
*/

func (a *API) Quit(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.Quit(), http.StatusOK)
}

/*
	<<< FOR STATS >>>
*/

func (a *API) GetStats(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetStats(), http.StatusOK)
}

/*
	<<< FOR CONNECTIONS >>>
*/

func (a *API) GetConnections(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetConnections(), http.StatusOK)
}

func (a *API) AddConnection(w http.ResponseWriter, r *http.Request) {
	if e := a.g.AddConnection(r.FormValue("address")); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

func (a *API) RemoveConnection(w http.ResponseWriter, r *http.Request) {
	if e := a.g.RemoveConnection(r.FormValue("address")); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

/*
	<<< FOR SUBSCRIPTIONS >>>
*/

func (a *API) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetSubscriptions(), http.StatusOK)
}

func (a *API) GetSubscription(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
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
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	a.g.Subscribe(r.FormValue("address"), bpk)
	sendResponse(w, true, http.StatusOK)
}

func (a *API) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
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
	upk, e := misc.GetPubKey(r.FormValue("user"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Set current user.
	if e := a.g.SetCurrentUser(upk); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, a.g.GetCurrentUser(), http.StatusOK)
}

func (a *API) GetMasterUsers(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetMasterUsers(), http.StatusOK)
}

func (a *API) NewMasterUser(w http.ResponseWriter, r *http.Request) {
	// Get alias and seed.
	uc := a.g.NewMasterUser(
		r.FormValue("alias"),
		r.FormValue("seed"),
	)
	sendResponse(w, uc, http.StatusOK)
}

func (a *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, a.g.GetUsers(), http.StatusOK)
}

func (a *API) NewUser(w http.ResponseWriter, r *http.Request) {
	// Get user public key.
	upk, e := misc.GetPubKey(r.FormValue("user"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get alias.
	alias := r.FormValue("alias")
	uc := a.g.NewUser(alias, upk)
	sendResponse(w, uc, http.StatusOK)
}

func (a *API) RemoveUser(w http.ResponseWriter, r *http.Request) {
	// Get user public key.
	upk, e := misc.GetPubKey(r.FormValue("user"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	if e := a.g.RemoveUser(upk); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

/*
	<<< FOR BOARD META >>>
*/

func (a *API) GetSubmissionAddresses(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	addresses, e := a.g.GetSubmissionAddresses(bpk)
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, addresses, http.StatusOK)
}

func (a *API) AddSubmissionAddress(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	address := r.FormValue("address")
	if address == "" {
		sendResponse(w, "no submission address provided", http.StatusBadRequest)
		return
	}
	if e := a.g.AddSubmissionAddress(bpk, address); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

func (a *API) RemoveSubmissionAddress(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	address := r.FormValue("address")
	if address == "" {
		sendResponse(w, "no submission address provided", http.StatusBadRequest)
		return
	}
	if e := a.g.RemoveSubmissionAddress(bpk, address); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
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
	// Obtain Board Meta.
	meta := new(typ.BoardMeta)
	meta.SubmissionAddresses =
		strings.Split(r.FormValue("submission_addresses"), ",")
	meta.Trim()

	// Generate board.
	board := &typ.Board{
		Name: r.FormValue("name"),
		Desc: r.FormValue("description"),
	}
	board.SetMeta(meta)

	// Create board in cxo.
	bi, e := a.g.NewBoard(board, r.FormValue("seed"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
	} else {
		sendResponse(w, bi, http.StatusOK)
	}
}

func (a *API) RemoveBoard(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	e = a.g.RemoveBoard(bpk)
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
	} else {
		sendResponse(w, true, http.StatusOK)
	}
}

func (a *API) GetBoardPage(w http.ResponseWriter, r *http.Request) {
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	bpv, e := a.g.GetBoardPage(bpk)
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, bpv, http.StatusOK)
}

func (a *API) GetThreads(w http.ResponseWriter, r *http.Request) {
	bpkStr := r.FormValue("board")
	if bpkStr == "" {
		sendResponse(w, a.g.GetThreads(), http.StatusOK)
		return
	}
	bpk, e := misc.GetPubKey(bpkStr)
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, a.g.GetThreads(bpk), http.StatusOK)
}

func (a *API) NewThread(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread values.
	thread := &typ.Thread{
		Name:        r.FormValue("name"),
		Desc:        r.FormValue("description"),
		MasterBoard: bpk.Hex(),
	}
	if e := a.g.NewThread(bpk, thread); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, thread, http.StatusOK)
}

func (a *API) RemoveThread(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	if e := a.g.RemoveThread(bpk, tRef); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

func (a *API) GetThreadPage(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread page.
	threadPage, e := a.g.GetThreadPage(bpk, tRef)
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, threadPage, http.StatusOK)
}

func (a *API) GetPosts(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	posts, e := a.g.GetPosts(bpk, tRef)
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, posts, http.StatusOK)
}

func (a *API) NewPost(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get post values.
	post := &typ.Post{
		Title: r.FormValue("title"),
		Body:  r.FormValue("body"),
	}
	if e := a.g.NewPost(bpk, tRef, post); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, post, http.StatusOK)
}

func (a *API) RemovePost(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// // Get post reference.
	pRef, e := misc.GetReference(r.FormValue("post"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}

	if e := a.g.RemovePost(bpk, tRef, pRef); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

func (a *API) ImportThread(w http.ResponseWriter, r *http.Request) {
	// Get "from" board's public key.
	fromBpk, e := misc.GetPubKey(r.FormValue("from_board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread's reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get "to" board's public key.
	toBpk, e := misc.GetPubKey(r.FormValue("to_board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Import thread.
	if e := a.g.ImportThread(fromBpk, toBpk, tRef); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

/*
	<<< VOTES >>>
*/

func (a *API) GetVotesForThread(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get votes.
	vv, e := a.g.GetVotesForThread(bpk, tRef)
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, vv, http.StatusOK)
}

func (a *API) GetVotesForPost(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get post reference.
	pRef, e := misc.GetReference(r.FormValue("post"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get posts.
	vv, e := a.g.GetVotesForPost(bpk, pRef)
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, vv, http.StatusOK)
}

func (a *API) AddVoteForThread(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get vote mode (up/down vote).
	mode, e := strconv.Atoi(r.FormValue("mode"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Prepare vote.
	vote := &typ.Vote{Mode: int8(mode), Tag: []byte(r.FormValue("tag"))}
	if e := a.g.VoteForThread(bpk, tRef, vote); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

func (a *API) AddVoteForPost(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get post reference.
	pRef, e := misc.GetReference(r.FormValue("post"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get vote mode (up/down vote).
	mode, e := strconv.Atoi(r.FormValue("mode"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Prepare vote.
	vote := &typ.Vote{Mode: int8(mode), Tag: []byte(r.FormValue("tag"))}
	if e := a.g.VoteForPost(bpk, pRef, vote); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

/*
	<<< HEX >>>
*/

func (a *API) GetThreadPageAsHex(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread page as hex.
	tph, e := a.g.GetThreadPageAsHex(bpk, tRef)
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, *tph, http.StatusOK)
}

func (a *API) GetThreadPageWithTpRefAsHex(w http.ResponseWriter, r *http.Request) {
	// Get thread page reference.
	tpRef, e := misc.GetReference(r.FormValue("threadpage"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread page as hex.
	tph, e := a.g.GetThreadPageWithTpRefAsHex(tpRef)
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, *tph, http.StatusOK)
}

func (a *API) NewThreadWithHex(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread data.
	tData, e := misc.GetBytes(r.FormValue("raw_thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Inject.
	if e := a.g.NewThreadWithHex(bpk, tData); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

func (a *API) NewPostWithHex(w http.ResponseWriter, r *http.Request) {
	// Get board public key.
	bpk, e := misc.GetPubKey(r.FormValue("board"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get thread reference.
	tRef, e := misc.GetReference(r.FormValue("thread"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Get post data.
	pData, e := misc.GetBytes(r.FormValue("raw_post"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	// Inject.
	if e := a.g.NewPostWithHex(bpk, tRef, pData); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
}

/*
	<<< TESTS >>>
*/

func (a *API) TestNewFilledBoard(w http.ResponseWriter, r *http.Request) {
	seed := r.FormValue("seed")

	threads, e := strconv.Atoi(r.FormValue("threads"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}

	minPosts, e := strconv.Atoi(r.FormValue("min_posts"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}

	maxPosts, e := strconv.Atoi(r.FormValue("max_posts"))
	if e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}

	if e := a.g.TestNewFilledBoard(seed, threads, minPosts, maxPosts); e != nil {
		sendResponse(w, e.Error(), http.StatusBadRequest)
		return
	}
	sendResponse(w, true, http.StatusOK)
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
