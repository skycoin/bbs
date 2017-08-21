package content_view

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
)

const (
	Board      = "Board"
	BoardPage  = "BoardPage"
	ThreadPage = "ThreadPage"
)

func (v *ContentView) Get(id string, a ...interface{}) (object.Lockable, error) {
	switch {
	case id == Board:
		return v.getBoard()

	case id == BoardPage:
		return v.getBoardPage()

	case id == ThreadPage && len(a) == 1:
		return v.getThreadPage(a[0].(cipher.SHA256))

	default:
		return nil, boo.Newf(boo.NotAllowed,
			"invalid get request 's' (%v)", id, a)
	}
}

func (v *ContentView) getBoard() (*BoardRep, error) {
	return v.board, nil
}

type BoardPageOut struct {
	Board   *BoardRep    `json:"board"`
	Threads []*ThreadRep `json:"threads"`
}

func (o *BoardPageOut) Lock() {
	if o.Board != nil {
		o.Board.Lock()
	}
	for _, t := range o.Threads {
		if t != nil {
			t.Lock()
		}
	}
}

func (o *BoardPageOut) Unlock() {
	if o.Board != nil {
		o.Board.Unlock()
	}
	for _, t := range o.Threads {
		if t != nil {
			t.Unlock()
		}
	}
}

func (v *ContentView) getBoardPage() (*BoardPageOut, error) {
	out := new(BoardPageOut)
	out.Board = v.board
	out.Threads = make([]*ThreadRep, len(v.board.Threads))
	for i, tHash := range v.board.Threads {
		out.Threads[i] = v.tMap[tHash]
	}
	return out, nil
}

type ThreadPageOut struct {
	Board  *BoardRep  `json:"board"`
	Thread *ThreadRep `json:"thread"`
	Posts  []*PostRep `json:"posts"`
}

func (o *ThreadPageOut) Lock() {
	if o.Board != nil {
		o.Board.Lock()
	}
	if o.Thread != nil {
		o.Thread.Lock()
	}
	for _, p := range o.Posts {
		if p != nil {
			p.Lock()
		}
	}
}

func (o *ThreadPageOut) Unlock() {
	if o.Board != nil {
		o.Board.Unlock()
	}
	if o.Thread != nil {
		o.Thread.Unlock()
	}
	for _, p := range o.Posts {
		if p != nil {
			p.Unlock()
		}
	}
}

func (v *ContentView) getThreadPage(threadHash cipher.SHA256) (*ThreadPageOut, error) {
	out := new(ThreadPageOut)
	out.Board = v.board
	if out.Thread = v.tMap[threadHash]; out.Thread != nil {
		out.Posts = make([]*PostRep, len(out.Thread.Posts))
		for i, pHash := range out.Thread.Posts {
			out.Posts[i] = v.pMap[pHash]
		}
	}
	return out, nil
}
