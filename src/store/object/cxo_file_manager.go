package object

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/bbs/src/misc/typ"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
	"log"
	"os"
	"sync"
)

const (
	cxoFileManagerLogPrefix = "CXOFILEMANAGER"
)

var (
	defaultSubscriptions = []string{
		"03588a2c8085e37ece47aec50e1e856e70f893f7f802cb4f92d52c81c4c3212742",
	}
	defaultMessengers = []string{
		"messenger.skycoin.net:8080",
	}
	defaultDevMessengers = []string{
		"127.0.0.1:8080",
	}
)

// CXOFileManagerConfig configures the CXOFileManager.
type CXOFileManagerConfig struct {
	Memory   *bool // Whether to run in memory mode.
	Defaults *bool
	Dev      *bool
}

// CXOFileManager manages the CXOFile.
// This is a file containing saved connections and board keys.
type CXOFileManager struct {
	c          *CXOFileManagerConfig
	l          *log.Logger
	mux        sync.Mutex
	hasChanges bool // has changes.
	masters    *typ.List
	remotes    *typ.List
	messengers *typ.List
}

// NewCXOFileManager creates a new file manager with provided configuration.
func NewCXOFileManager(config *CXOFileManagerConfig) *CXOFileManager {
	return &CXOFileManager{
		c:          config,
		l:          inform.NewLogger(true, os.Stdout, cxoFileManagerLogPrefix),
		masters:    typ.NewList(),
		remotes:    typ.NewList(),
		messengers: typ.NewList(),
	}
}

// Load loads the configuration from file (if not in memory mode).
func (m *CXOFileManager) Load(path string) error {
	defer m.lock()()
	if m.memMode() == false {
		return m.load(path)
	}
	return nil
}

// Save saves the configuration to file (if not in memory mode).
func (m *CXOFileManager) Save(path string) error {
	defer m.lock()()
	if m.memMode() == false && m.hasChanges {
		m.untagChanges()
		return m.save(path)
	}
	return nil
}

// AddMasterSub adds a master subscription to file.
func (m *CXOFileManager) AddMasterSub(pk cipher.PubKey, sk cipher.SecKey) error {
	defer m.lock()()

	// Make sure this public key does not exist in remotes.
	m.remotes.DelOfKey(pk)

	// Append to masters.
	if m.masters.Append(pk, &Subscription{PK: pk, SK: sk}) == false {
		return boo.Newf(boo.AlreadyExists,
			"board of public key '%s' already exists in master subscriptions",
			pk.Hex())
	}

	// Record changes and return.
	m.tagChanges()
	return nil
}

// AddRemoteSub adds a remote subscription to file.
func (m *CXOFileManager) AddRemoteSub(pk cipher.PubKey) error {
	defer m.lock()()

	if m.masters.HasKey(pk) {
		return boo.Newf(boo.NotAllowed,
			"file already has master subscription to '%s'", pk.Hex())
	}

	if m.remotes.Append(pk, &Subscription{PK: pk}) == false {
		return boo.Newf(boo.AlreadyExists,
			"file already has remote subscription to '%s'", pk.Hex())
	}

	m.tagChanges()
	return nil
}

// RemoveSub removes a subscription from file, whether master or remote.
func (m *CXOFileManager) RemoveSub(pk cipher.PubKey) error {
	defer m.lock()()

	m.masters.DelOfKey(pk)
	m.remotes.DelOfKey(pk)

	m.tagChanges()
	return nil
}

// GetMasterSubs returns list of master subscriptions.
func (m *CXOFileManager) GetMasterSubs() ([]*Subscription, error) {
	defer m.lock()()

	out := make([]*Subscription, m.masters.Len())
	e := m.masters.Range(typ.Ascending, func(i int, _, v interface{}) (bool, error) {
		var ok bool
		if out[i], ok = v.(*Subscription); !ok {
			return false, boo.Newf(boo.Internal,
				"failed to extract from master_subscriptions[%d]", i)
		}
		return false, nil
	})

	return out, e
}

// GetRemoteSubs returns list of remote subscriptions.
func (m *CXOFileManager) GetRemoteSubs() ([]*Subscription, error) {
	defer m.lock()()

	out := make([]*Subscription, m.remotes.Len())
	e := m.remotes.Range(typ.Ascending, func(i int, _, v interface{}) (bool, error) {
		var ok bool
		if out[i] = v.(*Subscription); !ok {
			return false, boo.Newf(boo.Internal,
				"failed to extract from remote_subscriptions[%d]", i)
		}
		return false, nil
	})

	return out, e
}

// MasterSubsLen gets the number of master subscriptions in file.
func (m *CXOFileManager) MasterSubsLen() int {
	defer m.lock()()
	return m.masters.Len()
}

