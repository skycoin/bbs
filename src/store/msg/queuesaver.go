package msg

import (
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/cmd/bbsnode/args"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/rpc"
	"github.com/skycoin/bbs/src/store/cxo"
	"github.com/skycoin/bbs/src/store/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// QueueSaverFileName represents the filename of the queue configuration file.
const QueueSaverFileName = "bbs_queue.json"

type QueueSaver struct {
	sync.Mutex
	config   *args.Config
	c        *cxo.Container
	subIndex map[string]int
	queue    []*QueueItem
	done     []*QueueItem
	quit     chan struct{}
}

func NewQueueSaver(config *args.Config, container *cxo.Container) (*QueueSaver, error) {
	qs := QueueSaver{
		config:   config,
		c:        container,
		subIndex: make(map[string]int),
		quit:     make(chan struct{}),
	}
	qs.load()
	if e := qs.save(); e != nil {
		return nil, e
	}
	go qs.serve()
	return &qs, nil
}

func (qs *QueueSaver) absConfigDir() string {
	return filepath.Join(qs.config.ConfigDir(), QueueSaverFileName)
}

func (qs *QueueSaver) load() error {
	// Don't load if specified not to.
	if !qs.config.SaveConfig() {
		return nil
	}
	if e := util.LoadJSON(qs.absConfigDir(), &qs.queue); e != nil {
		return e
	}
	return nil
}

func (qs *QueueSaver) save() error {
	// Don't save if specified.
	if !qs.config.SaveConfig() {
		return nil
	}
	return util.SaveJSON(qs.absConfigDir(), qs.queue, os.FileMode(0700))
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

func (qs *QueueSaver) Ping(address string) error {
	timer := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-timer.C:
			return errors.New("timeout")
		default:
			rpcClient, e := rpc.NewClient(address)
			if e != nil {
				break
			}
			if rpcClient.PingPong(); e != nil {
				break
			}
			return nil
		}
	}
}

func (qs *QueueSaver) Process() {
	qs.Lock()
	defer qs.Unlock()

	if len(qs.queue) == 0 {
		return
	}
	log.Println("[QUEUESAVER] Processing queue ...")
	doneList := []int{}

	connectRemote := func(bpk cipher.PubKey, qi *QueueItem) (*rpc.Client, error) {
		b, e := qs.c.GetBoard(bpk)
		if e != nil {
			log.Printf("[QUEUESAVER] \t- Failed to get board '%s': %s",
				bpk.Hex(), e.Error())
			qi.GetBoardFails += 1
			return nil, e
		}
		subAddr, e := qs.getSubAddr(b)
		if e != nil {
			log.Printf("[QUEUESAVER] \t- Failed to get submission address of board '%s': '%s'",
				b.PubKey, e.Error())
			return nil, e
		}
		rpcClient, e := rpc.NewClient(subAddr)
		if e != nil {
			log.Printf("[QUEUESAVER] \t- Failed to connect to '%s': %s",
				subAddr, e.Error())
			qi.ConnectionFails += 1
			return nil, e
		}
		return rpcClient, nil
	}

	for i, qi := range qs.queue {
		switch {
		case qi.ReqNewPost != nil:
			req := qi.ReqNewPost
			log.Println("[QUEUESAVER] RESENDING: (New Post Request)")

			rpcClient, e := connectRemote(req.BoardPubKey, qi)
			if e != nil {
				break
			}

			ref, e := rpcClient.NewPost(req)
			log.Printf("[QUEUESAVER] \t- (COMPLETE!) Ref: '%s', Err: '%s'", ref, e.Error())
			doneList = append(doneList, i)

		case qi.ReqNewThread != nil:
			req := qi.ReqNewThread
			log.Println("[QUEUESAVER] (New Thread Request)")

			rpcClient, e := connectRemote(req.BoardPubKey, qi)
			if e != nil {
				break
			}

			ref, e := rpcClient.NewThread(req)
			log.Printf("[QUEUESAVER] \t- (COMPLETE!) Ref: '%s', Err: '%s'", ref, e.Error())
			doneList = append(doneList, i)

		case qi.ReqVotePost != nil:
			req := qi.ReqVotePost
			log.Println("[QUEUESAVER] RESENDING: (Vote Post Request)")

			rpcClient, e := connectRemote(req.BoardPubKey, qi)
			if e != nil {
				break
			}

			ok, e := rpcClient.VotePost(req)
			log.Printf("[QUEUESAVER] \t- (COMPLETE!) Status: '%v', Err: '%s'", ok, e.Error())
			doneList = append(doneList, i)

		case qi.ReqVoteThread != nil:
			req := qi.ReqVoteThread
			log.Println("[QUEUESAVER] RESENDING: (Vote Thread Request)")

			rpcClient, e := connectRemote(req.BoardPubKey, qi)
			if e != nil {
				break
			}

			ok, e := rpcClient.VoteThread(req)
			log.Printf("[QUEUESAVER] \t- (COMPLETE!) Status: '%v', Err: '%s'", ok, e.Error())
			doneList = append(doneList, i)
		}
	}

	for _, i := range doneList {
		qs.done = append(qs.done, qs.queue[i])
	}
	for i := range misc.ReverseIntSlice(doneList) {
		qs.queue = append(qs.queue[:i], qs.queue[i+1:]...)
	}
	if len(doneList) > 0 {
		qs.save()
	}
}

