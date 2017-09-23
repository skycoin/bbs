package transfer

import "github.com/skycoin/cxo/skyobject"

type Generator interface {
	Pack() *skyobject.Pack
	NewRootPage() RootPage
	NewBoardPage() BoardPage
	NewThreadPage() ThreadPage
	NewBoard() Board
	NewThread() Thread
	NewPost() Post
}

type RootPage interface {
	ToRep() (*RootPageRep, error)
	FromRep(rpRep *RootPageRep) error
}

type BoardPage interface {
	ExportBoard() (Board, error)
	ExportThreadPages() ([]ThreadPage, error)

	ImportBoard(p *skyobject.Pack, b Board) error
	ImportThreadPages(p *skyobject.Pack, tps []ThreadPage) error
}

type ThreadPage interface {
	ExportThread() (Thread, error)
	ExportPosts() ([]Post, error)

	ImportThread(p *skyobject.Pack, t Thread) error
	ImportPosts(p *skyobject.Pack, ps []Post) error
}

type Board interface {
	ToRep() (*BoardRep, error)
	FromRep(bRep *BoardRep) error
}

type Thread interface {
	ToRep() (*ThreadRep, error)
	FromRep(tRep *ThreadRep) error
}

type Post interface {
	ToRep() (*PostRep, error)
	FromRep(pRep *PostRep) error
}
