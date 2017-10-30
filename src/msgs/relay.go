package msgs

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/object/revisions/r0"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
	"os"
	"sync"
	"time"
)

const (
	ReceiverPrefix = "MSGRECEIVER"
	ServiceName    = "skycoin_bbs"
)

type MsgType byte

const (
	// New Content.
	MsgNewThread MsgType = iota << 0
	MsgNewPost
	MsgNewVote
	MsgNewContentResponse

	// Ask Boards.
	MsgDiscoverer
)

var MsgTypeStr = [...]string{
	MsgNewThread:          "New Thread",
	MsgNewPost:            "New Post",
	MsgNewVote:            "New Vote",
	MsgNewContentResponse: "New Content Response",
	MsgDiscoverer:         "Boards Discoverer",
}

func (mt MsgType) String() string {
	if int(mt) >= len(MsgTypeStr) {
		return "Unknown"
	}
	return MsgTypeStr[mt]
}

type RelayConfig struct {
	Addresses         []string
	ReconnectInterval *time.Duration
}

type Relay struct {
	c          *RelayConfig
	l          *log.Logger
	factory    *factory.MessengerFactory
	compiler   *state.Compiler
	incomplete *Incomplete
	//discoverer  *BoardDiscoverer
	in          chan *BBSMessage
	onConnected chan struct{}
	quit        chan struct{}
	wg          sync.WaitGroup
}

func NewRelay(config *RelayConfig) *Relay {
	return &Relay{
		c:           config,
		l:           inform.NewLogger(true, os.Stdout, ReceiverPrefix),
		factory:     factory.NewMessengerFactory(),
		incomplete:  NewIncomplete(),
		in:          make(chan *BBSMessage),
		onConnected: make(chan struct{}),
		quit:        make(chan struct{}),
	}
}

func (r *Relay) Open(compiler *state.Compiler) error {
	r.compiler = compiler
	if e := r.setup(); e != nil {
		r.l.Panicln("failed to setup 'Relay':", e)
	}
	//r.discoverer = NewBoardDiscoverer(r.send, r.compiler.RangeMasterSubs)
	go r.service()
	return nil
}

func (r *Relay) Close() {
	for {
		select {
		case r.quit <- struct{}{}:
		default:
			r.factory.Close()
			r.wg.Wait()
			return
		}
	}
}

func (r *Relay) GetKeys() []cipher.PubKey {
	var keys []cipher.PubKey
	r.factory.ForEachConn(func(conn *factory.Connection) {
		keys = append(keys, conn.GetKey())
	})
	return keys
}

func (r *Relay) setup() error {
	if len(r.c.Addresses) == 0 {
		return boo.New(boo.InvalidInput,
			"no messenger addresses provided")
	}

	for _, address := range r.c.Addresses {
		r.factory.ConnectWithConfig(address, &factory.ConnConfig{
			Reconnect:     true,
			ReconnectWait: *r.c.ReconnectInterval,
			FindServiceNodesByAttributesCallback: func(resp *factory.QueryByAttrsResp) {
				for i, pk := range resp.Result[ServiceName] {
					r.l.Printf(" - [%d] found service node '%s'.", i, pk.Hex()[:5]+"...")
					//r.discoverer.AddNode(pk)
				}
			},
			OnConnected: func(conn *factory.Connection) {
				r.l.Println("Connected!", conn.GetKey().Hex()[:5]+"...")

				if e := conn.OfferService(ServiceName); e != nil {
					r.l.Printf("failed to offer service '%s'", ServiceName)
				}

				// TODO: Create a waiting function that wraps this.
				if e := conn.FindServiceNodesByAttributes(ServiceName); e != nil {
					r.l.Printf("failed to find services of '%s'", ServiceName)
				}

				go func(wg *sync.WaitGroup, quit chan struct{}) {
					wg.Add(1)
					defer wg.Done()

					r.onConnected <- struct{}{}

					for {
						select {
						case <-quit:
							return

						case data, ok := <-conn.GetChanIn():

							if !ok {
								r.l.Printf("'%s' disconnected.",
									conn.GetKey().Hex()[:5]+"...")
								return
							}

							msg := Message(data)

							if e := msg.Check(); e != nil {
								r.l.Printf("'%s' skipping invalid message: %v",
									conn.GetKey().Hex()[:5]+"...", e)
								continue
							}

							sendMsg, e := msg.ToSendMessage()

							if e != nil {
								r.l.Printf("'%s' skipping invalid message: %v",
									conn.GetKey().Hex()[:5]+"...", e)
								continue
							}

							bbsMsg, e := sendMsg.ToBBSMessage()

							if e != nil {
								r.l.Printf("'%s' skipping invalid message: %v",
									conn.GetKey().Hex()[:5]+"...", e)
								continue
							}

							r.in <- bbsMsg
						}
					}
				}(&r.wg, r.quit)
			},
		})
	}
	return nil
}

