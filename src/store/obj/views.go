package obj

type BoardView struct {
	Board
	PublicKey     string              `json:"public_key"`
	ExternalRoots []*ExternalRootView `json:"external_roots,omitempty"`
	Threads       []*ThreadView       `json:"threads,omitempty"`
}

type ExternalRootView struct {
	ExternalRoot
	PublicKey string `json:"public_key"`
}

type ThreadView struct {
	Thread
	Ref            string       `json:"reference"`
	AuthorRef      string       `json:"author_reference,omitempty"`
	AuthorAlias    string       `json:"author_alias,omitempty"`
	MasterBoardRef string       `json:"master_board_reference"`
	Votes          *VoteSummary `json:"votes,omitempty"`
	Posts          []*PostView  `json:"posts,omitempty"`
}

type PostView struct {
	Post
	Ref         string       `json:"reference"`
	AuthorRef   string       `json:"author_reference,omitempty"`
	AuthorAlias string       `json:"author_alias,omitempty"`
	Votes       *VoteSummary `json:"votes"`
}

type VoteSummary struct {
	Up   VoteView `json:"up"`
	Down VoteView `json:"down"`
	Spam VoteView `json:"spam"`
}

type VoteView struct {
	Voted bool `json:"voted"`
	Count int  `json:"count"`
}

type UserView struct {
	User
	PublicKey string       `json:"public_key,omitempty"`
	SecretKey string       `json:"secret_key,omitempty"`
	Votes     *VoteSummary `json:"votes"`
}

type SubscriptionView struct {
	PubKey string `json:"public_key"`
	SecKey string `json:"secret_key,omitempty"`
}

type ConnectionView struct {
	Address string `json:"address"`
	Active  bool   `json:"active"`
}
