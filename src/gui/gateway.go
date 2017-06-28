package gui

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/cmd/bbsnode/args"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/msg"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"time"
)

// Gateway represents the intermediate between External calls and internal processing.
// It can be seen as a security layer.
type Gateway struct {
	config     *args.Config
	container  *cxo.Container
	boardSaver *store.BoardSaver
	userSaver  *store.UserSaver
	queueSaver *msg.QueueSaver
	quit       chan int
}

// NewGateway creates a new Gateway.
func NewGateway(
	config *args.Config,
	container *cxo.Container,
	boardSaver *store.BoardSaver,
	userSaver *store.UserSaver,
	queueSaver *msg.QueueSaver,
	quit chan int,
) *Gateway {
	return &Gateway{
		config:     config,
		container:  container,
		boardSaver: boardSaver,
		userSaver:  userSaver,
		queueSaver: queueSaver,
		quit:       quit,
	}
}

/*
	<<< MISC >>>
*/

// Quit quits the node entirely.
func (g *Gateway) Quit() bool {
	timer := time.NewTimer(10 * time.Second)
	select {
	case g.quit <- 0:
		return true
	case <-timer.C:
		return false
	}
}

// PingSubmissionAddress pings a submission address.
func (g *Gateway) PingSubmissionAddress(address string) error {
	return g.queueSaver.Ping(address)
}

/*
	<<< FOR STATS >>>
*/

type StatsView struct {
	NodeIsMaster   bool   `json:"node_is_master"`
	NodeCXOAddress string `json:"node_cxo_address"`
}

func (g *Gateway) GetStats() *StatsView {
	return &StatsView{
		NodeIsMaster:   g.config.Master(),
		NodeCXOAddress: g.container.GetAddress(),
	}
}

/*
	<<< FOR CONNECTIONS >>>
*/

// GetConnection lists all connections.
func (g *Gateway) GetConnections() []string {
	return g.container.GetConnections()
}

// AddConnection adds a connection.
func (g *Gateway) AddConnection(addr string) error {
	return g.container.Connect(addr)
}

// RemoveConnection removes a connection.
func (g *Gateway) RemoveConnection(addr string) error {
	boards := g.boardSaver.GetOfAddress(addr)
	if len(boards) == 0 {
		return g.container.Disconnect(addr)
	}
	return errors.Errorf("currently subscribed to %d boards under address %s",
		len(boards), addr)
}

/*
	<<< FOR SUBSCRIPTIONS >>>
*/

// GetSubscriptions lists all subscriptions.
func (g *Gateway) GetSubscriptions() []store.BoardInfo {
	return g.boardSaver.List()
}

// GetSubscription gets a subscription.
func (g *Gateway) GetSubscription(bpk cipher.PubKey) (store.BoardInfo, bool) {
	return g.boardSaver.Get(bpk)
}

// Subscribe subscribes to a board.
func (g *Gateway) Subscribe(addr string, bpk cipher.PubKey) {
	g.boardSaver.Add(addr, bpk)
}

// Unsubscribe unsubscribes from a board.
func (g *Gateway) Unsubscribe(bpk cipher.PubKey) {
	g.boardSaver.Remove(bpk)
}

/*
	<<< FOR USERS >>>
*/

// GetCurrentUser returns the currently active user.
func (g *Gateway) GetCurrentUser() store.UserConfig {
	return g.userSaver.GetCurrent()
}

// SetCurrentUser sets the currently active user.
func (g *Gateway) SetCurrentUser(upk cipher.PubKey) error {
	return g.userSaver.SetCurrent(upk)
}

// GetMasterUsers lists all users this node is master of.
func (g *Gateway) GetMasterUsers() []store.UserConfig {
	return g.userSaver.ListMasters()
}

// NewMasterUser creates a new user configuration of a master user.
// It will replace if user already exists.
func (g *Gateway) NewMasterUser(alias, seed string) store.UserConfig {
	pk, sk := cipher.GenerateDeterministicKeyPair([]byte(seed))
	g.userSaver.MasterAdd(alias, pk, sk)
	uc, _ := g.userSaver.Get(pk)
	return uc
}

