package msg

import (
	"github.com/evanlinjin/bbs/extern/rpc"
	"github.com/evanlinjin/bbs/misc"
	"time"
)

type QueueItem struct {
	ID              string            `json:"id"`
	Submitted       int64             `json:"submitted"`
	GetBoardFails   int               `json:"get_board_fails"`
	ConnectionFails int               `json:"connection_fails"`
	ReqNewPost      *rpc.ReqNewPost   `json:"new_post_request,omitempty"`
	ReqNewThread    *rpc.ReqNewThread `json:"new_thread_request,omitempty"`
}

func NewQueueItem() *QueueItem {
	return &QueueItem{
		ID:        misc.MakeTimeStampedRandomID(100).Hex(),
		Submitted: time.Now().UnixNano(),
	}
}

func (qi *QueueItem) SetPost(req *rpc.ReqNewPost) *QueueItem {
	qi.ReqNewThread = nil
	qi.ReqNewPost = req
	return qi
}

func (qi *QueueItem) SetThread(req *rpc.ReqNewThread) *QueueItem {
	qi.ReqNewPost = nil
	qi.ReqNewThread = req
	return qi
}
