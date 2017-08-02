package v1

import (
	"context"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/content"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
	"os"
	"sync"
	"github.com/skycoin/bbs/src/store/state/states"
)

type seqWaiter struct {
	seq  uint64
	done chan struct{}
}

// BoardState represents an internal state of a board.
type BoardState struct {
	l                  *log.Logger
	bpk, user          cipher.PubKey
	tMux, pMux         sync.Mutex
	t, p               map[skyobject.Reference]*object.VoteSummary
	seq                uint64        // Last processed sequence of root.
	tsh, tdh, psh, pdh cipher.SHA256 // Hashes (thread (store, deleted), post (store, deleted)).
	workers            chan<- func()
	seqWaiters         chan *seqWaiter
	newRoots           chan *node.Root
	quit               chan struct{}
	wg                 sync.WaitGroup
}

func NewBoardState(bpk, upk cipher.PubKey, workerChan chan<- func()) states.State {
	s := &BoardState{
		l:          inform.NewLogger(false, os.Stdout, "BOARDSTATE:"+bpk.Hex()),
		bpk:        bpk,
		user:       upk,
		t:          make(map[skyobject.Reference]*object.VoteSummary),
		p:          make(map[skyobject.Reference]*object.VoteSummary),
		workers:    workerChan,
		seqWaiters: make(chan *seqWaiter),
		newRoots:   make(chan *node.Root),
		quit:       make(chan struct{}),
	}
	go s.service()
	return s
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

func (s *BoardState) Trigger(ctx context.Context, root *node.Root) {
	select {
	case s.newRoots <- root:
	case <-ctx.Done():
		s.l.Println(ctx.Err())
	}
}

func (s *BoardState) service() {
	s.wg.Add(1)
	defer s.wg.Done()

	var seqWaiters []*seqWaiter

	for {
		select {
		case root := <-s.newRoots:
			for len(s.newRoots) > 1 {
				s.l.Printf("SKIPPING root of seq (%d).", root.Seq())
				root = <-s.newRoots
			}
			if root.Seq() <= s.seq {
				s.l.Printf("SKIPPING root of seq (%d).", root.Seq())
				break
			}
			s.l.Printf("PROCESSING root of seq (%d)", root.Seq())
			s.seq = root.Seq()
			s.processRoot(root)

			for i := len(seqWaiters) - 1; i >= 0; i++ {
				if seqWaiters[i].seq < root.Seq() {
					select {
					case seqWaiters[i].done <- struct{}{}:
					default:
					}
					seqWaiters[i], seqWaiters[0] =
						seqWaiters[0], seqWaiters[i]
					seqWaiters = seqWaiters[1:]
				}
			}

		case temp := <-s.seqWaiters:
			s.l.Printf("WAITING for seq > %d ...", temp.seq)
			if temp.seq <= s.seq {
				select {
				case temp.done <- struct{}{}:
				default:
				}
			} else {
				seqWaiters = append(seqWaiters, temp)
			}

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

	// Get pages hash.
	temp := struct{ tsh, tdh, psh, pdh cipher.SHA256 }{
		tsh: cipher.SumSHA256(encoder.Serialize(result.ThreadVotesPage.Store)),
		tdh: cipher.SumSHA256(encoder.Serialize(result.ThreadVotesPage.Deleted)),
		psh: cipher.SumSHA256(encoder.Serialize(result.PostVotesPage.Store)),
		pdh: cipher.SumSHA256(encoder.Serialize(result.PostVotesPage.Deleted)),
	}

	var pageWG sync.WaitGroup

	if temp.tsh != s.tsh {
		pageWG.Add(len(result.ThreadVotesPage.Store))
		for _, page := range result.ThreadVotesPage.Store {
			s.processVotesPage(root, &page, &pageWG, s.GetThreadVotes, s.setThreadVotes)
		}
		s.tsh = temp.tsh
	}

	if temp.tdh != s.tdh {
		for _, dRef := range result.ThreadVotesPage.Deleted {
			s.tMux.Lock()
			delete(s.t, skyobject.Reference(dRef))
			s.tMux.Unlock()
		}
		s.tdh = temp.tdh
	}

	if temp.psh != s.psh {
		pageWG.Add(len(result.PostVotesPage.Store))
		for _, page := range result.PostVotesPage.Store {
			s.processVotesPage(root, &page, &pageWG, s.GetPostVotes, s.setPostVotes)
		}
		s.psh = temp.psh
	}

	if temp.pdh != s.pdh {
		for _, dRef := range result.PostVotesPage.Deleted {
			s.pMux.Lock()
			delete(s.p, skyobject.Reference(dRef))
			s.pMux.Unlock()
		}
		s.pdh = temp.pdh
	}

	pageWG.Wait()
}

func (s *BoardState) processVotesPage(
	root *node.Root,
	page *object.VotesPage,
	pageWG *sync.WaitGroup,
	get func(skyobject.Reference) *object.VoteSummary,
	set func(skyobject.Reference, *object.VoteSummary),
) {
	summary := &object.VoteSummary{
		VotesHash: cipher.SumSHA256(encoder.Serialize(page)),
	}

	if summary.VotesHash != get(skyobject.Reference(page.Ref)).VotesHash {
		var summaryWG sync.WaitGroup
		summaryWG.Add(len(page.Votes))

		for _, vRef := range page.Votes {
			s.processVote(root, vRef, summary, &summaryWG)
		}

		saveSummary := func() {
			defer pageWG.Done()
			summaryWG.Wait()
			set(skyobject.Reference(page.Ref), summary)
		}
		s.workers <- saveSummary
	} else {
		pageWG.Done()
	}
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
		if e := i.Run(); e != nil {
			s.l.Printf("Error: %s", e.Error())
		}
	}
	s.workers <- updateSummary
}

func (s *BoardState) GetThreadVotes(tRef skyobject.Reference) *object.VoteSummary {
	s.l.Printf("GetThreadVotes : thread '%s'.", tRef.String())
	summary, has := s.getThreadVotes(tRef)
	if !has {
		s.l.Println("/t- (NOT HAS)")
		summary = new(object.VoteSummary)
	}
	return summary
}

func (s *BoardState) GetThreadVotesSeq(
	ctx context.Context, tRef skyobject.Reference, seq uint64) *object.VoteSummary {
	seqW := &seqWaiter{
		seq:  seq,
		done: make(chan struct{}),
	}
	go func(ctx context.Context) {
		select {
		case s.seqWaiters <- seqW:
		case <-ctx.Done():
		}
	}(ctx)
	<-seqW.done
	return s.GetThreadVotes(tRef)
}

func (s *BoardState) GetPostVotes(pRef skyobject.Reference) *object.VoteSummary {
	summary, has := s.getPostVotes(pRef)
	if !has {
		summary = new(object.VoteSummary)
	}
	return summary
}

func (s *BoardState) GetPostVotesSeq(
	ctx context.Context, pRef skyobject.Reference, seq uint64) *object.VoteSummary {
	seqW := &seqWaiter{
		seq:  seq,
		done: make(chan struct{}),
	}
	go func(ctx context.Context) {
		select {
		case s.seqWaiters <- seqW:
		case <-ctx.Done():
		}
	}(ctx)
	<-seqW.done
	return s.GetPostVotes(pRef)
}

func (s *BoardState) getThreadVotes(
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
	s.l.Printf("setThreadVotes : thread '%s' data '%v'", tRef.String(), data)
	s.t[tRef] = data
}

func (s *BoardState) getPostVotes(
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
	s.p[pRef] = data
}
