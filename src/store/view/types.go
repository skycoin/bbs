package view

import "github.com/skycoin/bbs/src/store/obj"

type Board struct {
	obj.Board
	PublicKey     string         `json:"public_key"`
	ExternalRoots []ExternalRoot `json:"external_roots"`
	Threads       []Thread       `json:"threads"`
	// TODO: BoardMeta part.
}

type ExternalRoot struct {
	obj.ExternalRoot
	PublicKey string `json:"public_key"`
}

type Thread struct {
	obj.Thread
	Ref            string `json:"reference"`
	AuthorRef      string `json:"author_reference,omitempty"`
	AuthorAlias    string `json:"author_alias,omitempty"`
	MasterBoardRef string `json:"master_board_reference"`
	Posts          []Post `json:"posts"`
	// TODO: ThreadMeta part.
}

type Post struct {
	obj.Post
	Ref         string `json:"reference"`
	AuthorRef   string `json:"author_reference,omitempty"`
	AuthorAlias string `json:"author_alias,omitempty"`
	// TODO: PostMeta part.
	// TODO: Votes.
}