func (qs *QueueSaver) AddNewPostReq(bpk cipher.PubKey, tRef skyobject.Reference, post *typ.Post) error {
	qs.Lock()
	defer qs.Unlock()
	log.Println("[QUEUESAVER] Sending new post request...")
	req := &rpc.ReqNewPost{bpk, tRef, post}
	b, e := qs.c.GetBoard(bpk)
	if e != nil {
		e = errors.Wrap(e, "failed to obtain board")
		log.Printf("[QUEUESAVER] \t- Error: '%s'", e.Error())
		return e
	}
	subAddr, e := qs.getSubAddr(b)
	if e != nil {
		e = errors.Wrap(e, "failed to obtain submission address")
		log.Printf("[QUEUESAVER] \t- Error: '%s'", e.Error())
		return e
	}
	rpcClient, e := rpc.NewClient(subAddr)
	if e != nil {
		// Add to queue.
		log.Printf("[QUEUESAVER] \t- rpc error: '%s'", e.Error())
		log.Printf("[QUEUESAVER] \t- adding request to queue: %v", req)
		qs.queue = append(qs.queue, NewQueueItem().SetReqNewPost(req))
		qs.save()
		return nil
	}
	if pRef, e := rpcClient.NewPost(req); e != nil {
		e = errors.Wrap(e, "error from remote")
		log.Println("[QUEUESAVER] \t- reply:", e)
		return e
	} else {
		log.Println("[QUEUESAVER] \t- (COMPLETE!) Ref:", pRef)
		post.Ref = pRef
		return nil
	}
}

func (qs *QueueSaver) AddNewThreadReq(bpk, upk cipher.PubKey, usk cipher.SecKey, thread *typ.Thread) error {
	qs.Lock()
	defer qs.Unlock()
	log.Println("[QUEUESAVER] Sending new thread request...")
	req := &rpc.ReqNewThread{bpk, upk, thread.Sign(usk), thread}
	b, e := qs.c.GetBoard(bpk)
	if e != nil {
		e = errors.Wrap(e, "failed to obtain board")
		log.Printf("[QUEUESAVER] \t- Error: '%s'", e.Error())
		return e
	}
	subAddr, e := qs.getSubAddr(b)
	if e != nil {
		e = errors.Wrap(e, "failed to obtain submission address")
		log.Printf("[QUEUESAVER] \t- Error: '%s'", e.Error())
		return e
	}
	rpcClient, e := rpc.NewClient(subAddr)
	if e != nil {
		// Add to queue.
		log.Printf("[QUEUESAVER] \t- rpc error: '%s'", e.Error())
		log.Printf("[QUEUESAVER] \t- adding request to queue: %v", req)
		qs.queue = append(qs.queue, NewQueueItem().SetReqNewThread(req))
		qs.save()
		return e
	}
	if tRef, e := rpcClient.NewThread(req); e != nil {
		e = errors.Wrap(e, "error from remote")
		log.Println("[QUEUESAVER] \t- reply:", e)
		return e
	} else {
		log.Println("[QUEUESAVER] \t- (COMPLETE!) Ref:", tRef)
		thread.Ref = tRef
		return nil
	}
}

