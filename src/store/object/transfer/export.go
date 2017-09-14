package transfer

type RootPage interface {
	Export() (*RootPageRep, error)
}

type BoardPage interface {
	DumpBoard() (Board, error)
	DumpThreadPages() ([]ThreadPage, error)
}

type ThreadPage interface {
	DumpThread() (Thread, error)
	DumpPosts() ([]Post, error)
}

type Board interface {
	Export() (*BoardRep, error)
}

type Thread interface {
	Export() (*ThreadRep, error)
}

type Post interface {
	Export() (*PostRep, error)
}

