package transfer

import (
	"github.com/skycoin/bbs/src/misc/boo"
)

type RootRep struct {
	Root  *RootPageRep  `json:"root"`
	Board *BoardPageRep `json:"board"`
}

func (r *RootRep) Fill(rp RootPage, bp BoardPage) error {
	var e error
	if r.Root, e = rp.ToRep(); e != nil {
		return e
	}
	r.Board = new(BoardPageRep)
	return r.Board.Fill(bp)
}

func (r *RootRep) Dump(g Generator) error {
	rp, e := r.Root.Dump(g)
	if e != nil {
		return e
	}
	if e := g.Pack().SetRefByIndex(0, rp); e != nil {
		return e
	}
	bp, e := r.Board.Dump(g)
	if e != nil {
		return e
	}
	if e := g.Pack().SetRefByIndex(1, bp); e != nil {
		return e
	}
	return nil
}

type RootPageRep struct {
	Type     string          `json:"root_type"`     // Type of root (eg. Board root).
	Revision uint64          `json:"root_revision"` // Revision of root type (eg. Board root r2).
	Deleted  bool            `json:"deleted"`       // Whether root is deleted.
	Summary  RootPageSummary `json:"summary,string"`
}

func (r *RootPageRep) Dump(g Generator) (RootPage, error) {
	out := g.NewRootPage()
	if e := out.FromRep(r); e != nil {
		return nil, e
	}
	return out, nil
}

type RootPageSummary struct {
	Name    string   `json:"name"`
	Body    string   `json:"body"`
	Created int64    `json:"created"`
	Tags    []string `json:"tags"`
}

type BoardPageRep struct {
	Board   *BoardRep        `json:"board"`
	Threads []*ThreadPageRep `json:"threads"`
}

func (r *BoardPageRep) Fill(bp BoardPage) error {
	if b, e := bp.ExportBoard(); e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to obtain board")
	} else {
		if r.Board, e = b.ToRep(); e != nil {
			return boo.WrapType(e, boo.InvalidRead,
				"failed to export board")
		}
	}
	if ts, e := bp.ExportThreadPages(); e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to obtain threads")
	} else {
		r.Threads = make([]*ThreadPageRep, len(ts))
		for i, tp := range ts {
			r.Threads[i] = new(ThreadPageRep)
			if e := r.Threads[i].Fill(tp); e != nil {
				return boo.WrapTypef(e, boo.InvalidRead,
					"failed to fill thread page at index %d", i)
			}
		}
	}
	return nil
}

func (r *BoardPageRep) Dump(g Generator) (BoardPage, error) {
	out := g.NewBoardPage()
	g.Pack().Ref(out)

	// Get board.
	b, e := r.Board.Dump(g)
	if e != nil {
		return nil, e
	}
	if e := out.ImportBoard(g.Pack(), b); e != nil {
		return nil, e
	}

	// Get threads.
	ts := make([]ThreadPage, len(r.Threads))
	for i, tRep := range r.Threads {
		if ts[i], e = tRep.Dump(g); e != nil {
			return nil, e
		}
	}
	if e := out.ImportThreadPages(g.Pack(), ts); e != nil {
		return nil, e
	}

	// Return.
	return out, nil
}

type ThreadPageRep struct {
	Thread *ThreadRep `json:"thread"`
	Posts  []*PostRep `json:"posts"`
}

func (r *ThreadPageRep) Fill(tp ThreadPage) error {
	if t, e := tp.ExportThread(); e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to obtain thread")
	} else {
		if r.Thread, e = t.ToRep(); e != nil {
			return boo.WrapType(e, boo.InvalidRead,
				"failed to export thread")
		}
	}
	if ps, e := tp.ExportPosts(); e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to obtain posts")
	} else {
		r.Posts = make([]*PostRep, len(ps))
		for i, p := range ps {
			if r.Posts[i], e = p.ToRep(); e != nil {
				return boo.WrapTypef(e, boo.InvalidRead,
					"failed to export post at index %d", i)
			}
		}
	}
	return nil
}

func (r *ThreadPageRep) Dump(g Generator) (ThreadPage, error) {
	out := g.NewThreadPage()
	g.Pack().Ref(out)

	// Get thread.
	t, e := r.Thread.Dump(g)
	if e != nil {
		return nil, e
	}
	if e := out.ImportThread(g.Pack(), t); e != nil {
		return nil, e
	}

	// Get posts.
	ps := make([]Post, len(r.Posts))
	for i, pRep := range r.Posts {
		if ps[i], e = pRep.Dump(g); e != nil {
			return nil, e
		}
	}
	if e := out.ImportPosts(g.Pack(), ps); e != nil {
		return nil, e
	}

	// Return.
	return out, nil
}

type BoardRep struct {
	Name    string   `json:"name"`
	Body    string   `json:"body"`
	Created int64    `json:"created"`
	Tags    []string `json:"tags"`
}

func (r *BoardRep) Dump(g Generator) (Board, error) {
	out := g.NewBoard()
	if e := out.FromRep(r); e != nil {
		return nil, e
	}
	return out, nil
}

type ThreadRep struct {
	Name    string `json:"name"`
	Body    string `json:"body"`
	Created int64  `json:"created"`
	Creator string `json:"creator"`
}

func (r *ThreadRep) Dump(g Generator) (Thread, error) {
	out := g.NewThread()
	if e := out.FromRep(r); e != nil {
		return nil, e
	}
	return out, nil
}

type PostRep struct {
	OfPost  string `json:"of_post"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Created int64  `json:"created"`
	Creator string `json:"creator"`
}

func (r *PostRep) Dump(g Generator) (Post, error) {
	out := g.NewPost()
	if e := out.FromRep(r); e != nil {
		return nil, e
	}
	return out, nil
}
