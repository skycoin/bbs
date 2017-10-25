package accord

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
	"os"
	"sync"
)

const (
	LogPrefix = "ACCORD_RELAY"
)

type Relay struct {
	l              *log.Logger
	factory        *factory.MessengerFactory
	compiler       *state.Compiler
	incomplete     *Incomplete
	disconnectChan chan string
	quit           chan struct{}
	wg             sync.WaitGroup
}

func NewRelay() *Relay {
	return &Relay{
		l:              inform.NewLogger(true, os.Stdout, LogPrefix),
		factory:        factory.NewMessengerFactory(),
		incomplete:     NewIncomplete(),
		disconnectChan: make(chan string),
		quit:           make(chan struct{}),
	}
}

func (r *Relay) Open(compiler *state.Compiler) error {
	r.compiler = compiler
	return nil
}

func (r *Relay) Close() {
	close(r.quit)
	r.wg.Wait()
}

func (r *Relay) Connect(address string) (cipher.PubKey, error) {
	conn, e := r.factory.ConnectWithConfig(address, &factory.ConnConfig{
		Reconnect:   false,
		OnConnected: r.connectionService,
	})
	if e != nil {
		return cipher.PubKey{}, e
	} else {
		return conn.GetKey(), nil
	}
}

func (r *Relay) Disconnect(key cipher.PubKey) bool {
	conn, ok := r.factory.GetConnection(key)
	if !ok {
		return false
	} else {
		conn.Close()
		return true
	}
}

func (r *Relay) connectionService(conn *factory.Connection) {
	r.l.Printf("service started for connection: address(%s) public_key(%s)",
		conn.GetRemoteAddr().String(), conn.GetKey().Hex())

	var (
		address = conn.GetRemoteAddr().String()
		pk      = conn.GetKey()
	)

	go func(r *Relay, conn *factory.Connection) {
		r.wg.Add(1)
		defer r.wg.Done()

		for {
			select {
			case <-r.quit:
				return

			case data, ok := <-conn.GetChanIn():

				if !ok {
					r.disconnectChan <- address
					r.l.Printf("(%s:%s) disconnected",
						address, pk.Hex()[:5]+"...")
					return
				}

				wrapper, e := NewWrapper(data)
				if e != nil {
					r.l.Printf("(%s:%s) received invalid message, error: %v",
						address, pk.Hex()[:5]+"...", e)
					continue
				} else {
					r.l.Printf("(%s:%s) received message of type '%s'",
						address, pk.Hex()[:5]+"...", wrapper.GetType().String())
				}

				switch wrapper.GetType() {
				case SubmissionType:
					e := send(conn, wrapper.GetFromPK(), SubmissionResponseType,
						NewSubmissionResponse(r.processSubmission(conn, address, pk, wrapper)).Serialize())
					if e != nil {
						r.l.Println("failed to send message, error:", e)
					}

				case SubmissionResponseType:
					if res, e := wrapper.ToSubmissionResponse(); e != nil {
						r.l.Println("failed to obtain submission response, error:", e)
					} else {
						r.incomplete.Satisfy(res)
					}
				}
			}
		}
	}(r, conn)

	r.l.Printf("(%s:%s) connected", address, pk.Hex()[:5]+"...")
}

func (r *Relay) processSubmission(conn *factory.Connection, address string, pk cipher.PubKey, wrapper *Wrapper) (
	hash cipher.SHA256, goal uint64, e error,
) {
	submission, e := wrapper.ToSubmission()
	if e != nil {
		e = boo.WrapType(e, boo.InvalidRead, "failed to extract submission")
		return
	}
	transport, e := submission.ToTransport()
	if e != nil {
		e = boo.WrapType(e, boo.InvalidRead, "failed to extract transport")
		return
	} else {
		hash = transport.Header.GetHash()
	}
	bi, e := r.compiler.GetBoard(transport.GetOfBoard())
	if e != nil {
		e = boo.WrapType(e, boo.InvalidRead, "failed to obtain board instance")
		return
	}
	if bi.IsMaster() == false {
		e = boo.WrapType(e, boo.NotAllowed, "node does not own this board")
		return
	}
	if goal, e = bi.Submit(transport); e != nil {
		e = boo.WrapType(e, boo.Type(e), "submission failed")
		return
	}
	return
}

func send(conn *factory.Connection, toPK cipher.PubKey, t Type, body []byte) error {
	return conn.Send(toPK, append([]byte{byte(t)}, body...))
}

func (r *Relay) SubmitToRemote(ctx context.Context, fromPK, toPK cipher.PubKey, data interface{}) (uint64, error) {
	switch t := data.(type) {
	case *Submission:
		conn, ok := r.factory.GetConnection(fromPK)
		if !ok {
			return 0, boo.Newf(boo.NotFound, "failed to find connection of public key %s", fromPK.Hex())
		}

		hash := data.(*Submission).GetHash()
		resChan, e := r.incomplete.Add(hash)
		if e != nil {
			return 0, e
		}
		defer r.incomplete.Remove(hash)

		if e := send(conn, toPK, SubmissionType, encoder.Serialize(data)); e != nil {
			return 0, e
		}
		for {
			select {
			case <-ctx.Done():
				return 0, ctx.Err()
			case res := <-resChan:
				return res.Seq, res.Error()
			}
		}
	default:
		return 0, boo.Newf(boo.InvalidInput, "invalid type '%T'", t)
	}
}

func (r *Relay) SubmissionKeys() []*object.MessengerSubKeyTransport {
	var out []*object.MessengerSubKeyTransport
	r.factory.ForEachConn(func(conn *factory.Connection) {
		out = append(out, &object.MessengerSubKeyTransport{
			Address: conn.GetRemoteAddr().String(),
			PubKey:  conn.GetKey(),
		})
	})
	r.l.Println("got submission keys")
	return out
}

func (r *Relay) Disconnections() chan string {
	return r.disconnectChan
}
