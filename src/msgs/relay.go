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
	"github.com/skycoin/bbs/src/misc/keys"
)

const (
	ReceiverPrefix = "MSGRECEIVER"
)

type MsgType byte

const (
	MsgNewThread MsgType = iota << 0
	MsgNewPost
	MsgNewVote
	MsgResponse
)

var MsgTypeStr = [...]string{
	MsgNewThread: "New Thread",
	MsgNewPost: "New Post",
	MsgNewVote: "New Vote",
	MsgResponse: "Response",
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
	in         chan *Message
	quit       chan struct{}
	wg         sync.WaitGroup
}

func NewRelay(config *RelayConfig) *Relay {
	return &Relay{
		c:          config,
		l:          inform.NewLogger(true, os.Stdout, ReceiverPrefix),
		factory:    factory.NewMessengerFactory(),
		incomplete: NewIncomplete(),
		in:         make(chan *Message),
		quit:       make(chan struct{}),
	}
}

func (r *Relay) Open(compiler *state.Compiler) error {
	r.compiler = compiler
	if e := r.setup(); e != nil {
		r.l.Panicln("failed to setup 'Relay':", e)
	}
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
			OnConnected: func(conn *factory.Connection) {
				go func(wg *sync.WaitGroup, quit chan struct{}) {
					wg.Add(1)
					defer wg.Done()

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
							if msg, e := NewMessage(data); e != nil {
								r.l.Printf("'%s' skipping invalid message: %v",
									conn.GetKey().Hex()[:5]+"...", e)
							} else {
								r.in <- msg
							}
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

	for {
		select {
		case <-r.quit:
			return

		case in := <-r.in:
			if e := r.receiveMessage(in); e != nil {
				r.l.Println(e)
			}
		}
	}
}

func (r *Relay) send(toPK cipher.PubKey, t MsgType, d1 []byte) error {
	sent := false
	errors := []error{}

	d0 := []byte{byte(t)}
	//d1 := encoder.Serialize(v)

	r.factory.ForEachConn(func(conn *factory.Connection) {
		if sent {
			return
		}
		out := append(append(d0, keys.PubKeyToSlice(conn.GetKey())...), d1...)
		if e := conn.Send(toPK, out); e != nil {
			r.l.Printf("'%s' send error: %v",
				conn.GetKey().Hex()[:5]+"...", e)
			errors = append(errors, e)
		} else {
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

func (r *Relay) receiveMessage(msg *Message) error {

	switch msg.GetMsgType() {
	case MsgNewThread:
		goal, e := r.processNewThread(msg)
		return r.send(msg.GetOutPK(), MsgResponse,
			encoder.Serialize(NewResponse(msg, goal, e)))

	case MsgNewPost:
		goal, e := r.processNewPost(msg)
		return r.send(msg.GetOutPK(), MsgResponse,
			encoder.Serialize(NewResponse(msg, goal, e)))

	case MsgNewVote:
		goal, e := r.processNewVote(msg)
		return r.send(msg.GetOutPK(), MsgResponse,
			encoder.Serialize(NewResponse(msg, goal, e)))

	case MsgResponse:
		return r.processResponse(msg)

	default:
		return boo.Newf(boo.NotAllowed, "unknown message type '%v'",
			msg.GetMsgType())
	}
}

func (r *Relay) processNewThread(msg *Message) (uint64, error) {
	thread, e := msg.ExtractThread()
	if e != nil {
		return 0, e
	}
	bi, e := r.compiler.GetBoard(thread.OfBoard)
	if e != nil {
		return 0, e
	}
	if !bi.IsMaster() {
		return 0, notMasterErr(thread.OfBoard)
	}
	return bi.NewThread(thread)
}

func (r *Relay) processNewPost(msg *Message) (uint64, error) {
	post, e := msg.ExtractPost()
	if e != nil {
		return 0, e
	}
	bi, e := r.compiler.GetBoard(post.OfBoard)
	if e != nil {
		return 0, e
	}
	if !bi.IsMaster() {
		return 0, notMasterErr(post.OfBoard)
	}
	return bi.NewPost(post)
}

func (r *Relay) processNewVote(msg *Message) (uint64, error) {
	vote, e := msg.ExtractVote()
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

func (r *Relay) processResponse(msg *Message) error {
	res, e := msg.ExtractResponse()
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

func (r *Relay) multiSendRequest(ctx context.Context, toPKs []cipher.PubKey, mt MsgType, data []byte) (uint64, error) {
	if len(toPKs) == 0 {
		return 0, boo.New(boo.InvalidInput, "no submission public keys provided")
	}
	r.l.Printf("Sending '%s' request ...", MsgTypeStr[mt])
	var goal uint64
	var e error
	for i, pk := range toPKs {
		ctxReq, _ :=  context.WithTimeout(ctx, time.Second * 10)
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
				return 0, boo.New(res.ErrTyp, res.ErrMsg)
			}
		}
	}
}

func notMasterErr(bpk cipher.PubKey) error {
	return boo.Newf(boo.NotAllowed,
		"node is not owner of board of public key '%s'", bpk.Hex())
}
