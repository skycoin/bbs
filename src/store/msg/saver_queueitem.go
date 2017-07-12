package msg

import (
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/rpc"
	"time"
)

// QueueItem represents an incomplete action in queue.
type QueueItem struct {
	ID              string             `json:"id"`
	Submitted       int64              `json:"submitted"`
	GetBoardFails   int                `json:"get_board_fails"`
	ConnectionFails int                `json:"connection_fails"`
	ReqNewPost      *rpc.ReqNewPost    `json:"new_post_request,omitempty"`
	ReqNewThread    *rpc.ReqNewThread  `json:"new_thread_request,omitempty"`
	ReqVotePost     *rpc.ReqVotePost   `json:"vote_post_request,omitempty"`
	ReqVoteThread   *rpc.ReqVoteThread `json:"vote_thread_request,omitempty"`
	ReqVoteUser     *rpc.ReqVoteUser   `json:"vote_user_request,omitempty"`
}

// NewQueueItem creates a new QueueItem.
func NewQueueItem() *QueueItem {
	return &QueueItem{
		ID:        misc.MakeTimeStampedRandomID(100).Hex(),
		Submitted: time.Now().UnixNano(),
	}
}

// Clear clears all fields of QueueItem.
func (qi *QueueItem) Clear() {
	qi.ReqNewPost = nil
	qi.ReqNewThread = nil
	qi.ReqVotePost = nil
	qi.ReqVoteThread = nil
	qi.ReqVoteUser = nil
}

func (qi *QueueItem) SetReqNewPost(req *rpc.ReqNewPost) *QueueItem {
	qi.ReqNewPost = req
	return qi
}

func (qi *QueueItem) SetReqNewThread(req *rpc.ReqNewThread) *QueueItem {
	qi.ReqNewThread = req
	return qi
}

func (qi *QueueItem) SetReqVotePost(req *rpc.ReqVotePost) *QueueItem {
	qi.ReqVotePost = req
	return qi
}

func (qi *QueueItem) SetReqVoteThread(req *rpc.ReqVoteThread) *QueueItem {
	qi.ReqVoteThread = req
	return qi
}

func (qi *QueueItem) SetReqVoteUser(req *rpc.ReqVoteUser) *QueueItem {
	qi.ReqVoteUser = req
	return qi
}
