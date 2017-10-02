package msgs

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"log"
	"os"
	"sync"
	"time"
)

// pkMap stores public keys in order and avoids duplicates.
// It is also thread-safe.
type pkMap struct {
	mux  sync.RWMutex
	dict map[cipher.PubKey]bool
	list []cipher.PubKey
}

// newPKMap creates a new pkMap.
func newPKMap() *pkMap {
	return &pkMap{
		dict: make(map[cipher.PubKey]bool),
	}
}

// Add appends a public key to the list, avoiding duplicates.
func (m *pkMap) Add(key cipher.PubKey) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if !m.dict[key] {
		m.dict[key] = true
		m.list = append(m.list, key)
	}
}

// Del removes a public key from the list (if it exists).
func (m *pkMap) Del(key cipher.PubKey) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.dict[key] {
		delete(m.dict, key)
		for i, v := range m.list {
			if v == key {
				m.list = append(m.list[:i],
					m.list[i+1:]...)
				return
			}
		}
	}
}

// Get returns a list of public keys, given a start index and count.
func (m *pkMap) Get(start, count int) []cipher.PubKey {
	m.mux.RLock()
	defer m.mux.RUnlock()

	l := len(m.list)
	if start >= l {
		return nil
	} else if diff := l - start - count; diff < 0 {
		count += diff
	}

	out := make([]cipher.PubKey, count)
	for i := 0; i < count; i++ {
		out[i] = m.list[start+i]
	}
	return out
}

// Range ranges the public keys from start to end and applies an action to each.
func (m *pkMap) Range(action func(key cipher.PubKey)) {
	m.mux.RLock()
	defer m.mux.RUnlock()
	for _, key := range m.list {
		action(key)
	}
}

// Clear empties the list.
func (m *pkMap) Clear() {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.dict = make(map[cipher.PubKey]bool)
	m.list = []cipher.PubKey{}
}

// SendFunc represents the function to send message to messenger.
type SendFunc func(toPK cipher.PubKey, t MsgType, body []byte) error

// RangeFunc wraps a master subscription action.
type RangeFunc func(action object.MasterSubAction) error

// DiscovererType represents the type of discoverer message.
type DiscovererType byte

const (
	// DiscovererAsk occurs when a node requests for board public keys.
	DiscovererAsk DiscovererType = iota << 0

	// DiscovererResponse occurs when a node broadcasts all it's master boards.
	DiscovererResponse
)

// DiscovererMsgTypeStr represents a DiscovererType as a string.
var DiscovererMsgTypeStr = [...]string{
	DiscovererAsk:      "Discoverer Ask",
	DiscovererResponse: "Discoverer Response",
}

// DiscoveredBoard represents a discovered board.
// It is signed to determine legitimacy of the provided board.
type DiscoveredBoard struct {
	PubKey    cipher.PubKey `verify:"upk"`
	Sig       cipher.Sig    `verify:"sig"`
	TimeStamp int64
}

// Fill prepares the DiscoveredBoard with provided timestamp and keys.
func (b *DiscoveredBoard) Fill(ts int64, pk cipher.PubKey, sk cipher.SecKey) *DiscoveredBoard {
	b.TimeStamp = ts
	tag.Sign(b, pk, sk)
	return b
}

// Verify verifies the validity of the DiscoveredBoard with provided timestamp.
func (b DiscoveredBoard) Verify(ts int64) error {
	if b.TimeStamp != ts {
		return boo.New(boo.NotAuthorised, "unmatched timestamps")
	}
	return tag.Verify(&b)
}

// DiscovererMsg represents an outgoing message of nature "discoverer".
type DiscovererMsg struct {
	Type      DiscovererType
	TimeStamp int64
	Boards    []DiscoveredBoard
}

// Add adds a Board to 'DiscovererMsg.Boards'. This is only relevant if
// message is of type 'DiscovererResponse'.
func (m *DiscovererMsg) Add(pk cipher.PubKey, sk cipher.SecKey) {
	m.Boards = append(m.Boards,
		*new(DiscoveredBoard).Fill(m.TimeStamp, pk, sk))
}

