package v1

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/misc/obtain"
	"github.com/skycoin/bbs/src/store/io"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state/states"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"os"
	"sync"
)

const (
	indexContent     = 0
	indexDeleted     = 1
	indexThreadVotes = 2
	indexPostVotes   = 3
	indexUserVotes   = 4
	countRootRefs    = 5
)

type pubPass struct {
	run  states.PublishFunc
	done chan error
}

type State struct {
	c          *states.StateConfig
	l          *log.Logger
	flag       skyobject.Flag
	node       *node.Node
	currentMux sync.Mutex
	current    *PackInstance
	inChan     chan *io.State   // Obtains new sequence.
	outChan    chan *io.Changes // Outputs updates.
	pubChan    chan *pubPass    // Need to generate.
	quitChan   chan struct{}    // Need to generate.
	wg         sync.WaitGroup
}

func (s *State) Init(config *states.StateConfig, node *node.Node) error {
	s.c = config
	s.l = inform.NewLogger(true, os.Stdout, "STATE:"+s.c.PubKey.Hex())

	// Prepare flags for unpacking root.
	s.flag = skyobject.HashTableIndex | skyobject.EntireTree
	if !s.c.Master {
		s.flag |= skyobject.ViewOnly
	}
	s.node = node

	root, e := node.Container().LastFull(config.PubKey)
	if e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to obtain root")
	}
	pack, e := s.unpackRoot(node, root)
	if e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to unpack root")
	}
	s.current, e = NewPackInstance(nil, pack)
	if e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to initiate pack instance")
	}

	s.inChan = make(chan *io.State, 10)
	s.outChan = make(chan *io.Changes)
	s.pubChan = make(chan *pubPass)
	s.quitChan = make(chan struct{})

	go s.service()
	return nil
}

func (s *State) unpackRoot(node *node.Node, root *skyobject.Root) (*skyobject.Pack, error) {
	return node.Container().Unpack(
		root,
		s.flag,
		node.Container().CoreRegistry().Types(),
		s.c.SecKey,
	)
}

func (s *State) Close() {
	for {
		select {
		case s.quitChan <- struct{}{}:
		default:
			s.wg.Wait()
			return
		}
	}
}

func (s *State) IncomingChan() chan<- *io.State {
	return s.inChan
}

func (s *State) ChangesChan() <-chan *io.Changes {
	return s.outChan
}

func (s *State) Publish(ctx context.Context, publish states.PublishFunc) error {
	done := make(chan error)
	select {
	case s.pubChan <- &pubPass{
		run:  publish,
		done: done,
	}:
		select {
		case e := <-done:
			return boo.WrapType(e, boo.Internal,
				"failed to publish root")

		case <-ctx.Done():
			return boo.WrapType(ctx.Err(), boo.Internal,
				"failed to publish root")
		}
	case <-ctx.Done():
		return boo.WrapType(ctx.Err(), boo.Internal,
			"failed to publish root")
	}
}

func (s *State) service() {
	s.wg.Add(1)
	defer s.wg.Done()

	for {
		select {
		case in := <-s.inChan:
			// Get rid of accumulated.
			if len(s.inChan) > 1 {
				in = <-s.inChan
			}
			// Process.
			if e := s.processIncoming(in); e != nil {
				s.l.Printf("Failed to process root %s[%d]",
					in.Root.Pub.Hex(), in.Root.Seq)
			}

		case pubPass := <-s.pubChan:
			pubPass.done <- pubPass.
				run(s.node, s.getCurrent().pack)

		case <-s.quitChan:
			return
		}
	}
}

func (s *State) processIncoming(in *io.State) error {
	pack, e := s.unpackRoot(in.Node, in.Root)
	if e != nil {
		return e
	}
	newCurrent, e := NewPackInstance(s.current, pack)
	if e != nil {
		return e
	}
	s.outChan <- s.
		replaceCurrent(newCurrent).
		changes

	return nil
}

func (s *State) getCurrent() *PackInstance {
	s.currentMux.Lock()
	defer s.currentMux.Unlock()
	return s.current
}

func (s *State) replaceCurrent(pi *PackInstance) *PackInstance {
	s.currentMux.Lock()
	defer s.currentMux.Unlock()
	s.current = pi
	return pi
}

func (s *State) GetBoardPage(ctx context.Context) (*io.BoardPageOut, error) {
	out := new(io.BoardPageOut)
	if e := s.getCurrent().Do(getBoardPage(out)); e != nil {
		return nil, e
	}
	return out, nil
}