// GetUsers lists all users, master or not.
func (g *Gateway) GetUsers() []store.UserConfig {
	return g.userSaver.List()
}

// NewUser creates a new user configuration for a user we are not master of.
// It will replace if user already exists.
func (g *Gateway) NewUser(alias string, upk cipher.PubKey) store.UserConfig {
	g.userSaver.Add(alias, upk)
	uc, _ := g.userSaver.Get(upk)
	return uc
}

// RemoveUser removes a user configuration, master or not.
func (g *Gateway) RemoveUser(upk cipher.PubKey) error {
	return g.userSaver.Remove(upk)
}

/*
	<<< FOR BOARD META >>>
*/

// GetBoardMeta obtains the board meta.
func (g *Gateway) GetBoardMeta(bpk cipher.PubKey) (*typ.BoardMeta, error) {
	board, e := g.container.GetBoard(bpk)
	if e != nil {
		return nil, errors.Wrap(e, "unable to obtain board")
	}
	meta, e := board.GetMeta()
	return meta, errors.Wrap(e, "unable to obtain board meta")
}

// GetSubmissionAddresses gets submission addresses from board meta.
func (g *Gateway) GetSubmissionAddresses(bpk cipher.PubKey) ([]string, error) {
	list, e := g.container.GetSubmissionAddresses(bpk)
	return list, errors.Wrap(e, "failed to obtain submission addresses of board")
}

// AddSubmissionAddress adds a submission address to board meta.
func (g *Gateway) AddSubmissionAddress(bpk cipher.PubKey, address string) error {
	bi, has := g.boardSaver.Get(bpk)
	if !has {
		return errors.New("board not found")
	}
	if !bi.Config.Master {
		return errors.New("not master of board")
	}
	return g.container.AddSubmissionAddress(bpk, bi.Config.GetSK(), address)
}

// RemoveSubmissionAddress removes a submission address from board meta.
func (g *Gateway) RemoveSubmissionAddress(bpk cipher.PubKey, address string) error {
	bi, has := g.boardSaver.Get(bpk)
	if !has {
		return errors.New("board not found")
	}
	if !bi.Config.Master {
		return errors.New("not master of board")
	}
	return g.container.RemoveSubmissionAddress(bpk, bi.Config.GetSK(), address)
}

/*
	<<< FOR BOARDS, THREADS & POSTS >>>
*/

// GetBoards lists all boards.
func (g *Gateway) GetBoards() []*typ.Board {
	return g.container.GetBoards(g.boardSaver.ListKeys()...)
}

// GetBoard obtains a single board.
func (g *Gateway) GetBoard(bpk cipher.PubKey) (*typ.Board, error) {
	board, e := g.container.GetBoard(bpk)
	return board, errors.Wrap(e, "unable to obtain board")
}

// NewBoard creates a new board.
func (g *Gateway) NewBoard(board *typ.Board, seed string) (bi store.BoardInfo, e error) {
	bm, e := board.GetMeta()
	if e != nil {
		e = errors.Wrap(e, "failed to obtain board meta")
		return
	}
	if len(bm.SubmissionAddresses) == 0 {
		bm.AddSubmissionAddress(g.config.RPCRemAdr())
		if e = board.SetMeta(bm); e != nil {
			e = errors.Wrap(e, "failed to re-set board meta")
			return
		}
	}
	bpk, bsk := board.TouchWithSeed([]byte(seed))
	if _, e = board.Check(); e != nil {
		e = errors.Wrap(e, "failed to create board")
		return
	}
	if e = g.boardSaver.MasterAdd(bpk, bsk); e != nil {
		return
	}
	//board.URL = g.config.RPCRemAdr()
	if e = g.container.NewBoard(board, bpk, bsk); e != nil {
		return
	}
	bi, _ = g.boardSaver.Get(bpk)
	return
}

// RemoveBoard removes a board
func (g *Gateway) RemoveBoard(bpk cipher.PubKey) error {
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to the board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == true {
		// Via Container.
		log.Println("[GUI GW] Master, remove board!")
		if e := g.container.RemoveBoard(bpk, bi.Config.GetSK()); e != nil {
			return e
		}
	} else {
		// threads and posts are only to be deleted from master.
		return errors.New("not owning the board")
	}
	return nil
}