// Verify determines the legitimacy of the message.
func (m *DiscovererMsg) Verify() error {
	for i, b := range m.Boards {
		if e := b.Verify(m.TimeStamp); e != nil {
			return boo.WrapTypef(e, boo.NotAuthorised,
				"failed at index %d", i)
		}
	}
	return nil
}

func (m *DiscovererMsg) RangeValid(ts int64, action func(board *DiscoveredBoard)) {
	m.TimeStamp = ts
	for i, board := range m.Boards {
		if e := board.Verify(m.TimeStamp); e != nil {
			log.Printf("failed at index %d: %v", i, e)
		} else {
			action(&board)
		}
	}
}

// BoardDiscoverer is responsible for discovering boards.
type BoardDiscoverer struct {
	l *log.Logger

	nodes   *pkMap
	boards  *pkMap
	doSend  SendFunc
	doRange RangeFunc

	tsMux sync.RWMutex
	ts    int64
}

// NewBoardDiscoverer creates a new BoardDiscoverer.
func NewBoardDiscoverer(doSend SendFunc, doRange RangeFunc) *BoardDiscoverer {
	return &BoardDiscoverer{
		l:       inform.NewLogger(true, os.Stdout, "BOARD_DISCOVERER"),
		nodes:   newPKMap(),
		boards:  newPKMap(),
		doSend:  doSend,
		doRange: doRange,
	}
}

func (bd *BoardDiscoverer) ClearNodes() {
	bd.nodes.Clear()
}

func (bd *BoardDiscoverer) AddNode(npk cipher.PubKey) {
	bd.nodes.Add(npk)
}

// SendAsk sends a message to ask for boards.
func (bd *BoardDiscoverer) SendAsk() {
	bd.tsMux.Lock()
	defer bd.tsMux.Unlock()

	bd.ts = int64(time.Now().UnixNano())
	outMsg := &DiscovererMsg{
		Type:      DiscovererAsk,
		TimeStamp: bd.ts,
	}

	bd.nodes.Range(func(key cipher.PubKey) {
		bd.doSend(key, MsgDiscoverer, encoder.Serialize(outMsg))
	})
}

// SendResponse should be called after a SendAsk is received.
func (bd *BoardDiscoverer) SendResponse(nodeKey cipher.PubKey, ts int64) {

	respMsg := &DiscovererMsg{
		Type:      DiscovererResponse,
		TimeStamp: ts,
	}

	bd.doRange(func(pk cipher.PubKey, sk cipher.SecKey) {
		respMsg.Boards = append(respMsg.Boards,
			*new(DiscoveredBoard).Fill(ts, pk, sk))
	})

	bd.doSend(nodeKey, MsgDiscoverer, encoder.Serialize(respMsg))
}

// Process processes message.
func (bd *BoardDiscoverer) Process(fromNode cipher.PubKey, msg *DiscovererMsg) {
	switch msg.Type {
	case DiscovererAsk:
		bd.SendResponse(fromNode, msg.TimeStamp)

	case DiscovererResponse:
		msg.RangeValid(bd.getTimeStamp(), func(board *DiscoveredBoard) {
			bd.boards.Add(board.PubKey)
		})
	}
}

// GetBoards.
func (bd *BoardDiscoverer) GetBoards() []string {
	var out []string
	bd.boards.Range(func(key cipher.PubKey) {
		out = append(out, key.Hex())
	})
	return out
}

func (bd *BoardDiscoverer) getTimeStamp() int64 {
	bd.tsMux.RLock()
	defer bd.tsMux.RUnlock()
	return bd.ts
}

func (bd *BoardDiscoverer) setTimeStamp(ts int64) {
	bd.tsMux.Lock()
	defer bd.tsMux.Unlock()
	bd.ts = ts
}