func (s *State) GetThreadPage(ctx context.Context, tRef cipher.SHA256) (*io.ThreadPageOut, error) {
	out := new(io.ThreadPageOut)
	if e := s.getCurrent().Do(getThreadPage(out, tRef)); e != nil {
		return nil, e
	}
	return out, nil
}

func (s *State) GetFollowPage(ctx context.Context, upk cipher.PubKey) (*io.FollowPageOut, error) {
	out := new(io.FollowPageOut)
	e := s.getCurrent().Do(func(pi *PackInstance) error {
		out.Seq = pi.pack.Root().Seq
		out.BoardPubKey = pi.pack.Root().Pub
		fPage, e := pi.followStore.Get(upk)
		out.FollowPage = fPage
		return e
	})
	if e != nil {
		return nil, e
	}
	return out, nil
}

func (s *State) GetUserVotes(ctx context.Context, upk cipher.PubKey) (*io.VoteUserOut, error) {
	out := new(io.VoteUserOut)
	e := s.getCurrent().Do(func(pi *PackInstance) error {
		out.Seq = pi.pack.Root().Seq
		out.BoardPubKey = pi.pack.Root().Pub
		out.UserPubKey = upk
		vs, e := pi.uVotesStore.Get(upk)
		out.VoteSummary = vs
		return e
	})
	if e != nil {
		return nil, e
	}
	return out, nil
}

func (s *State) NewThread(ctx context.Context, thread *object.Content) (*io.BoardPageOut, error) {
	if e := thread.Verify(); e != nil {
		return nil, e
	}
	out := new(io.BoardPageOut)
	e := s.getCurrent().Do(func(pi *PackInstance) error {

		// Get thread ref.
		threadRef := pi.pack.Ref(thread)
		thread.R = threadRef.Hash

		// Check existence.
		if _, has := pi.gotStore.GetThread(threadRef.Hash); has {
			return boo.Newf(boo.AlreadyExists,
				"thread of hash '%s' already exists", threadRef.Hash)
		}

		// Append to thread pages.
		tPages, e := pi.GetThreadPages()
		if e != nil {
			return e
		}
		tPage := &object.ThreadPage{Thread: threadRef}
		if e := tPages.ThreadPages.Append(tPage); e != nil {
			return e
		}
		if e := pi.pack.SetRefByIndex(indexContent, tPages); e != nil {
			return e
		}

		// Add to GotStore.
		if e := pi.gotStore.AddThread(threadRef.Hash, tPage); e != nil {
			return e
		}

		// If no such VoteSummary, append it.
		if e := pi.AppendThreadVotesPage(threadRef.Hash); e != nil {
			return e
		}

		// Record changes.
		// TODO: Record new thread.

		// Prepare output.
		return getBoardPage(out)(pi)
	})
	if e != nil {
		return nil, e
	}
	return out, nil
}

func (s *State) NewPost(ctx context.Context, post *object.Content) (*io.ThreadPageOut, error) {
	if e := post.Verify(); e != nil {
		return nil, e
	}
	out := new(io.ThreadPageOut)
	e := s.getCurrent().Do(func(pi *PackInstance) error {

		// Get post ref.
		postRef := pi.pack.Ref(post)
		post.R = postRef.Hash

		// Check existence.
		if pi.gotStore.GetPostOrigin(postRef.Hash) != (cipher.SHA256{}) {
			return boo.Newf(boo.AlreadyExists,
				"post of hash '%s' already exists", postRef.Hash)
		}

		// Get thread info.
		gotThread, has := pi.gotStore.GetThread(post.Refer)
		if !has {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' is not found", post.Refer.Hex())
		}

		// Get thread.
		tPages, e := pi.GetThreadPages()
		if e != nil {
			return e
		}
		tPageRef, e := tPages.ThreadPages.RefByHash(gotThread.tPageHash)
		if e != nil {
			return e
		}
		tPage, e := obtain.ThreadPage(tPageRef)
		if e != nil {
			return e
		}
		thread, e := obtain.Content(&tPage.Thread)
		if e != nil {
			return e
		}

		// Submit post.
		if e := tPage.Posts.Append(post); e != nil {
			return e
		}
		if e := tPageRef.SetValue(tPage); e != nil {
			return e
		}
		if e := pi.pack.SetRefByIndex(indexContent, tPages); e != nil {
			return e
		}

		// Add to GotStore.
		if e := pi.gotStore.AddPost(thread.R, post.R); e != nil {
			return e
		}

		// If no such VoteSummary, append it.
		if e := pi.AppendPostVotesPage(postRef.Hash); e != nil {
			return e
		}

		// Record Changes.
		pi.changes.RecordNewPost(post.R, post)

		// Prepare output.
		return getThreadPage(out, thread.R)(pi)
	})
	if e != nil {
		return nil, e
	}
	return out, nil
}

