package msg

import (
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/extern/rpc"
	"github.com/evanlinjin/bbs/intern/cxo"
	"github.com/evanlinjin/bbs/intern/typ"
	"github.com/evanlinjin/bbs/misc"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"sync"
	"time"
)

// QueueConfigFileName represents the filename of the queue configuration file.
const QueueConfigFileName = "bbs_queue.json"

type QueueItem struct {
	ID           string            `json:"id"`
	Submitted    int64             `json:"submitted"`
	ReqNewPost   *rpc.ReqNewPost   `json:"new_post_request,omitempty"`
	ReqNewThread *rpc.ReqNewThread `json:"new_thread_request,omitempty"`
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

type QueueSaver struct {
	sync.Mutex
	config *cmd.Config
	c      *cxo.Container
	queue  []*QueueItem
	done   []*QueueItem
	quit   chan struct{}
}

func NewQueueSaver(config *cmd.Config, container *cxo.Container) (*QueueSaver, error) {
	qs := QueueSaver{
		config: config,
		c:      container,
		quit:   make(chan struct{}),
	}
	qs.load()
	if e := qs.save(); e != nil {
		return nil, e
	}
	go qs.serve()
	return &qs, nil
}

func (qs *QueueSaver) load() error {
	if e := util.LoadJSON(QueueConfigFileName, &qs.queue); e != nil {
		return e
	}
	return nil
}

func (qs *QueueSaver) save() error {
	return util.SaveJSON(QueueConfigFileName, &qs.queue, 0600)
}

func (qs *QueueSaver) serve() {
	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-ticker.C:
			qs.Process()

		case <-qs.quit:
			return
		}
	}
}

func (qs *QueueSaver) Close() {
	select {
	case qs.quit <- struct{}{}:
	default:
	}
}

func (qs *QueueSaver) Process() {
	//TODO: screams for refactoring
	qs.Lock()
	defer qs.Unlock()
	if len(qs.queue) == 0 {
		return
	}
	log.Println("[QUEUESAVER] Processing queue...")

	toRemoveList := []int{}

	for i, qi := range qs.queue {
		switch {
		case qi.ReqNewPost != nil:
			req := qi.ReqNewPost
			b, e := qs.c.GetBoard(req.BoardPubKey)
			if e != nil {
				toRemoveList = append(toRemoveList, i)
				break
			}
			rpcClient, e := rpc.NewClient(b.URL)
			if e != nil {
				break
			}
			rpcClient.NewPost(req)
			toRemoveList = append(toRemoveList, i)

		case qi.ReqNewThread != nil:
			req := qi.ReqNewThread
			b, e := qs.c.GetBoard(req.BoardPubKey)
			if e != nil {
				toRemoveList = append(toRemoveList, i)
				break
			}
			rpcClient, e := rpc.NewClient(b.URL)
			if e != nil {
				break
			}
			rpcClient.NewThread(req)
			toRemoveList = append(toRemoveList, i)
		}
	}
	for _, i := range toRemoveList {
		qs.done = append(qs.done, qs.queue[i])
	}
	for i := range misc.ReverseIntSlice(toRemoveList) {
		qs.queue = append(qs.queue[:i], qs.queue[i+1:]...)
	}
}

func (qs *QueueSaver) AddNewPostReq(bpk cipher.PubKey, tRef skyobject.Reference, post *typ.Post) error {
	qs.Lock()
	defer qs.Unlock()
	log.Println("[QUEUESAVER] New Post Request.")
	req := &rpc.ReqNewPost{bpk, tRef, post}
	b, e := qs.c.GetBoard(bpk)
	if e != nil {
		return e
	}
	rpcClient, e := rpc.NewClient(b.URL)
	if e != nil {
		// Add to queue.
		qs.queue = append(qs.queue, NewQueueItem().SetPost(req))
		qs.save()
		return nil
	}
	_, e = rpcClient.NewPost(req)
	return e
}

func (qs *QueueSaver) AddNewThreadReq(bpk, upk cipher.PubKey, usk cipher.SecKey, thread *typ.Thread) error {
	qs.Lock()
	defer qs.Unlock()
	log.Println("[QUEUESAVER] New Thread Request.")
	req := &rpc.ReqNewThread{bpk, upk, thread.Sign(usk), thread}
	b, e := qs.c.GetBoard(bpk)
	if e != nil {
		log.Println("[QUEUESAVER]", e)
		return e
	}
	log.Println("[QUEUESAVER] Got Board.")
	rpcClient, e := rpc.NewClient(b.URL)
	if e != nil {
		// Add to queue.
		log.Println("[QUEUESAVER]", e)
		qs.queue = append(qs.queue, NewQueueItem().SetThread(req))
		qs.save()
		return e
	}
	_, e = rpcClient.NewThread(req)
	if e != nil {
		log.Println("[QUEUESAVER]", e)
		return e
	}
	return e
}
