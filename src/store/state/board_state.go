package state

import (
	"fmt"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/content"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"os"
	"sync"
)

// BoardState represents an internal state of a board.
type BoardState struct {
	l          *log.Logger
	bpk, user  cipher.PubKey
	tMux, pMux sync.Mutex
	t, p       map[skyobject.Reference]*object.VoteSummary
	seq        uint64 // Last processed sequence of root.
	workers    chan<- func()
	newRoots   chan *node.Root
	quit       chan struct{}
	wg         sync.WaitGroup
}

func NewBoardState(bpk, user cipher.PubKey, workers chan<- func()) *BoardState {
	bs := &BoardState{
		l:        inform.NewLogger(true, os.Stdout, "BOARDSTATE:"+bpk.Hex()),
		bpk:      bpk,
		user:     user,
		t:        make(map[skyobject.Reference]*object.VoteSummary),
		p:        make(map[skyobject.Reference]*object.VoteSummary),
		workers:  workers,
		newRoots: make(chan *node.Root),
		quit:     make(chan struct{}),
	}
	go bs.service()
	return bs
}

func (s *BoardState) Close() {
	s.tMux.Lock()
	s.pMux.Lock()
	s.l.Println("Closing...")
	defer s.tMux.Unlock()
	defer s.pMux.Unlock()
	defer s.l.Println("Closed.")
	for {
		select {
		case s.quit <- struct{}{}:
		default:
			s.wg.Wait()
			return
		}
	}
}

func (s *BoardState) service() {
	s.wg.Add(1)
	defer s.wg.Done()
	for {
		select {
		case root := <-s.newRoots:
			for len(s.newRoots) > 1 || root.Seq() <= s.seq {
				s.l.Printf("SKIPPING root of seq (%d).", root.Seq())
				root = <-s.newRoots
			}
			s.l.Printf("PROCESSING root of seq (%d)", root.Seq())
			s.seq = root.Seq()
			s.processRoot(root)

		case <-s.quit:
			return
		}
	}
}

func (s *BoardState) processRoot(root *node.Root) {
	result := content.NewResult(root).
		GetPages(false, true, true)

	if e := result.Error(); e != nil {
		s.l.Printf("PROCESSING result error: %s", e.Error())
		return
	}

	var pageWG sync.WaitGroup
	pageWG.Add(
		len(result.ThreadVotesPage.Store) +
		len(result.PostVotesPage.Store),
	)

	for _, page := range result.ThreadVotesPage.Store {
		s.processVotesPage(root, &page, &pageWG, s.setThreadVotes)
	}

	for _, page := range result.PostVotesPage.Store {
		s.processVotesPage(root, &page, &pageWG, s.setPostVotes)
	}

	pageWG.Wait()
}

func (s *BoardState) processVotesPage(
	root *node.Root,
	page *object.VotesPage,
	pageWG *sync.WaitGroup,
	setter func(skyobject.Reference, *object.VoteSummary),
) {
	summary := new(object.VoteSummary)
	var summaryWG sync.WaitGroup
	summaryWG.Add(len(page.Votes))

	for _, vRef := range page.Votes {
		s.processVote(root, vRef, summary, &summaryWG)
	}

	saveSummary := func() {
		defer pageWG.Done()
		summaryWG.Wait()
		setter(skyobject.Reference(page.Ref), summary)
	}
	s.workers <- saveSummary
}

func (s *BoardState) processVote(
	root *node.Root,
	vRef skyobject.Reference,
	summary *object.VoteSummary,
	summaryWG *sync.WaitGroup,
) {
	data, has := root.Get(vRef)
	if !has {
		summaryWG.Done()
		return
	}
	i := &Instruction{
		user:    &s.user,
		data:    data,
		summary: summary,
	}
	updateSummary := func() {
		defer summaryWG.Done()
		e := i.Run()
		fmt.Println(e)
	}
	s.workers <- updateSummary
}

func (s *BoardState) votePageWorker() {
	for {
		select {
		case <-s.quit:
			return
		}
	}
}

func (s *BoardState) GetThreadVotes(
	tRef skyobject.Reference) (*object.VoteSummary, bool,
) {
	s.tMux.Lock()
	defer s.tMux.Unlock()
	out, has := s.t[tRef]
	return out, has
}

func (s *BoardState) setThreadVotes(
	tRef skyobject.Reference, data *object.VoteSummary,
) {
	s.tMux.Lock()
	defer s.tMux.Unlock()
	//s.l.Printf("SETTING votes for thread '%s' : %v", tRef.String(), data)
	s.t[tRef] = data
}

func (s *BoardState) GetPostVotes(
	pRef skyobject.Reference) (*object.VoteSummary, bool,
) {
	s.pMux.Lock()
	defer s.pMux.Unlock()
	out, has := s.p[pRef]
	return out, has
}

func (s *BoardState) setPostVotes(
	pRef skyobject.Reference, data *object.VoteSummary,
) {
	s.pMux.Lock()
	defer s.pMux.Unlock()
	//s.l.Printf("SETTING votes for post '%s' : %v", pRef.String(), data)
	s.p[pRef] = data
}