// RemoteSubsLen gets the number of remote subscriptions in file.
func (m *CXOFileManager) RemoteSubsLen() int {
	defer m.lock()()
	return m.remotes.Len()
}

// MasterSubAction performs an action to master subscription.
type MasterSubAction func(pk cipher.PubKey, sk cipher.SecKey)

// RangeMasterSubs ranges all master subscriptions.
func (m *CXOFileManager) RangeMasterSubs(action MasterSubAction) error {
	defer m.lock()()
	return m.masters.Range(typ.Ascending, func(i int, k, v interface{}) (bool, error) {
		sub, ok := v.(*Subscription)
		if !ok {
			return false, boo.Newf(boo.Internal,
				"failed to extract value from master_subscriptions[%d:'%s']",
				i, k.(cipher.PubKey).Hex()[:5]+"...")
		}
		action(sub.PK, sub.SK)
		return false, nil
	})
}

// RemoteSubAction performs an action on a remote subscription.
type RemoteSubAction func(pk cipher.PubKey)

// RangeRemoteSubs ranges all remote subscriptions.
func (m *CXOFileManager) RangeRemoteSubs(action RemoteSubAction) error {
	defer m.lock()()
	for i, k := range m.remotes.Keys() {
		pk, ok := k.(cipher.PubKey)
		if !ok {
			return boo.Newf(boo.Internal,
				"failed to extract key from remote_subscriptions[%d]", i)
		}
		action(pk)
	}
	return nil
}

// GetSubKeyList gets a list of all subscription public keys.
func (m *CXOFileManager) GetSubKeyList() ([]cipher.PubKey, error) {
	defer m.lock()()
	out := make([]cipher.PubKey, m.masters.Len()+m.remotes.Len())
	m.masters.Range(typ.Ascending, func(i int, key, _ interface{}) (bool, error) {
		out[i] = key.(cipher.PubKey)
		return false, nil
	})
	m.remotes.Range(typ.Ascending, func(i int, key, _ interface{}) (bool, error) {
		out[i+m.masters.Len()] = key.(cipher.PubKey)
		return false, nil
	})
	return out, nil
}

// HasMasterSub determines whether we are subscribed to this as master.
func (m *CXOFileManager) HasMasterSub(pk cipher.PubKey) bool {
	return m.masters.HasKey(pk)
}

// HasRemoteSub determines whether we are subscribed to this as remote.
func (m *CXOFileManager) HasRemoteSub(pk cipher.PubKey) bool {
	return m.remotes.HasKey(pk)
}

// GetMasterSubSecKey obtains secret key of master if it exists.
func (m *CXOFileManager) GetMasterSubSecKey(pk cipher.PubKey) (cipher.SecKey, bool) {
	defer m.lock()()

	v, ok := m.masters.GetOfKey(pk)
	if !ok {
		return cipher.SecKey{}, false
	}

	return v.(*Subscription).SK, true
}

// AddMessenger adds a messenger server address to connect to.
func (m *CXOFileManager) AddMessenger(address string) error {
	defer m.lock()()

	// Append.
	if m.messengers.Append(address, cipher.PubKey{}) == false {
		return boo.Newf(boo.AlreadyExists,
			"address '%s' already exists in messenger servers",
			address)
	}

	// Record changes and return.
	m.tagChanges()
	return nil
}

// RemoveMessenger removes a messenger server address.
func (m *CXOFileManager) RemoveMessenger(address string) error {
	defer m.lock()()

	m.messengers.DelOfKey(address)
	m.tagChanges()
	return nil
}

// MessengerAction represents an action to perform on a messenger address and key.
type MessengerAction func(address string, pk cipher.PubKey)

// RangeMessengers ranges the list of messenger server addresses.
func (m *CXOFileManager) RangeMessengers(action MessengerAction) error {
	defer m.lock()()
	m.messengers.Range(typ.Ascending, func(i int, key, value interface{}) (bool, error) {
		action(key.(string), value.(cipher.PubKey))
		return false, nil
	})
	return nil
}

// SetMessengerPK sets the messenger's public key.
func (m *CXOFileManager) SetMessengerPK(address string, pk cipher.PubKey) error {
	defer m.lock()()
	return m.UnsafeSetMessengerPK(address, pk)
}

func (m *CXOFileManager) UnsafeSetMessengerPK(address string, pk cipher.PubKey) error {
	if m.messengers.Replace(address, pk) == false {
		return boo.Newf(boo.InvalidInput,
			"messenger server address '%s' is not found", address)
	}
	return nil
}

