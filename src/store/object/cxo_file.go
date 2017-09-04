package object

import (
	"github.com/skycoin/skycoin/src/cipher"
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
	if s.SK == (cipher.SecKey{}) {
		return SubscriptionView{
			PK: s.PK.Hex(),
		}
	} else {
		return SubscriptionView{
			PK: s.PK.Hex(),
			SK: s.SK.Hex(),
		}
	}
}

type CXOFile struct {
	MasterSubs  []SubscriptionView `json:"master_subscriptions"`
	RemoteSubs  []SubscriptionView `json:"remote_subscriptions"`
	Connections []string           `json:"connections"`
}

//func (f *CXOFile) Load(path string) error {
//	f.mux.Lock()
//	defer f.mux.Unlock()
//
//	if e := file.LoadJSON(path, f); e != nil {
//		return boo.WrapType(e, boo.InvalidRead,
//			"failed to read CXO file")
//	}
//
//	f.masterSubMap = make(map[cipher.PubKey]*Subscription)
//	f.remoteSubMap = make(map[cipher.PubKey]*Subscription)
//	f.connectionMap = make(map[string]bool)
//
//	for _, sub := range f.MasterSubs {
//		f.masterSubMap[sub.PK] = sub
//	}
//	for _, sub := range f.RemoteSubs {
//		f.remoteSubMap[sub.PK] = sub
//	}
//	for _, address := range f.Connections {
//		f.connectionMap[address] = true
//	}
//
//	return nil
//}
//
//func (f *CXOFile) Save(path string) error {
//	f.mux.Lock()
//	defer f.mux.Unlock()
//
//	if e := file.SaveJSON(path, f, os.FileMode(0600)); e != nil {
//		return boo.WrapType(e, boo.NotAllowed,
//			"failed to save CXO file")
//	}
//	return nil
//}
//
//func (f *CXOFile) AddConnection(address string) bool {
//	f.mux.Lock()
//	defer f.mux.Unlock()
//
//	if _, has := f.connectionMap[address]; !has {
//		f.connectionMap[address] = true
//		f.Connections = append(f.Connections, address)
//		return true
//	}
//	return false
//}
//
//func (f *CXOFile) RemoveConnection(address string) bool {
//	f.mux.Lock()
//	defer f.mux.Unlock()
//
//	if _, has := f.connectionMap[address]; has {
//		for i := len(f.Connections); i >= 0; i-- {
//			if f.Connections[i] == address {
//				f.Connections[i], f.Connections[0] =
//					f.Connections[0], f.Connections[i]
//				f.Connections = f.Connections[1:]
//			}
//		}
//		delete(f.connectionMap, address)
//		return true
//	}
//	return false
//}
//
//func (f *CXOFile) AddRemoteSub(pk cipher.PubKey) bool {
//	f.mux.Lock()
//	defer f.mux.Unlock()
//
//	if f.hasSub(pk) == false {
//		sub := &Subscription{PK: pk}
//		f.RemoteSubs = append(f.RemoteSubs, sub)
//		f.remoteSubMap[pk] = sub
//		return true
//	}
//	return false
//}
//
//func (f *CXOFile) RemoveRemoteSub(pk cipher.PubKey) bool {
//	f.mux.Lock()
//	defer f.mux.Unlock()
//
//	if _, has := f.remoteSubMap[pk]; has {
//		for i, sub := range f.RemoteSubs {
//			if sub.PK == pk {
//				f.RemoteSubs[i], f.RemoteSubs[0] =
//					f.RemoteSubs[0], f.RemoteSubs[i]
//				f.RemoteSubs = f.RemoteSubs[1:]
//			}
//		}
//		delete(f.remoteSubMap, pk)
//		return true
//	}
//	return false
//}
//
//func (f *CXOFile) AddMasterSub(pk cipher.PubKey, sk cipher.SecKey) bool {
//	f.mux.Lock()
//	defer f.mux.Unlock()
//
//	if f.hasSub(pk) {
//		return false
//	}
//	f.MasterSubs = append(f.MasterSubs,
//		&Subscription{PK: pk, SK: sk})
//	return true
//}
//
//func (f *CXOFile) RemoveMasterSub(pk cipher.PubKey) bool {
//	f.Lock()
//	defer f.Unlock()
//
//	for i, sub := range f.MasterSubs {
//		if sub.PK == pk {
//			f.MasterSubs[0], f.MasterSubs[i] =
//				f.MasterSubs[i], f.MasterSubs[0]
//			f.MasterSubs = f.MasterSubs[1:]
//			return true
//		}
//	}
//	return false
//}
//
//func (f *CXOFile) HasSub(pk cipher.PubKey) bool {
//	f.mux.Lock()
//	defer f.mux.Unlock()
//	return f.hasSub(pk)
//}
//
//func (f *CXOFile) hasSub(pk cipher.PubKey) bool {
//	return f.hasMasterSub(pk) || f.hasRemoteSub(pk)
//}
//
//func (f *CXOFile) hasMasterSub(pk cipher.PubKey) bool {
//	_, has := f.masterSubMap[pk]
//	return has
//}
//
//func (f *CXOFile) hasRemoteSub(pk cipher.PubKey) bool {
//	_, has := f.remoteSubMap[pk]
//	return has
//}
//
//func (f *CXOFile) GetSub(pk cipher.PubKey) (*Subscription, bool, error) {
//	f.mux.Lock()
//	defer f.mux.Unlock()
//
//	if sub, has := f.masterSubMap[pk]; has {
//		return sub, true, nil
//	}
//	if sub, has := f.remoteSubMap[pk]; has {
//		return sub, false, nil
//	}
//
//	return nil, false, boo.Newf(boo.NotFound,
//		"subscription '%s' not found", pk.Hex())
//}