func (r *Relay) service() {
	r.wg.Add(1)
	defer r.wg.Done()

	//nodeInterval := time.Minute
	//nodeRefreshTicker := time.NewTicker(nodeInterval)
	//defer nodeRefreshTicker.Stop()

	for {
		select {
		case <-r.quit:
			return

		case in := <-r.in:
			if e := r.receiveMessage(in); e != nil {
				r.l.Println(e)
			}

		//case <-nodeRefreshTicker.C:
		//r.discoverer.ClearNodes()
		//r.factory.ForEachConn(func(conn *factory.Connection) {
		//	if e := conn.FindServiceNodesByAttributes(ServiceName); e != nil {
		//		r.l.Printf("failed to find services of '%s'", ServiceName)
		//	}
		//})
		//go func() {
		//	time.Sleep(nodeInterval/2)
		//r.discoverer.SendAsk()
		//}()

		case <-r.onConnected:
			if e := r.compiler.EnsureSubmissionKeys(r.GetKeys()); e != nil {
				r.l.Println(e)
			}
		}
	}
}

func (r *Relay) send(toPK cipher.PubKey, t MsgType, body []byte) error {
	sent := false
	errors := []error{}

	out := append([]byte{byte(t)}, body...)

	r.factory.ForEachConn(func(conn *factory.Connection) {
		if sent {
			return
		}
		if e := conn.Send(toPK, out); e != nil {
			r.l.Printf("'%s' send error: %v",
				conn.GetKey().Hex()[:5]+"...", e)
			errors = append(errors, e)
		} else {
			r.l.Printf("'%s' sent message: body_len(%d)",
				conn.GetKey().Hex()[:5]+"...", len(body))
			sent = true
		}
	})

	if !sent {
		if len(errors) == 0 {
			return boo.New(boo.NotFound,
				"not connected")
		}
		return boo.WrapType(errors[0], boo.NotFound,
			"failed to send message")
	}
	return nil
}

func (r *Relay) receiveMessage(msg *BBSMessage) error {
	r.l.Printf("message received: type(%s) body_len(%d)",
		msg.GetMsgType().String(), len(msg.GetBody()))

	switch msg.GetMsgType() {
	case MsgNewThread:
		goal, e := r.processNewThread(msg)
		return r.send(msg.GetFromPubKey(), MsgNewContentResponse,
			encoder.Serialize(GenerateNewContentResponse(msg, goal, e)))

	case MsgNewPost:
		goal, e := r.processNewPost(msg)
		return r.send(msg.GetFromPubKey(), MsgNewContentResponse,
			encoder.Serialize(GenerateNewContentResponse(msg, goal, e)))

	case MsgNewVote:
		goal, e := r.processNewVote(msg)
		return r.send(msg.GetFromPubKey(), MsgNewContentResponse,
			encoder.Serialize(GenerateNewContentResponse(msg, goal, e)))

	case MsgNewContentResponse:
		return r.processResponse(msg)

	case MsgDiscoverer:
		if _, e := msg.ExtractDiscovererMsg(); e != nil {
			return e
		} else {
			//r.discoverer.Process(msg.GetFromPubKey(), out)
			return nil
		}

	default:
		return boo.Newf(boo.NotAllowed, "unknown message type '%v'",
			msg.GetMsgType())
	}
}

