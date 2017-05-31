package typ

import (
	"github.com/pkg/errors"
	"github.com/skycoin/cxo/skyobject"
)

// ThreadVotes is an element of ThreadVotesContainer.
type ThreadVotes struct {
	Thread skyobject.Reference  `skyobject:"schema=Thread"`
	Votes  skyobject.References `skyobject:"schema=Vote"`
}

// ThreadVotesContainer contains the votes of threads.
type ThreadVotesContainer struct {
	Threads []ThreadVotes
}

// GetThreadVotes obtains the thread vote references for specified thread.
func (tvc *ThreadVotesContainer) GetThreadVotes(tRef skyobject.Reference) (*ThreadVotes, error) {
	for i := range tvc.Threads {
		if tvc.Threads[i].Thread == tRef {
			return &tvc.Threads[i], nil
		}
	}
	return nil, errors.New("thread votes not found")
}

// AddThread adds a thread to ThreadVotesContainer.
func (tvc *ThreadVotesContainer) AddThread(tRef skyobject.Reference) {
	for _, t := range tvc.Threads {
		if t.Thread == tRef {
			return
		}
	}
	tvc.Threads = append(tvc.Threads, ThreadVotes{Thread: tRef})
}

// RemoveThread removes a thread from ThreadVotesContainer.
func (tvc *ThreadVotesContainer) RemoveThread(tRef skyobject.Reference) {
	for i, t := range tvc.Threads {
		if t.Thread != tRef {
			continue
		}
		// Swap i'th and last element.
		tvc.Threads[i], tvc.Threads[len(tvc.Threads)-1] =
			tvc.Threads[len(tvc.Threads)-1], tvc.Threads[i]
		// Remove last element.
		tvc.Threads = tvc.Threads[:len(tvc.Threads)-1]
		return
	}
}