func (qs *QueueSaver) AddVotePostReq(bpk cipher.PubKey, pRef skyobject.Reference, vote *typ.Vote) error {
	qs.Lock()
	defer qs.Unlock()
	log.Println("[QUEUESAVER] Sending vote post request...")
	req := &rpc.ReqVotePost{bpk, pRef, vote}
	b, e := qs.c.GetBoard(bpk)
	if e != nil {
		e = errors.Wrap(e, "failed to obtain board")
		log.Printf("[QUEUESAVER] \t- Error: '%s'", e.Error())
		return e
	}
	subAddr, e := qs.getSubAddr(b)
	if e != nil {
		e = errors.Wrap(e, "failed to obtain submission address")
		log.Printf("[QUEUESAVER] \t- Error: '%s'", e.Error())
		return e
	}
	rpcClient, e := rpc.NewClient(subAddr)
	if e != nil {
		// Add to queue.
		log.Printf("[QUEUESAVER] \t- rpc error: '%s'", e.Error())
		log.Printf("[QUEUESAVER] \t- adding request to queue: %v", req)
		qs.queue = append(qs.queue, NewQueueItem().SetReqVotePost(req))
		qs.save()
		return nil
	}
	if ok, e := rpcClient.VotePost(req); e != nil {
		e = errors.Wrap(e, "error from remote")
		log.Println("[QUEUESAVER] \t- reply:", e)
		return e
	} else {
		log.Println("[QUEUESAVER] \t- (COMPLETE!) Status:", ok)
		return nil
	}
}

func (qs *QueueSaver) AddVoteThreadReq(bpk cipher.PubKey, tRef skyobject.Reference, vote *typ.Vote) error {
	qs.Lock()
	defer qs.Unlock()
	log.Println("[QUEUESAVER] Sending vote thread request...")
	req := &rpc.ReqVoteThread{bpk, tRef, vote}
	b, e := qs.c.GetBoard(bpk)
	if e != nil {
		e = errors.Wrap(e, "failed to obtain board")
		log.Printf("[QUEUESAVER] \t- Error: '%s'", e.Error())
		return e
	}
	subAddr, e := qs.getSubAddr(b)
	if e != nil {
		e = errors.Wrap(e, "failed to obtain submission address")
		log.Printf("[QUEUESAVER] \t- Error: '%s'", e.Error())
		return e
	}
	rpcClient, e := rpc.NewClient(subAddr)
	if e != nil {
		// Add to queue.
		log.Printf("[QUEUESAVER] \t- rpc error: '%s'", e.Error())
		log.Printf("[QUEUESAVER] \t- adding request to queue: %v", req)
		qs.queue = append(qs.queue, NewQueueItem().SetReqVoteThread(req))
		qs.save()
		return nil
	}
	if ok, e := rpcClient.VoteThread(req); e != nil {
		e = errors.Wrap(e, "error from remote")
		log.Println("[QUEUESAVER] \t- reply:", e)
		return e
	} else {
		log.Println("[QUEUESAVER] \t- (COMPLETE!) Status:", ok)
		return nil
	}
}

func (qs *QueueSaver) getSubAddr(board *typ.Board) (string, error) {
	meta, e := board.GetMeta()
	if e != nil {
		return "", e
	}
	addrs := meta.SubmissionAddresses

	i, has := qs.subIndex[board.PubKey]
	if has == false {
		i, _ = misc.MakeIntBetween(0, len(addrs)-1)
		qs.subIndex[board.PubKey] = i
	}
	if i += 1; i >= len(addrs) {
		i = 0
	}
	return addrs[i], nil
}
