package dev

import (
	"fmt"
	"github.com/evanlinjin/bbs/cmd/bbsnode/args"
	"github.com/evanlinjin/bbs/extern/gui"
	"github.com/evanlinjin/bbs/intern/store"
	"github.com/evanlinjin/bbs/intern/typ"
	"github.com/evanlinjin/bbs/misc"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"time"
)

// Tester represents a tester.
// It autonomously creates threads and posts.
type Tester struct {
	config *args.Config
	g      *gui.Gateway
	bpk    cipher.PubKey
	tRefs  skyobject.References
	users  []store.UserConfig
	pCount int
	quit   chan struct{}
}

// NewTester creates a new tester.
func NewTester(config *args.Config, gateway *gui.Gateway) (*Tester, error) {
	t := &Tester{
		config: config,
		g:      gateway,
		quit:   make(chan struct{}),
	}
	if e := t.setupUsers(); e != nil {
		return nil, e
	}
	if e := t.setupBoard(); e != nil {
		return nil, e
	}
	go t.service()
	return t, nil
}

func (t *Tester) setupUsers() error {
	nGoal := t.config.TestModeUsers()
	nNow := len(t.g.GetMasterUsers())
	for i := 0; i < nGoal-nNow; i++ {
		t.g.NewMasterUser(misc.MakeRandomAlias(), misc.MakeTimeStampedRandomID(100).Hex())
	}
	t.users = t.g.GetMasterUsers()
	return nil
}

func (t *Tester) setupBoard() error {
	seed := misc.MakeTimeStampedRandomID(100).Hex()
	pk, _ := cipher.GenerateDeterministicKeyPair([]byte(seed))
	if e := t.g.TestNewFilledBoard(seed, t.config.TestModeThreads(), 1, 1); e != nil {
		return e
	}
	t.bpk = pk
	threads := t.g.GetThreads(t.bpk)
	t.tRefs = make(skyobject.References, len(threads))
	for i, thread := range threads {
		t.tRefs[i] = thread.GetRef()
	}
	return nil
}

func (t *Tester) service() {
	for {
		choice, _ := misc.MakeIntBetween(0, 3)
		switch choice {
		case 0:
			t.actionChangeUser()
		case 1:
			t.actionNewPost()
		case 2:
			t.actionVotePost()
		case 3:
			t.actionVoteThread()
		}
		select {
		case <-t.quit:
			return
		default:
			time.Sleep(t.getInterval())
			continue
		}
	}
}

func (t *Tester) Close() {
	go func() {
		t.quit <- struct{}{}
	}()
}

func (t *Tester) getInterval() time.Duration {
	i, e := misc.MakeIntBetween(
		t.config.TestModeMinInterval(),
		t.config.TestModeMaxInterval(),
	)
	if e != nil {
		log.Panic(e)
	}
	return time.Duration(i) * time.Second
}

func (t *Tester) getRandomThreadRef() skyobject.Reference {
	log.Printf("[TESTER] getRandomThreadRef: 0, %d", len(t.tRefs)-1)
	i, e := misc.MakeIntBetween(0, len(t.tRefs)-1)
	if e != nil {
		log.Panic(e)
	}
	return t.tRefs[i]
}

func (t *Tester) getRandomPostRef(tRef skyobject.Reference) skyobject.Reference {
	posts, e := t.g.GetPosts(t.bpk, tRef)
	if e != nil {
		log.Panic(e)
	}
	i, e := misc.MakeIntBetween(0, len(posts)-1)
	if e != nil {
		log.Panic(e)
	}
	ref, e := misc.GetReference(posts[i].Ref)
	if e != nil {
		log.Panic(e)
	}
	return ref
}

func (t *Tester) actionChangeUser() {
	i, e := misc.MakeIntBetween(0, len(t.users)-1)
	if e != nil {
		log.Panic(e)
	}
	if e := t.g.SetCurrentUser(t.users[i].GetPK()); e != nil {
		log.Panic(e)
	}
}

func (t *Tester) actionNewPost() {
	user := t.g.GetCurrentUser()
	tRef := t.getRandomThreadRef()
	post := &typ.Post{
		Title: fmt.Sprintf("Test Post %d", t.pCount),
		Body:  fmt.Sprintf("This is a test post by test user %s.", user.Alias),
	}
	if e := post.Sign(user.GetPK(), user.GetSK()); e != nil {
		log.Panic(e)
	}
	if e := t.g.NewPost(t.bpk, tRef, post); e != nil {
		log.Panic(e)
	}
}

func (t *Tester) actionVotePost() {
	user := t.g.GetCurrentUser()
	pRef := t.getRandomPostRef(t.getRandomThreadRef())
	vMode, e := misc.MakeIntBetween(-1, +1)
	if e != nil {
		log.Panic(e)
	}
	vote := &typ.Vote{Mode: int8(vMode)}
	if e := vote.Sign(user.GetPK(), user.GetSK()); e != nil {
		log.Panic(e)
	}
	if e := t.g.VoteForPost(t.bpk, pRef, vote); e != nil {
		log.Panic(e)
	}
}

func (t *Tester) actionVoteThread() {
	user := t.g.GetCurrentUser()
	tRef := t.getRandomThreadRef()
	vMode, e := misc.MakeIntBetween(-1, +1)
	if e != nil {
		log.Panic(e)
	}
	vote := &typ.Vote{Mode: int8(vMode)}
	if e := vote.Sign(user.GetPK(), user.GetSK()); e != nil {
		log.Panic(e)
	}
	if e := t.g.VoteForThread(t.bpk, tRef, vote); e != nil {
		log.Panic(e)
	}
}
