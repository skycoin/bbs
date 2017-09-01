package object

import (
	"github.com/skycoin/skycoin/src/cipher"
	"sync"
)

type SubscriptionView struct {
	PK string `json:"public_key"`
	SK string `json:"secret_key,omitempty"`
}

type Subscription struct {
	PK cipher.PubKey
	SK cipher.SecKey
}

func (s *Subscription) View() SubscriptionView {
	return SubscriptionView{
		PK: s.PK.Hex(),
		SK: s.SK.Hex(),
	}
}

type CXOFileView struct {
	MasterSubs  []SubscriptionView
	RemoteSubs  []SubscriptionView
	Connections []string
}

type CXOFile struct {
	sync.Mutex
	MasterSubs  []Subscription
	RemoteSubs  []Subscription
	Connections []string
}

func (f *CXOFile) View() CXOFileView {
	f.Lock()
	defer f.Unlock()

	view := CXOFileView{
		Connections: f.Connections,
	}
	view.MasterSubs =
		make([]SubscriptionView, len(f.MasterSubs))
	for i, s := range f.MasterSubs {
		view.MasterSubs[i] = s.View()
	}
	view.RemoteSubs =
		make([]SubscriptionView, len(f.RemoteSubs))
	for i, s := range f.RemoteSubs {
		view.RemoteSubs[i] = s.View()
	}
	return view
}

func (f *CXOFile) AddConnection(address string) bool {
	f.Lock()
	defer f.Unlock()

	for _, conn := range f.Connections {
		if conn == address {
			return false
		}
	}
	f.Connections = append(f.Connections,
		address)
	return true
}

func (f *CXOFile) RemoveConnection(address string) bool {
	f.Lock()
	defer f.Unlock()

	for i, conn := range f.Connections {
		if conn == address {
			f.Connections[i], f.Connections[0] =
				f.Connections[0], f.Connections[i]
			f.Connections = f.Connections[1:]
			return true
		}
	}
	return false
}

func (f *CXOFile) AddRemoteSub(pk cipher.PubKey) bool {
	f.Lock()
	defer f.Unlock()

	if f.hasSub(pk) {
		return false
	}
	f.RemoteSubs = append(f.RemoteSubs,
		Subscription{PK: pk})
	return true
}

func (f *CXOFile) RemoveRemoteSub(pk cipher.PubKey) bool {
	f.Lock()
	defer f.Unlock()

	for i, sub := range f.RemoteSubs {
		if sub.PK == pk {
			f.RemoteSubs[i], f.RemoteSubs[0] =
				f.RemoteSubs[0], f.RemoteSubs[i]
			f.RemoteSubs = f.RemoteSubs[1:]
			return true
		}
	}
	return false
}

func (f *CXOFile) AddMasterSub(pk cipher.PubKey, sk cipher.SecKey) bool {
	f.Lock()
	defer f.Unlock()

	if f.hasSub(pk) {
		return false
	}
	f.MasterSubs = append(f.MasterSubs,
		Subscription{PK: pk, SK: sk})
	return true
}

func (f *CXOFile) RemoveMasterSub(pk cipher.PubKey) bool {
	f.Lock()
	defer f.Unlock()

	for i, sub := range f.MasterSubs {
		if sub.PK == pk {
			f.MasterSubs[0], f.MasterSubs[i] =
				f.MasterSubs[i], f.MasterSubs[0]
			f.MasterSubs = f.MasterSubs[1:]
			return true
		}
	}
	return false
}

func (f *CXOFile) HasSub(pk cipher.PubKey) bool {
	f.Lock()
	defer f.Unlock()
	return f.hasSub(pk)
}

func (f *CXOFile) hasSub(pk cipher.PubKey) bool {
	for _, sub := range append(f.RemoteSubs, f.MasterSubs...) {
		if pk == sub.PK {
			return true
		}
	}
	return false
}