func (s *State) DeleteThread(ctx context.Context, tRef cipher.SHA256) (*io.BoardPageOut, error) {
	out := new(io.BoardPageOut)
	e := s.getCurrent().Do(func(pi *PackInstance) error {

		// Check existence of thread.
		gotThread, has := pi.gotStore.GetThread(tRef)
		if !has {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' is not found", tRef.Hex())
		}

		// Get thread page.
		tPages, e := pi.GetThreadPages()
		if e != nil {
			return e
		}
		tPageRef, e := tPages.ThreadPages.RefByHash(gotThread.tPageHash)
		if e != nil {
			return e
		}
		tPage, e := obtain.ThreadPage(tPageRef)
		if e != nil {
			return e
		}

		// Get post votes pages.
		pvPages, e := pi.GetPostVotesPages()
		if e != nil {
			return e
		}

		// Delete posts.
		e = tPage.Posts.Ascend(func(_ int, pRef *skyobject.Ref) error {
			for i, vPage := range pvPages.Posts {
				if vPage.Ref == pRef.Hash {
					pvPages.Posts[0], pvPages.Posts[i] =
						pvPages.Posts[i], pvPages.Posts[0]
					pvPages.Posts = pvPages.Posts[1:]
					break
				}
			}
			pi.pVotesStore.Delete(pRef.Hash)
			return nil
		})
		if e != nil {
			return e
		}

		// Save post votes pages.
		if e := pi.pack.SetRefByIndex(indexPostVotes, pvPages); e != nil {
			return e
		}

		// Get thread votes pages.
		tvPages, e := pi.GetThreadVotesPages()
		if e != nil {
			return e
		}

		// Delete thread votes page.
		for i, tvPage := range tvPages.Threads {
			if tvPage.Ref == tRef {
				tvPages.Threads[0], tvPages.Threads[i] =
					tvPages.Threads[i], tvPages.Threads[0]
				tvPages.Threads = tvPages.Threads[1:]
				break
			}
		}
		pi.tVotesStore.Delete(tRef)
		pi.gotStore.DeleteThread(tRef)

		// Save thread votes pages.
		if e := pi.pack.SetRefByIndex(indexThreadVotes, tvPages); e != nil {
			return e
		}

		// Delete thread page.
		if e := tPageRef.Clear(); e != nil {
			return e
		}
		if e := pi.pack.SetRefByIndex(indexContent, tPages); e != nil {
			return e
		}

		// Record changes.
		pi.changes.RecordDeleteThread(tRef)

		// Prepare output.
		return getBoardPage(out)(pi)
	})
	if e != nil {
		return nil, e
	}
	return out, nil
}

func (s *State) DeletePost(ctx context.Context, pHash cipher.SHA256) (*io.ThreadPageOut, error) {
	out := new(io.ThreadPageOut)
	e := s.getCurrent().Do(func(pi *PackInstance) error {

		// Get post's origin.
		tHash := pi.gotStore.GetPostOrigin(pHash)
		if tHash == (cipher.SHA256{}) {
			return boo.Newf(boo.NotFound,
				"post of hash '%s' has no thread of origin", pHash.Hex())
		}

		gotThread, has := pi.gotStore.GetThread(tHash)
		if !has {
			return boo.Newf(boo.NotFound,
				"thread of hash '%s' is not found in GotStore", tHash.Hex())
		}

		// Delete post.
		tPages, e := pi.GetThreadPages()
		if e != nil {
			return e
		}
		tPageRef, e := tPages.ThreadPages.RefByHash(gotThread.tPageHash)
		if e != nil {
			return e
		}
		tPage, e := obtain.ThreadPage(tPageRef)
		if e != nil {
			return e
		}
		pRef, e := tPage.Posts.RefByHash(pHash)
		if e != nil {
			return e
		}
		if e := pRef.Clear(); e != nil {
			return e
		}
		if e := tPageRef.SetValue(tPage); e != nil {
			return e
		}
		if e := pi.pack.SetRefByIndex(indexContent, tPages); e != nil {
			return e
		}

		// Delete post votes.
		pvPages, e := pi.GetPostVotesPages()
		if e != nil {
			return e
		}
		for i, vPage := range pvPages.Posts {
			if vPage.Ref == pHash {
				pvPages.Posts[0], pvPages.Posts[i] =
					pvPages.Posts[i], pvPages.Posts[0]
				pvPages.Posts = pvPages.Posts[1:]
				break
			}
		}
		if e := pi.pack.SetRefByIndex(indexPostVotes, pvPages); e != nil {
			return e
		}

		// Delete compiled stuff.
		if e := pi.gotStore.DeletePost(pHash); e != nil {
			return e
		}
		pi.pVotesStore.Delete(pHash)

		// Record changes.
		pi.changes.RecordDeletePost(tHash, pHash)

		// Prepare output
		return getThreadPage(out, tHash)(pi)
	})
	if e != nil {
		return nil, e
	}
	return out, nil
}

