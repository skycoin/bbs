package cxo

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type Try struct {
	Connections   map[string]bool
	Subscriptions map[cipher.PubKey]bool
	cMux, sMux    sync.Mutex
}

func NewTry() *Try {
	return &Try{
		Connections:   make(map[string]bool),
		Subscriptions: make(map[cipher.PubKey]bool),
	}
}

func (t *Try) fill(f *object.CXOFile) {
	f.Lock()
	defer f.Unlock()
	for _, sub := range f.MasterSubs {
		t.AddSub(sub.PK)
	}
	for _, sub := range f.RemoteSubs {
		t.AddSub(sub.PK)
	}
	for _, conn := range f.Connections {
		t.AddCon(conn)
	}
}

func (t *Try) AddCon(connection string) {
	t.cMux.Lock()
	defer t.cMux.Unlock()
	t.Connections[connection] = true
}

func (t *Try) DelCon(connection string) {
	t.cMux.Lock()
	defer t.cMux.Unlock()
	delete(t.Connections, connection)
}

func (t *Try) HasCon(connection string) bool {
	t.cMux.Lock()
	defer t.cMux.Unlock()
	return t.Connections[connection]
}

func (t *Try) AddSub(subscription cipher.PubKey) {
	t.sMux.Lock()
	defer t.sMux.Unlock()
	t.Subscriptions[subscription] = true
}

func (t *Try) DelSub(subscription cipher.PubKey) {
	t.sMux.Lock()
	defer t.sMux.Unlock()
	delete(t.Subscriptions, subscription)
}

func (t *Try) HasSub(subscription cipher.PubKey) bool {
	t.sMux.Lock()
	defer t.sMux.Unlock()
	return t.Subscriptions[subscription]
}
