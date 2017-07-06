package gui

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/msg"
	"github.com/skycoin/bbs/src/store/typ"
	"net/http"
	"time"
)

// Gateway represents the intermediate between External calls and internal processing.
// It can be seen as a security layer.
type Gateway struct {
	config     *HTTPConfig
	container  *store.CXO
	boardSaver *store.BoardSaver
	userSaver  *store.UserSaver
	queueSaver *msg.QueueSaver
	quitChan   chan int

	Stats
	Connections
	Subscriptions
	Users
	Boards
	Threads
	Posts
	Tests
}

// NewGateway creates a new Gateway.
func NewGateway(
	cf *HTTPConfig, ct *store.CXO,
	bs *store.BoardSaver, us *store.UserSaver, qs *msg.QueueSaver,
	q chan int,
) *Gateway {
	g := &Gateway{config: cf, container: ct, boardSaver: bs, userSaver: us, queueSaver: qs, quitChan: q}
	g.Stats.Gateway = g
	g.Connections.Gateway = g
	g.Subscriptions.Gateway = g
	g.Users.Gateway = g
	g.Users.Masters.Gateway = g
	g.Users.Masters.Current.Gateway = g
	g.Boards.Gateway = g
	g.Boards.Meta.Gateway = g
	g.Boards.Meta.SubmissionAddresses.Gateway = g
	g.Boards.Page.Gateway = g
	g.Threads.Gateway = g
	g.Threads.Page.Gateway = g
	g.Threads.Votes.Gateway = g
	g.Posts.Gateway = g
	g.Posts.Votes.Gateway = g
	g.Tests.Gateway = g
	return g
}

// Quit quits the node entirely.
func (g *Gateway) Quit(w http.ResponseWriter, r *http.Request) {
	send(w, g.quit(), http.StatusOK)
}

func (g *Gateway) quit() bool {
	timer := time.NewTimer(10 * time.Second)
	select {
	case g.quitChan <- 0:
		return true
	case <-timer.C:
		return false
	}
}

// PingSubmissionAddress pings a submission address.
func (g *Gateway) pingSubmissionAddress(address string) error {
	return g.queueSaver.Ping(address)
}

/*
	<<< VIEWS >>>
*/

// StatsView represents the stats json structure as displayed to end user.
type StatsView struct {
	NodeIsMaster   bool   `json:"node_is_master"`
	NodeCXOAddress string `json:"node_cxo_address"`
}

// BoardPageView represents a board page as json as displayed to end user.
type BoardPageView struct {
	Board   *typ.Board    `json:"board"`
	Threads []*typ.Thread `json:"threads"`
}

// ThreadPageView represents a thread page as json when displayed to end user,
type ThreadPageView struct {
	Board  *typ.Board  `json:"board"`
	Thread *typ.Thread `json:"thread"`
	Posts  []*typ.Post `json:"posts"`
}

// VotesView represents a votes view as json when displayed to end user.
type VotesView struct {
	UpVotes             int  `json:"up_votes"`
	DownVotes           int  `json:"down_votes"`
	CurrentUserVoted    bool `json:"current_user_voted"`
	CurrentUserVoteMode int  `json:"current_user_vote_mode,omitempty"`
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func send(w http.ResponseWriter, v interface{}, httpStatus int) error {
	w.Header().Set("Content-Type", "application/json")
	respData, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.WriteHeader(httpStatus)
	w.Write(respData)
	return nil
}