type BoardPageView struct {
	Board   *typ.Board    `json:"board"`
	Threads []*typ.Thread `json:"threads"`
}

func (g *Gateway) GetBoardPage(bpk cipher.PubKey) (*BoardPageView, error) {
	board, e := g.container.GetBoard(bpk)
	if e != nil {
		return nil, errors.Wrap(e, "unable to obtain board")
	}
	threads, e := g.container.GetThreads(bpk)
	if e != nil {
		return nil, errors.Wrap(e, "unable to obtain threads")
	}
	return &BoardPageView{board, threads}, nil
}

// GetThreads obtains threads of boards we are subscribed to.
// Input `bpks` acts as a filter.
// If no `bpks` are specified, threads of all boards will be obtained.
// If one or more `bpks` are specified, only threads under those boards will be returned.
func (g *Gateway) GetThreads(bpks ...cipher.PubKey) []*typ.Thread {
	tMap := make(map[string]*typ.Thread)
	switch len(bpks) {
	case 0:
		for _, bpk := range g.boardSaver.ListKeys() {
			ts, e := g.container.GetThreads(bpk)
			if e != nil {
				continue
			}
			for _, t := range ts {
				tMap[t.Ref] = t
			}
		}
	default:
		for _, bpk := range bpks {
			if _, has := g.boardSaver.Get(bpk); has == false {
				return nil
			}
			ts, e := g.container.GetThreads(bpk)
			if e != nil {
				continue
			}
			for _, t := range ts {
				tMap[t.Ref] = t
			}
		}
	}
	threads, i := make([]*typ.Thread, len(tMap)), 0
	for _, t := range tMap {
		threads[i] = t
		i += 1
	}
	return threads
}

// NewThread creates a new thread and sets the board of public key as it's master.
func (g *Gateway) NewThread(bpk cipher.PubKey, thread *typ.Thread) error {
	// Check thread.
	if e := thread.Check(); e != nil {
		return e
	}
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to the board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == true {
		// Via Container.
		return g.container.NewThread(bpk, bi.Config.GetSK(), thread)
	} else {
		// Via RPC Client.
		uc := g.userSaver.GetCurrent()
		return g.queueSaver.AddNewThreadReq(bpk, uc.GetPK(), uc.GetSK(), thread)
	}
}

// RemoveThread removes a thread
func (g *Gateway) RemoveThread(bpk cipher.PubKey, tRef skyobject.Reference) error {
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == true {
		// Via Container.

		// Obtain thread.
		thread, e := g.container.GetThread(tRef)
		if e != nil {
			return e
		}
		// Obtain thread master's public key.
		masterPK, e := misc.GetPubKey(thread.MasterBoard)
		if e != nil {
			return e
		}
		// Remove dependency (if has).
		bi.Config.RemoveDep(masterPK, tRef)
		// Remove thread.
		if e := g.container.RemoveThread(bpk, bi.Config.GetSK(), tRef); e != nil {
			return e
		}
	} else {
		// threads and posts are only to be deleted from master.
		return errors.New("not owning the board")
	}
	return nil
}

type ThreadPageView struct {
	Board  *typ.Board  `json:"board"`
	Thread *typ.Thread `json:"thread"`
	Posts  []*typ.Post `json:"posts"`
}

func (g *Gateway) GetThreadPage(bpk cipher.PubKey, tRef skyobject.Reference) (*ThreadPageView, error) {
	b, e := g.container.GetBoard(bpk)
	if e != nil {
		return nil, errors.Wrap(e, "unable to obtain board")
	}
	thread, posts, e := g.container.GetThreadPage(bpk, tRef)
	if e != nil {
		return nil, errors.Wrap(e, "unable to obtain threadpage")
	}
	return &ThreadPageView{b, thread, posts}, nil
}

