package store

import (
	"github.com/evanlinjin/bbs/cmd"
	"github.com/evanlinjin/bbs/typ"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util"
)

// QueueConfigFileName represents the filename of the queue configuration file.
const QueueConfigFileName = "bbs_queue.json"

type ReqNewPost struct {
	BoardPubKey cipher.PubKey       `json:"board_public_key,string"`
	ThreadRef   skyobject.Reference `json:"thread_reference,string"`
	Post        *typ.Post           `json:"post"`
}

type ReqNewThread struct {
	BoardPubKey cipher.PubKey `json:"board_public_key,string"`
	Creator     cipher.PubKey `json:"creator,string"`
	Signature   cipher.Sig    `json:"signature,string"`
	Thread      *typ.Thread   `json:"thread"`
}

type QueueItem struct {
	Submitted    int64         `json:"submitted"`
	ReqNewPost   *ReqNewPost   `json:"new_post_request,omitempty"`
	ReqNewThread *ReqNewThread `json:"new_thread_request,omitempty"`
}

type QueueSaver struct {
	config *cmd.Config
	c      *Container
	q      []*QueueItem
	quit   chan struct{}
}

func NewQueueSaver(config *cmd.Config, container *Container) (*QueueSaver, error) {
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
	if e := util.LoadJSON(QueueConfigFileName, &qs.q); e != nil {
		return e
	}
	return nil
}

func (qs *QueueSaver) save() error {
	return util.SaveJSON(QueueConfigFileName, &qs.q, 0600)
}

func (qs *QueueSaver) serve() {
	for {
		select {
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