func (s *State) VoteThread(ctx context.Context, vote *object.Vote) (*io.VoteThreadOut, error) {
	out := new(io.VoteThreadOut)
	e := s.getCurrent().Do(func(pi *PackInstance) error {
		tvPages, e := pi.GetThreadVotesPages()
		if e != nil {
			return e
		}
		for _, vPage := range tvPages.Threads {
			if vPage.Ref == vote.OfContent {
				// Loop through votes.
				e := vPage.Votes.Ascend(func(i int, oldVoteRef *skyobject.Ref) error {
					oldVote, e := obtain.Vote(oldVoteRef)
					if e != nil {
						return e
					}
					//if oldVote.
					return nil
				})
				if e != nil {
					return e
				}
			}
		}
		return nil
	})
	if e != nil {
		return nil, e
	}
	return out, nil
}

func (s *State) VotePost(ctx context.Context, vote *cipher.SHA256) (*io.VotePostOut, error) {
	return nil, nil
}

func (s *State) VoteUser(ctx context.Context, vote *cipher.SHA256) (*io.VoteUserOut, error) {
	return nil, nil
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func getBoardPage(out *io.BoardPageOut) func(pi *PackInstance) error {
	return func(pi *PackInstance) error {
		// Get initials.
		out.Seq = pi.pack.Root().Seq
		out.BoardPubKey = pi.pack.Root().Pub

		// Initiate thread pages.
		tPages, e := pi.GetThreadPages()
		if e != nil {
			return e
		}
		tPagesLen, e := tPages.ThreadPages.Len()
		if e != nil {
			return e
		}
		out.Threads = make([]*object.Content, tPagesLen)
		out.ThreadVotes = make([]*object.VotesSummary, tPagesLen)

		// Get thread pages.
		return tPages.ThreadPages.Ascend(func(i int, tPageRef *skyobject.Ref) error {
			tPage, e := obtain.ThreadPage(tPageRef)
			if e != nil {
				return e
			}
			thread, e := obtain.Content(&tPage.Thread)
			if e != nil {
				return e
			}
			out.Threads[i] = thread
			vs, e := pi.tVotesStore.Get(thread.R)
			out.ThreadVotes[i] = vs
			return e
		})
	}
}

func getThreadPage(out *io.ThreadPageOut, tRef cipher.SHA256) func(pi *PackInstance) error {
	return func(pi *PackInstance) error {
		// Get initials.
		out.Seq = pi.pack.Root().Seq
		out.BoardPubKey = pi.pack.Root().Pub

		// Get thread pages.
		tPages, e := pi.GetThreadPages()
		if e != nil {
			return e
		}
		tPageRef, e := tPages.ThreadPages.RefByHash(tRef)
		if e != nil {
			return boo.WrapTypef(e, boo.NotFound,
				"thread of hash %s is not found", tRef.Hex())
		}

		// Get thread page.
		tPage, e := obtain.ThreadPage(tPageRef)
		if e != nil {
			return boo.WrapType(e, boo.InvalidRead,
				"ThreadPage is corrupt")
		}

		// Get thread.
		thread, e := obtain.Content(&tPage.Thread)
		if e != nil {
			return e
		}
		out.Thread = thread
		vs, e := pi.tVotesStore.Get(thread.R)
		out.ThreadVote = vs
		if e != nil {
			return e
		}

		// Initiate posts.
		postsLen, e := tPage.Posts.Len()
		if e != nil {
			return e
		}
		out.Posts = make([]*object.Content, postsLen)
		out.PostVotes = make([]*object.VotesSummary, postsLen)

		// Get posts.
		return tPage.Posts.Ascend(func(i int, pRef *skyobject.Ref) error {
			post, e := obtain.Content(pRef)
			if e != nil {
				return e
			}
			out.Posts[i] = post
			vs, e := pi.pVotesStore.Get(post.R)
			out.PostVotes[i] = vs
			return e
		})
	}
}