// GetPosts obtains posts of specified board and thread.
func (g *Gateway) GetPosts(bpk cipher.PubKey, tRef skyobject.Reference) ([]*typ.Post, error) {
	_, has := g.boardSaver.Get(bpk)
	if has == false {
		return nil, errors.New("not subscribed to board")
	}
	return g.container.GetPosts(bpk, tRef)
}

// NewPost creates a new post in specified board and thread.
func (g *Gateway) NewPost(bpk cipher.PubKey, tRef skyobject.Reference, post *typ.Post) (e error) {
	// Check post.
	uc := g.userSaver.GetCurrent()
	if e = post.Sign(uc.GetPK(), uc.GetSK()); e != nil {
		return
	}
	post.Touch()
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == true {
		// Via Container.
		return g.container.NewPost(bpk, bi.Config.GetSK(), tRef, post)
	} else {
		// Via RPC Client.
		return g.queueSaver.AddNewPostReq(bpk, tRef, post)
	}
}

// RemovePost removes a post in specified board and thread.
func (g *Gateway) RemovePost(bpk cipher.PubKey, tRef, pRef skyobject.Reference) (e error) {
	// Check board.
	bi, has := g.boardSaver.Get(bpk)
	if has == false {
		return errors.New("not subscribed to board")
	}
	// Check if this BBS Node owns the board.
	if bi.Config.Master == true {
		log.Println("[GUI GW] Master, remove the post!")
		if e = g.container.RemovePost(bpk, bi.Config.GetSK(), tRef, pRef); e != nil {
			fmt.Println(e)
			return e
		}
	} else {
		// threads and posts are only to be deleted from master.
		return errors.New("not owning the board")
	}
	return nil
}

// ImportThread imports a thread from one board to a board which this node is master of.
func (g *Gateway) ImportThread(fromBpk, toBpk cipher.PubKey, tRef skyobject.Reference) error {
	// Check "to" board.
	bi, has := g.boardSaver.Get(toBpk)
	if !has {
		return errors.New("not subscribed to board")
	}
	if bi.Config.Master == false {
		return errors.New("to board is not master")
	}
	// Add Dependency to BoardSaver.
	if e := g.boardSaver.AddBoardDep(toBpk, fromBpk, tRef); e != nil {
		return e
	}
	// Import thread.
	return g.container.ImportThread(fromBpk, toBpk, bi.Config.GetSK(), tRef)
}

/*
	<<< VOTES >>>
*/

type VotesView struct {
	UpVotes             int  `json:"up_votes"`
	DownVotes           int  `json:"down_votes"`
	CurrentUserVoted    bool `json:"current_user_voted"`
	CurrentUserVoteMode int  `json:"current_user_vote_mode,omitempty"`
}

func (g *Gateway) GetVotesForThread(bpk cipher.PubKey, tRef skyobject.Reference) (*VotesView, error) {
	// Get current user.
	cu := g.userSaver.GetCurrent()
	upk := cu.GetPK()
	// Get votes.
	votes, e := g.container.GetVotesForThread(bpk, tRef)
	if e != nil {
		return nil, e
	}
	vv := &VotesView{}
	for _, vote := range votes {
		switch vote.Mode {
		case +1:
			vv.UpVotes += 1
		case -1:
			vv.DownVotes += 1
		}
		if vote.User == upk {
			vv.CurrentUserVoted = true
			vv.CurrentUserVoteMode = int(vote.Mode)
		}
	}
	return vv, nil
}

func (g *Gateway) GetVotesForPost(bpk cipher.PubKey, pRef skyobject.Reference) (*VotesView, error) {
	// Get current user.
	cu := g.userSaver.GetCurrent()
	upk := cu.GetPK()
	// Get votes.
	votes, e := g.container.GetVotesForPost(bpk, pRef)
	if e != nil {
		return nil, e
	}
	vv := &VotesView{}
	for _, vote := range votes {
		switch vote.Mode {
		case +1:
			vv.UpVotes += 1
		case -1:
			vv.DownVotes += 1
		}
		if vote.User == upk {
			vv.CurrentUserVoted = true
			vv.CurrentUserVoteMode = int(vote.Mode)
		}
	}
	return vv, nil
}