// GetMessengerPK gets the messenger's public key.
func (m *CXOFileManager) GetMessengerPK(address string) (cipher.PubKey, bool) {
	defer m.lock()()
	if pk, ok := m.messengers.GetOfKey(address); !ok {
		return cipher.PubKey{}, false
	} else {
		return pk.(cipher.PubKey), true
	}
}

func (m *CXOFileManager) GetMessengersLen() int {
	defer m.lock()()
	return m.messengers.Len()
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func (m *CXOFileManager) load(path string) error {
	var fileData CXOFile

	// Load from file - If file does not exist, this is okay.
	if e := file.LoadJSON(path, &fileData); e != nil && !os.IsExist(e) {
		if !os.IsNotExist(e) {
			return boo.WrapTypef(e, boo.InvalidRead,
				"failed to read CXO file from '%s'", path)
		} else if *m.c.Defaults {
			// Load defaults.
			m.l.Println("First Run - Loading default subscriptions:")
			for i, pkStr := range defaultSubscriptions {
				pk, _ := keys.GetPubKey(pkStr)
				m.l.Printf(" - [%d] Subscription '%s'", i, pkStr[:5]+"...")
				m.remotes.Append(pk, &Subscription{PK: pk})
			}
			if *m.c.Dev {
				m.l.Println("First Run - Loading default messenger addresses (dev):")
				for i, address := range defaultDevMessengers {
					m.l.Printf(" - [%d] Messenger address '%s' (dev)", i, address)
					m.messengers.Append(address, cipher.PubKey{})
				}
			} else {
				m.l.Println("First Run - Loading default messenger addresses:")
				for i, address := range defaultMessengers {
					m.l.Printf(" - [%d] Messenger address '%s'", i, address)
					m.messengers.Append(address, cipher.PubKey{})
				}
			}
			m.tagChanges()
			return nil
		}
	}

	// Load to memory.

	// Range master subscriptions.
	for i, sub := range fileData.MasterSubs {

		// Get public key of master subscription.
		pk, e := keys.GetPubKey(sub.PK)
		if e != nil {
			return boo.WrapTypef(e, boo.InvalidRead,
				"invalid public key in file at master_subscriptions[%d]", i)
		}

		// Get private key of master subscription.
		sk, e := keys.GetSecKey(sub.SK)
		if e != nil {
			return boo.WrapTypef(e, boo.InvalidRead,
				"invalid private key in file at master_subscriptions[%d]", i)
		}

		// Append.
		m.masters.Append(pk, &Subscription{PK: pk, SK: sk})
	}

	// Range remote subscriptions.
	for i, sub := range fileData.RemoteSubs {

		// Get public key of remote subscription.
		pk, e := keys.GetPubKey(sub.PK)
		if e != nil {
			return boo.WrapTypef(e, boo.InvalidRead,
				"invalid public key in file at remote_subscriptions[%d]", i)
		}

		// Append.
		m.remotes.Append(pk, &Subscription{PK: pk})
	}

	// Range messenger addresses.
	for i, address := range fileData.MessengerAddresses {
		if e := tag.CheckAddress(address); e != nil {
			return boo.WrapType(e, boo.InvalidRead,
				"invalid address in file at messenger_addresses[%d]", i)
		}
		m.messengers.Append(address, cipher.PubKey{})
	}

	return nil
}

func (m *CXOFileManager) save(path string) error {
	var fileData CXOFile

	fileData.MasterSubs = make([]SubscriptionView, m.masters.Len())
	m.masters.Range(typ.Ascending, func(i int, _, v interface{}) (bool, error) {
		fileData.MasterSubs[i] = v.(*Subscription).View()
		return false, nil
	})

	fileData.RemoteSubs = make([]SubscriptionView, m.remotes.Len())
	m.remotes.Range(typ.Ascending, func(i int, _, v interface{}) (bool, error) {
		fileData.RemoteSubs[i] = v.(*Subscription).View()
		return false, nil
	})

	fileData.MessengerAddresses = make([]string, m.messengers.Len())
	m.messengers.Range(typ.Ascending, func(i int, k, _ interface{}) (bool, error) {
		fileData.MessengerAddresses[i] = k.(string)
		return false, nil
	})

	if e := file.SaveJSON(path, fileData, os.FileMode(0600)); e != nil {
		return boo.WrapTypef(e, boo.Internal,
			"failed to save CXO file to '%s'", path)
	}
	return nil
}

func (m *CXOFileManager) tagChanges() {
	m.hasChanges = true
}

func (m *CXOFileManager) untagChanges() {
	m.hasChanges = false
}

func (m *CXOFileManager) lock() func() {
	m.mux.Lock()
	return m.mux.Unlock
}

func (m *CXOFileManager) memMode() bool {
	return *m.c.Memory
}
