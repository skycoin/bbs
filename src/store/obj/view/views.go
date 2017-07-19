package view

import "github.com/skycoin/bbs/src/store/obj"

type BoardView struct {
	obj.Board
	PublicKey     string             `json:"public_key"`
	ExternalRoots []ExternalRootView `json:"external_roots"`
	Threads       []ThreadView       `json:"threads"`
	// TODO: BoardMeta part.
}

type ExternalRootView struct {
	obj.ExternalRoot
	PublicKey string `json:"public_key"`
}

type ThreadView struct {
	obj.Thread
	Ref            string     `json:"reference"`
	AuthorRef      string     `json:"author_reference,omitempty"`
	AuthorAlias    string     `json:"author_alias,omitempty"`
	MasterBoardRef string     `json:"master_board_reference"`
	Posts          []PostView `json:"posts"`
	// TODO: ThreadMeta part.
}

type PostView struct {
	obj.Post
	Ref         string `json:"reference"`
	AuthorRef   string `json:"author_reference,omitempty"`
	AuthorAlias string `json:"author_alias,omitempty"`
	// TODO: PostMeta part.
	// TODO: Votes.
}