func (g *Gateway) VoteForThread(bpk cipher.PubKey, tRef skyobject.Reference, vote *typ.Vote) error {
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
		// Via Container.
		switch vote.Mode {
		case 0:
			return g.container.RemoveVoteForThread(uc.GetPK(), bpk, bi.Config.GetSK(), tRef)
		case -1, +1:
			return g.container.AddVoteForThread(bpk, bi.Config.GetSK(), tRef, vote)
		}
	} else {
		return g.queueSaver.AddVoteThreadReq(bpk, tRef, vote)
	}
	return nil
}

func (g *Gateway) VoteForPost(bpk cipher.PubKey, pRef skyobject.Reference, vote *typ.Vote) error {
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
		// Via Container.
		switch vote.Mode {
		case 0:
			return g.container.RemoveVoteForPost(uc.GetPK(), bpk, bi.Config.GetSK(), pRef)
		case -1, +1:
			return g.container.AddVoteForPost(bpk, bi.Config.GetSK(), pRef, vote)
		}
	} else {
		return g.queueSaver.AddVotePostReq(bpk, pRef, vote)
	}
	return nil
}

/*
	<<< HEX >>>
*/

func (g *Gateway) GetThreadPageAsHex(bpk cipher.PubKey, tRef skyobject.Reference) (*cxo.ThreadPageHex, error) {
	return g.container.GetThreadPageAsHex(bpk, tRef)
}

func (g *Gateway) GetThreadPageWithTpRefAsHex(tpRef skyobject.Reference) (*cxo.ThreadPageHex, error) {
	return g.container.GetThreadPageWithTpRefAsHex(tpRef)
}

func (g *Gateway) NewThreadWithHex(bpk cipher.PubKey, tData []byte) error {
	bi, has := g.boardSaver.Get(bpk)
	if !has {
		return errors.New("not subscribed to board")
	}
	if !bi.Config.Master {
		return errors.New("not master of board")
	}
	return g.container.NewThreadWithHex(bpk, bi.Config.GetSK(), tData)
}

func (g *Gateway) NewPostWithHex(bpk cipher.PubKey, tRef skyobject.Reference, pData []byte) error {
	bi, has := g.boardSaver.Get(bpk)
	if !has {
		return errors.New("not subscribed to board")
	}
	if !bi.Config.Master {
		return errors.New("not master of board")
	}
	return g.container.NewPostWithHex(bpk, bi.Config.GetSK(), tRef, pData)
}

/*
	<<< TESTS >>>
*/

// TestNewFilledBoard creates a new board with given seed, filled with threads and posts.
func (g *Gateway) TestNewFilledBoard(seed string, threads, minPosts, maxPosts int) error {
	if threads < 0 || minPosts < 0 || maxPosts < 0 || maxPosts-minPosts < 0 {
		return errors.New("invalid inputs")
	}
	b := &typ.Board{
		Name: fmt.Sprintf("Test Board '%s'", seed),
		Desc: fmt.Sprintf("A board with '%s' as seed and %d threads.", seed, threads),
	}
	b.SetMeta(new(typ.BoardMeta))
	bi, e := g.NewBoard(b, seed)
	if e != nil {
		return e
	}
	bpk := bi.Config.GetPK()
	for i := 1; i <= threads; i++ {
		t := &typ.Thread{
			Name: fmt.Sprintf("Thread %d", i),
			Desc: fmt.Sprintf("A test thread on board with seed '%s'.", seed),
		}
		if e := g.NewThread(bpk, t); e != nil {
			return errors.New("on creating thread " + string(i) + "; " + e.Error())
		}
		nPosts, e := misc.MakeIntBetween(minPosts, maxPosts)
		if e != nil {
			return e
		}
		for j := 1; j <= nPosts; j++ {

			p := &typ.Post{
				Title: fmt.Sprintf("Post %d", j),
				Body:  fmt.Sprintf("This is post %d on thread %d.", j, i),
			}
			if e := g.NewPost(bpk, t.GetRef(), p); e != nil {
				return e
			}
		}
	}
	return nil
}
