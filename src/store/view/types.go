package view

import "github.com/skycoin/bbs/src/store/obj"

type Board struct {
	*obj.Board
	Ref             string    `json:"reference"`
	ContentVotesRef string    `json:"content_votes_reference"`
	UserVotesRef    string    `json:"user_votes_reference"`
	Threads         []*Thread `json:"thread"`
	// TODO: BoardMeta part.
}

type Thread struct {
	*obj.Thread
	AuthorRef      string  `json:"author_reference,omitempty"`
	AuthorAlias    string  `json:"author_alias,omitempty"`
	MasterBoardRef string  `json:"master_board_reference"`
	Posts          []*Post `json:"posts"`
	// TODO: ThreadMeta part.
}

type Post struct {
	*obj.Post
	AuthorRef   string `json:"author_reference,omitempty"`
	AuthorAlias string `json:"author_alias,omitempty"`
	// TODO: PostMeta part.
	// TODO: Votes.
}
