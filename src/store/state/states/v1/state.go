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
	var out = new(io.BoardPageOut)
	var e = s.getCurrent().Do(func(pi *PackInstance) error {

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
	})
	return out, e
}

func (s *State) GetThreadPage(ctx context.Context, tRef cipher.SHA256) (*io.ThreadPageOut, error) {
	var out = new(io.ThreadPageOut)
	var e = s.getCurrent().Do(func(pi *PackInstance) error {

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
	})
	return out, e
}

func (s *State) GetFollowPage(ctx context.Context, upk cipher.PubKey) (*io.FollowPageOut, error) {
	var out = new(io.FollowPageOut)
	var e = s.current.Do(func(pi *PackInstance) error {

		return nil
	})
	return out, e
}

func (s *State) GetUserVotes(ctx context.Context, upk cipher.PubKey) (*io.VoteUserOut, error) {
	return nil, nil
}

func (s *State) NewThread(ctx context.Context, thread *object.Content) (*io.BoardPageOut, error) {
	return nil, nil
}

func (s *State) NewPost(ctx context.Context, tRef cipher.SHA256, post *object.Content) (*io.ThreadPageOut, error) {
	return nil, nil
}

func (s *State) DeleteThread(ctx context.Context, tRef cipher.SHA256) (*io.BoardPageOut, error) {
	return nil, nil
}

func (s *State) DeletePost(ctx context.Context, tRef, pRef cipher.SHA256) (*io.ThreadPageOut, error) {
	return nil, nil
}

func (s *State) VoteThread(ctx context.Context, tRef cipher.SHA256, vote *cipher.SHA256) (*io.VoteThreadOut, error) {
	return nil, nil
}

func (s *State) VotePost(ctx context.Context, tRef, pRef cipher.SHA256, vote *cipher.SHA256) (*io.VotePostOut, error) {
	return nil, nil
}

func (s *State) VoteUser(ctx context.Context, upk cipher.PubKey, vote *cipher.SHA256) (*io.VoteUserOut, error) {
	return nil, nil
}
