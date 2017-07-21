package access

import (
	"encoding/json"
	"fmt"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store/obj"
	"github.com/skycoin/bbs/src/store/obj/view"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"sync"
)

const rootQueueCap = 1

// State represents a state of a board.
type State struct {
	mux       sync.Mutex
	board     view.BoardView
	boardRoot chan *skyobject.Root
	quit      chan struct{}
}

// NewState creates a new state.
func NewState() *State {
	state := &State{
		boardRoot: make(chan *skyobject.Root, rootQueueCap),
		quit:      make(chan struct{}),
	}
	go state.service()
	return state
}

// Close closes the state service.
func (s *State) Close() {
	s.quit <- struct{}{}
}

// PushNewBoardRoot pushes a new board root to process.
func (s *State) PushNewBoardRoot(root *skyobject.Root) {
	for {
		select {
		case <-s.boardRoot:
		default:
			s.boardRoot <- root
			return
		}
	}
}

// GetView obtains the current state view.
func (s *State) GetView() view.BoardView {
	s.mux.Lock()
	defer s.mux.Unlock()
	return s.board
}

func (s *State) service() {
	for {
		select {
		case r := <-s.boardRoot:
			if boardView, e := processBoardRoot(r); e != nil {
				fmt.Println("[STATE] Error on process board root:", e)
			} else {
				s.mux.Lock()
				s.board = *boardView
				data, _ := json.MarshalIndent(s.board, "", "    ")
				fmt.Println(string(data))
				s.mux.Unlock()
			}
		case <-s.quit:
			return
		}
	}
}

func processBoardRoot(r *skyobject.Root) (*view.BoardView, error) {
	schemaRefs, e := misc.GetSchemaRefsFromRoot(r)
	if e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead, "failed to obtain schema references")
	}

	for _, dyn := range r.Refs() {
		switch dyn.Schema {
		case schemaRefs.BoardPage:
			return getBoardView(r, dyn.Object)
		}
	}
	return nil, boo.New(boo.InvalidRead,
		"root has no readable elements")
}

func getBoardPage(r *skyobject.Root, ref skyobject.Reference) (*obj.BoardPage, error) {
	data, has := r.Get(ref)
	if !has {
		return nil, boo.Newf(boo.ObjectNotFound,
			"board page of ref '%s' not found", ref.String())
	}
	boardPage := new(obj.BoardPage)
	if e := encoder.DeserializeRaw(data, boardPage); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead, "failed to deserialize board page")
	}
	return boardPage, nil
}

func getBoard(r *skyobject.Root, ref skyobject.Reference) (*obj.Board, error) {
	data, has := r.Get(ref)
	if !has {
		return nil, boo.Newf(boo.ObjectNotFound, "board of ref '%s' not found", ref.String())
	}
	board := new(obj.Board)
	if e := encoder.DeserializeRaw(data, board); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead, "failed to deserialize board")
	}
	return board, nil
}

func getExtRootsView(board *obj.Board) ([]view.ExternalRootView, error) {
	extRootViews := make([]view.ExternalRootView, len(board.ExternalRoots))
	for i, extRoot := range board.ExternalRoots {
		extRootViews[i] = view.ExternalRootView{
			ExternalRoot: extRoot,
			PublicKey:    extRoot.PublicKey.Hex(),
		}
	}
	return extRootViews, nil
}

func getBoardView(r *skyobject.Root, ref skyobject.Reference) (*view.BoardView, error) {
	boardPage, e := getBoardPage(r, ref)
	if e != nil {
		return nil, e
	}
	board, e := getBoard(r, boardPage.Board)
	if e != nil {
		return nil, e
	}
	extRootsView, e := getExtRootsView(board)
	if e != nil {
		return nil, e
	}
	threadViews, e := getThreadViews(r, boardPage.ThreadPages)
	if e != nil {
		return nil, e
	}
	boardView := &view.BoardView{
		Board:         *board,
		PublicKey:     r.Pub().Hex(),
		ExternalRoots: extRootsView,
		Threads:       threadViews,
	}
	return boardView, nil
}

func getThreadPage(r *skyobject.Root, ref skyobject.Reference) (*obj.ThreadPage, error) {
	data, has := r.Get(ref)
	if !has {
		return nil, boo.Newf(boo.ObjectNotFound, "thread page of ref '%s' not found", ref.String())
	}
	threadPage := new(obj.ThreadPage)
	if e := encoder.DeserializeRaw(data, threadPage); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead, "failed to deserialize thread page")
	}
	return threadPage, nil
}

func getThread(r *skyobject.Root, ref skyobject.Reference) (*obj.Thread, error) {
	data, has := r.Get(ref)
	if !has {
		return nil, boo.Newf(boo.ObjectNotFound,
			"thread of ref '%s' not found", ref.String())
	}
	thread := new(obj.Thread)
	if e := encoder.DeserializeRaw(data, thread); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead, "failed to deserialize thread")
	}
	return thread, nil
}

func getThreadViews(r *skyobject.Root, refs skyobject.References) ([]view.ThreadView, error) {
	tViews := make([]view.ThreadView, len(refs))
	for i, ref := range refs {
		threadPage, e := getThreadPage(r, ref)
		if e != nil {
			return nil, e
		}
		thread, e := getThread(r, threadPage.Thread)
		if e != nil {
			return nil, e
		}
		postViews, e := getPostViews(r, threadPage.Posts)
		if e != nil {
			return nil, e
		}
		tViews[i].Thread = *thread
		tViews[i].Ref = threadPage.Thread.String()
		tViews[i].AuthorRef = thread.User.Hex()
		tViews[i].AuthorAlias = "" // TODO: Implement.
		tViews[i].MasterBoardRef = thread.MasterBoardRef.String()
		tViews[i].Posts = postViews
	}
	return tViews, nil
}

func getPost(r *skyobject.Root, ref skyobject.Reference) (*obj.Post, error) {
	data, has := r.Get(ref)
	if !has {
		return nil, boo.Newf(boo.ObjectNotFound, "post of ref '%s' not found", ref.String())
	}
	post := new(obj.Post)
	if e := encoder.DeserializeRaw(data, post); e != nil {
		return nil, boo.WrapType(e, boo.InvalidRead, "failed to deserialize post")
	}
	return post, nil
}

func getPostViews(r *skyobject.Root, refs skyobject.References) ([]view.PostView, error) {
	pViews := make([]view.PostView, len(refs))
	for i, ref := range refs {
		post, e := getPost(r, ref)
		if e != nil {
			return nil, e
		}
		pViews[i].Post = *post
		pViews[i].Ref = ref.String()
		pViews[i].AuthorRef = post.User.Hex()
		pViews[i].AuthorAlias = "" // TODO: Implement.
		// TODO: PostView meta and votes.
	}
	return pViews, nil
}
