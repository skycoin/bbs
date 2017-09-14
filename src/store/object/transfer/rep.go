package transfer

import "github.com/skycoin/bbs/src/misc/boo"

type RootRep struct {
	Root *RootPageRep `json:"root"`
	Board *BoardPageRep   `json:"board"`
}

func (r *RootRep) Fill(rp RootPage, bp BoardPage) error {
	var e error
	if r.Root, e = rp.Export(); e != nil {
		return e
	}
	r.Board = new(BoardPageRep)
	return r.Board.Fill(bp)
}

type RootPageRep struct {
	Type     string          `json:"root_type"`     // Type of root (eg. Board root).
	Revision uint64          `json:"root_revision"` // Revision of root type (eg. Board root r2).
	Deleted  bool            `json:"deleted"`       // Whether root is deleted.
	Summary  RootPageSummary `json:"summary,string"`
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
	if b, e := bp.DumpBoard(); e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to obtain board")
	} else {
		if r.Board, e = b.Export(); e != nil {
			return boo.WrapType(e, boo.InvalidRead,
				"failed to export board")
		}
	}
	if ts, e := bp.DumpThreadPages(); e != nil {
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

type ThreadPageRep struct {
	Thread *ThreadRep `json:"thread"`
	Posts  []*PostRep `json:"posts"`
}

func (r *ThreadPageRep) Fill(tp ThreadPage) error {
	if t, e := tp.DumpThread(); e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to obtain thread")
	} else {
		if r.Thread, e = t.Export(); e != nil {
			return boo.WrapType(e, boo.InvalidRead,
				"failed to export thread")
		}
	}
	if ps, e := tp.DumpPosts(); e != nil {
		return boo.WrapType(e, boo.InvalidRead,
			"failed to obtain posts")
	} else {
		r.Posts = make([]*PostRep, len(ps))
		for i, p := range ps {
			if r.Posts[i], e = p.Export(); e != nil {
				return boo.WrapTypef(e, boo.InvalidRead,
					"failed to export post at index %d", i)
			}
		}
	}
	return nil
}

type BoardRep struct {
	Name    string   `json:"name"`
	Body    string   `json:"body"`
	Created int64    `json:"created"`
	Tags    []string `json:"tags"`
}

type ThreadRep struct {
	Name    string `json:"name"`
	Body    string `json:"body"`
	Created int64  `json:"created"`
	Creator string `json:"creator"`
}

type PostRep struct {
	OfPost  string `json:"of_post"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Created int64  `json:"created"`
	Creator string `json:"creator"`
}