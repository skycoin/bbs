package typ

import "github.com/skycoin/cxo/skyobject"

// ThreadVotePage is an element of ThreadVoteContainer.
type ThreadVotePage struct {
	Thread skyobject.Reference  `skyobject:"schema=Thread"`
	Votes  skyobject.References `skyobject:"schema=Vote"`
}

// ThreadVoteContainer contains the votes of threads.
type ThreadVoteContainer struct {
	Threads []ThreadVotePage
}

// GetThread obtains a ThreadVotePage from thread reference.
func (tvc *ThreadVoteContainer) GetThreadIndex(tRef skyobject.Reference) (int, bool) {
	for i, tvp := range tvc.Threads {
		if tvp.Thread == tRef {
			return i, true
		}
	}
	return -1, false
}

// AddThread adds a thread to ThreadVoteContainer.
func (tvc *ThreadVoteContainer) AddThread(tRef skyobject.Reference) {
	for _, t := range tvc.Threads {
		if t.Thread == tRef {
			return
		}
	}
	tvc.Threads = append(tvc.Threads, ThreadVotePage{Thread: tRef})
}

// RemoveThread removes a thread from ThreadVoteContainer.
func (tvc *ThreadVoteContainer) RemoveThread(tRef skyobject.Reference) {
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
