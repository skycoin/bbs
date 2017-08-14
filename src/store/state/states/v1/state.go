package v1

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/io"
	"github.com/skycoin/bbs/src/store/state/states"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/skyobject"
	"log"
	"os"
	"sync"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/bbs/src/store/object"
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
	c         *states.StateConfig
	l         *log.Logger
	flag      skyobject.Flag
	node      *node.Node
	packStore *PackStore
	inChan    chan *io.State   // Need to generate.
	outChan   chan *io.Changes // Need to generate.
	pubChan   chan *pubPass    // Need to generate.
	quitChan  chan struct{}    // Need to generate.
	wg        sync.WaitGroup
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
			"unable to obtain root")
	}
	pack, e := s.unpackRoot(node, root)
	if e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"unable to unpack root")
	}
	s.packStore = NewPackStore(s.l)
	if _, e := s.packStore.Extract(pack, root.Seq, true, false); e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to extract board root")
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
			if e := s.processIncoming(in, true); e != nil {
				s.l.Printf("Failed to process root %s[%d]",
					in.Root.Pub.Hex(), in.Root.Seq)
			}
		case pubPass := <-s.pubChan:
			pubPass.done <- s.packStore.Run(
				func(pack *skyobject.Pack, seq uint64) error {
					return pubPass.run(s.node, pack)
				},
			)
		case <-s.quitChan:
			return
		}
	}
}

func (s *State) processIncoming(in *io.State, full bool) error {
	pack, e := s.unpackRoot(in.Node, in.Root)
	if e != nil {
		return e
	}
	changes, e := s.packStore.Extract(pack, in.Root.Seq, full, true)
	if e != nil {
		return e
	}
	s.outChan <- changes
	return nil
}

func (s *State) GetBoardPage(ctx context.Context) (*io.BoardPageOut, error) {
	return nil, nil
}

func (s *State) GetThreadPage(ctx context.Context, tRef cipher.SHA256) (*io.ThreadPageOut, error) {
	return nil, nil
}

func (s *State) GetFollowPage(ctx context.Context, upk cipher.PubKey) (*io.FollowPageOut, error) {
	return nil, nil
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