func (r *Relay) processNewThread(msg *BBSMessage) (uint64, error) {
	thread, e := msg.ExtractContentThread()
	if e != nil {
		return 0, e
	}
	tOfBoard := thread.GetData().GetOfBoard()
	bi, e := r.compiler.GetBoard(tOfBoard)
	if e != nil {
		return 0, e
	}
	if !bi.IsMaster() {
		return 0, notMasterErr(tOfBoard)
	}
	return bi.NewThread(thread)
}

func (r *Relay) processNewPost(msg *BBSMessage) (uint64, error) {
	post, e := msg.ExtractContentPost()
	if e != nil {
		return 0, e
	}
	pOfBoard := post.GetData().GetOfBoard()
	bi, e := r.compiler.GetBoard(pOfBoard)
	if e != nil {
		return 0, e
	}
	if !bi.IsMaster() {
		return 0, notMasterErr(pOfBoard)
	}
	return bi.NewPost(post)
}

func (r *Relay) processNewVote(msg *BBSMessage) (uint64, error) {
	vote, e := msg.ExtractContentVote()
	if e != nil {
		return 0, e
	}
	bi, e := r.compiler.GetBoard(vote.OfBoard)
	if e != nil {
		return 0, e
	}
	if !bi.IsMaster() {
		return 0, notMasterErr(vote.OfBoard)
	}
	return bi.NewVote(vote)
}

func (r *Relay) processResponse(msg *BBSMessage) error {
	res, e := msg.ExtractNewContentResponse()
	if e != nil {
		return e
	}
	r.incomplete.Satisfy(res)
	return nil
}

func (r *Relay) NewThread(ctx context.Context, toPKs []cipher.PubKey, thread *r0.Thread) (uint64, error) {
	return r.multiSendRequest(ctx, toPKs, MsgNewThread, encoder.Serialize(thread))
}

func (r *Relay) NewPost(ctx context.Context, toPKs []cipher.PubKey, post *r0.Post) (uint64, error) {
	return r.multiSendRequest(ctx, toPKs, MsgNewPost, encoder.Serialize(post))
}

func (r *Relay) NewVote(ctx context.Context, toPKs []cipher.PubKey, vote *r0.Vote) (uint64, error) {
	return r.multiSendRequest(ctx, toPKs, MsgNewVote, encoder.Serialize(vote))
}

func (r *Relay) GetBoards() []string {
	return []string{"not implemented"}
	//return r.discoverer.GetBoards()
}

func (r *Relay) multiSendRequest(ctx context.Context, toPKs []cipher.PubKey, mt MsgType, data []byte) (uint64, error) {
	if len(toPKs) == 0 {
		return 0, boo.New(boo.InvalidInput, "no submission public keys provided")
	}
	r.l.Printf("Sending '%s' request ...", MsgTypeStr[mt])
	var goal uint64
	var e error
	for i, pk := range toPKs {
		ctxReq, _ := context.WithTimeout(ctx, time.Second*30)
		if goal, e = r.sendRequest(ctxReq, pk, mt, data); e != nil {
			r.l.Printf(" - [%d] send request to '%s' failed with error: %v",
				i, pk.Hex()[:5]+"...", e)
			continue
		} else {
			r.l.Printf(" - [%d] send request to '%s' succeeded to seq: %v",
				i, pk.Hex()[:5]+"...", goal)
			return goal, nil
		}
	}
	return 0, e
}

func (r *Relay) sendRequest(ctx context.Context, toPK cipher.PubKey, mt MsgType, data []byte) (uint64, error) {
	hash := cipher.SumSHA256(data)

	resChan, e := r.incomplete.Add(hash)
	if e != nil {
		return 0, e
	}
	defer r.incomplete.Remove(hash)

	if e := r.send(toPK, mt, data); e != nil {
		return 0, e
	}
	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case res := <-resChan:
			if res.Okay {
				return res.Seq, nil
			} else {
				return 0, boo.New(int(res.ErrTyp), res.ErrMsg)
			}
		}
	}
}

func notMasterErr(bpk cipher.PubKey) error {
	return boo.Newf(boo.NotAllowed,
		"node is not owner of board of public key '%s'", bpk.Hex())
}
