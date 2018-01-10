package cxo

import (
	"context"
	"github.com/skycoin/bbs/src/accord"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/misc/tag"
	"github.com/skycoin/bbs/src/store/cxo/setup"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/node/log"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/skycoin/src/cipher"
	"io/ioutil"
	log2 "log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	LogPrefix                = "CXO"
	SubDir                   = "cxo_v5"
	FileName                 = "bbs.json"
	ExportSubDir             = "exports"
	ExportFileExt            = ".export"
	BashAutoCompleteFileName = "bash_autocomplete"
	RetryDuration            = time.Second * 5
)

// ManagerConfig represents the configuration for CXO Manager.
type ManagerConfig struct {
	Public                     *bool    // Whether to expose node publicly.
	Memory                     *bool    // Whether to enable memory mode.
	EnforcedMessengerAddresses []string // Messenger addresses.
	EnforcedSubscriptions      []string // Subscriptions
	Config                     *string  // Configuration directory.
	CXOPort                    *int     // CXO listening port.
	CXORPCEnable               *bool    // Whether to enable CXO RPC.
	CXORPCPort                 *int     // CXO RPC port.
}

// Manager manages interaction with CXO and storing/retrieving node configuration files.
type Manager struct {
	mux      sync.Mutex
	c        *ManagerConfig
	l        *log2.Logger
	file     *object.CXOFileManager
	node     *node.Node
	compiler *state.Compiler
	relay    *accord.Relay
	wg       sync.WaitGroup
	newRoots chan state.RootWrap
	quit     chan struct{}
}

// NewManager creates a new CXO manager.
func NewManager(config *ManagerConfig, compilerConfig *state.CompilerConfig) *Manager {
	manager := &Manager{
		c: config,
		l: inform.NewLogger(true, os.Stdout, LogPrefix),
		file: object.NewCXOFileManager(&object.CXOFileManagerConfig{
			Memory: config.Memory,
		}),
		relay:    accord.NewRelay(),
		newRoots: make(chan state.RootWrap, 10),
		quit:     make(chan struct{}),
	}

	// Prepare CXO node.
	if e := manager.prepareNode(); e != nil {
		manager.l.Panicln("failed to start CXO manager:", e)
	}

	// Prepare CXO compiler.
	manager.compiler = state.NewCompiler(compilerConfig, manager.file, manager.newRoots, manager.node)

	// Prepare messenger relay.
	if e := manager.relay.Open(manager.compiler); e != nil {
		manager.l.Panicln("failed to start CXO manager:", e)
	}

	// Prepare CXO file.
	if e := manager.prepareFile(); e != nil {
		manager.l.Panicln("failed to start CXO manager:", e)
	}

	// Init directories.
	if !*config.Memory {
		if e := os.MkdirAll(
			path.Join(*config.Config, SubDir),
			os.FileMode(0700),
		); e != nil {
			manager.l.Panicln(e)
		}
		if e := os.MkdirAll(
			path.Join(*config.Config, ExportSubDir),
			os.FileMode(0700),
		); e != nil {
			manager.l.Panicln(e)
		}
	}

	go manager.retryLoop()
	go manager.relayLoop()
	return manager
}

// Close quits the CXO manager.
func (m *Manager) Close() {
	close(m.quit)
	m.relay.Close()
	m.compiler.Close()
	m.wg.Wait()
	if e := m.node.Close(); e != nil {
		m.l.Println("Error on close:", e.Error())
	}
	<-m.node.Quiting()
}

// Relay obtains the messenger relay.
func (m *Manager) Relay() *accord.Relay {
	return m.relay
}

/*
	<<< HELPER FUNCTIONS >>>
*/

// Sets up the CXO node.
func (m *Manager) prepareNode() error {
	c := node.NewConfig()

	c.Log.Prefix = "[CXO] "
	c.Log.Debug = false
	c.Log.Pins = log.All // all
	c.Config.Logger = log.NewLogger(log.Config{Output: ioutil.Discard})

	c.Skyobject.Log.Debug = true
	c.Skyobject.Log.Pins = skyobject.PackSavePin // all
	c.Skyobject.Log.Prefix = "[CXO] "

	// CXO schema are registered here.
	c.Skyobject.Registry = skyobject.NewRegistry(
		setup.PrepareRegistry)

	//c.MaxMessageSize = 0 // TODO -> Adjust.
	c.InMemoryDB = *m.c.Memory
	c.DataDir = filepath.Join(*m.c.Config, SubDir)
	c.EnableListener = true
	c.PublicServer = *m.c.Public
	c.DiscoveryAddresses = m.c.EnforcedMessengerAddresses
	c.Listen = "[::]:" + strconv.Itoa(*m.c.CXOPort)
	c.EnableRPC = *m.c.CXORPCEnable
	c.RemoteClose = false
	c.RPCAddress = "[::]:" + strconv.Itoa(*m.c.CXORPCPort)

	c.OnRootReceived = func(c *node.Conn, root *skyobject.Root) {
		m.l.Printf("Receiving board '%s' (NOT FILLED)", root.Pub.Hex()[:5]+"...")
		root.IsFull = false
		m.newRoots <- state.RootWrap{Root: root}
	}

	c.OnRootFilled = func(c *node.Conn, root *skyobject.Root) {
		m.l.Printf("Receiving board '%s' (FILLED)", root.Pub.Hex()[:5]+"...")
		root.IsFull = true
		m.newRoots <- state.RootWrap{Root: root}
	}

	c.OnCreateConnection = func(c *node.Conn) {
		for _, pk := range m.node.Feeds() {
			c.Subscribe(pk)
		}
		m.file.SetConnectionStatus(c.Address(), true)
	}

	c.OnCloseConnection = func(c *node.Conn) {
		m.file.SetConnectionStatus(c.Address(), false)
		go func() {
			time.Sleep(time.Second * 2)
			m.node.Connect(c.Address())
		}()
	}

	c.OnSubscribeRemote = func(c *node.Conn, feed cipher.PubKey) error {
		return nil
	}

	var e error
	if m.node, e = node.NewNode(c); e != nil {
		return e
	}

	return nil
}

// Sets up CXO file.
func (m *Manager) prepareFile() error {
	//if e := m.file.EnsureBashAutoComplete(path.Join(*m.c.Config, BashAutoCompleteFileName)); e != nil {
	//	return e
	//}
	if e := m.file.Load(m.filePath()); e != nil {
		return e
	}

	// Ensure messenger addresses and subscriptions.
	for _, address := range m.c.EnforcedMessengerAddresses {
		m.ConnectToMessenger(address)
	}
	for _, pkStr := range m.c.EnforcedSubscriptions {
		pk, e := tag.GetPubKey(pkStr)
		if e != nil {
			return e
		}
		m.SubscribeRemote(pk)
	}

	if e := m.file.RangeMasterSubs(func(pk cipher.PubKey, sk cipher.SecKey) {
		if r, e := m.node.Container().LastRoot(pk); e != nil {
			m.l.Println("prepareFile() LastRoot failed with error:", e)
		} else {
			m.compiler.UpdateBoard(r)
		}
	}); e != nil {
		return e
	}
	if e := m.file.RangeRemoteSubs(func(pk cipher.PubKey) {
		if e := m.subscribeNode(pk); e != nil {
			m.l.Println("prepareFile() subscribeNode() failed with error:", e)
		}
		if r, e := m.node.Container().LastRoot(pk); e != nil {
			m.l.Println("prepareFile() LastRoot failed with error:", e)
		} else {
			m.compiler.UpdateBoard(r)
		}
	}); e != nil {
		return e
	}
	return nil
}

func (m *Manager) filePath() string {
	return path.Join(*m.c.Config, SubDir, FileName)
}

func (m *Manager) exportPath(name string) string {
	return path.Join(*m.c.Config, ExportSubDir, name+ExportFileExt)
}

func (m *Manager) retryLoop() {
	m.wg.Add(1)
	defer m.wg.Done()

	ticker := time.NewTicker(RetryDuration)
	defer ticker.Stop()

	for {
		select {
		case <-m.quit:
			m.file.Save(m.filePath())
			return
		case <-ticker.C:
			m.file.Save(m.filePath())
		}
	}
}

func (m *Manager) relayLoop() {
	m.wg.Add(1)
	defer m.wg.Done()

	syncTicker := time.NewTicker(time.Second * 2)
	defer syncTicker.Stop()

	for {
		select {
		case <-m.quit:
			return
		case <-syncTicker.C:
			m.file.RangeMessengers(func(address string, pk cipher.PubKey) {
				if pk == (cipher.PubKey{}) {
					if pk, e := m.relay.Connect(address); e != nil {
						m.l.Printf("FAILED: messenger server connection: address(%s), error: %v",
							address, e)
					} else {
						m.l.Printf("SUCCESS: messenger server connection: address(%s) pk(%s)",
							address, pk.Hex())
						m.file.UnsafeSetMessengerPK(address, pk)
						go m.compiler.EnsureSubmissionKeys(m.relay.SubmissionKeys())
					}
					m.node.ConnectToMessenger(address)
				}
			})
			m.file.RangeConnections(func(address string, status bool) {
				if status == false {
					m.node.Pool().Dial(address)
				}
			})
		case address := <-m.relay.Disconnections():
			m.file.SetMessengerPK(address, cipher.PubKey{})
			go m.compiler.EnsureSubmissionKeys(m.relay.SubmissionKeys())
		}
	}
}

/*
	<<< MESSENGER >>>
*/

func (m *Manager) ConnectToMessenger(address string) error {
	if e := m.file.AddMessenger(strings.TrimSpace(address)); e != nil {
		return e
	}
	return nil
}

func (m *Manager) DisconnectFromMessenger(address string) error {
	if e := m.file.RemoveMessenger(strings.TrimSpace(address)); e != nil {
		return e
	}
	if sd := m.node.Discovery(); sd != nil {
		sd.ForEachConn(func(conn *factory.Connection) {
			if conn.GetRemoteAddr().String() == strings.TrimSpace(address) {
				conn.Close()
			}
		})
	}
	return nil
}

func (m *Manager) GetMessengers() []*object.MessengerConnection {
	var out []*object.MessengerConnection
	m.file.RangeMessengers(func(address string, pk cipher.PubKey) {
		mc := &object.MessengerConnection{
			Address:   address,
			Connected: pk != (cipher.PubKey{}),
		}
		if mc.Connected {
			mc.PubKey = pk
			mc.PubKeyStr = pk.Hex()
		}
		out = append(out, mc)
	})
	return out
}

func (m *Manager) SubmitToRemote(
	ctx context.Context, subKeys []*object.MessengerSubKeyTransport, transport *object.Transport,
) (
	uint64, error,
) {
	m.l.Println("attempting to submit to remote...")

	// Ensure that we can submit.
	if len(subKeys) == 0 {
		return 0, boo.New(boo.NotAllowed, "no submission keys are provided")
	}

	// Obtain submission.
	submission := &accord.Submission{
		Raw: transport.Content.Body,
		Sig: transport.Header.GetSig(),
	}

	// See if we are connected to any of the submission keys.
	m.l.Println("\t- looping through provided submission keys...")

	for i, subKey := range subKeys {

		m.l.Printf("\t\t- [%d] key '%s'", i, subKey.ToMessengerSubKey())

		fromPK, ok := m.file.GetMessengerPK(subKey.Address)
		if !ok || fromPK == (cipher.PubKey{}) {
			m.l.Println("\t\t\t (SKIPPING)")
			continue
		}

		// Attempt to submit.
		goal, e := m.relay.SubmitToRemote(ctx, subKey.PubKey, submission)
		if e != nil {
			m.l.Println("\t\t\t- Failed to submit to remote, error: %v", e)
			m.l.Println("\t\t\t (SKIPPING)")
			continue
		}
		return goal, nil
	}

	// Attempt manual connections.
	for _, subKey := range subKeys {
		if _, e := m.relay.Connect(subKey.Address); e != nil {
			continue
		}
		goal, e := m.relay.SubmitToRemote(ctx, subKey.PubKey, submission)
		if e != nil {
			continue
		}
		return goal, nil
	}

	return 0, boo.New(boo.NotFound, "a valid connection to messenger server is not found")
}

func (m *Manager) GetAvailableBoards() []cipher.PubKey {

	temp := make(map[cipher.PubKey]struct{})
	if discovery := m.node.Discovery(); discovery != nil {
		discovery.ForEachConn(func(conn *factory.Connection) {
			for _, service := range conn.GetServices().Services {
				temp[service.Key] = struct{}{}
			}
		})
	}

	out := make([]cipher.PubKey, len(temp))
	i := 0
	for pk := range temp {
		out[i] = pk
		i += 1
	}
	return out
}

/*
	<<< CONNECTION >>>
*/

func (m *Manager) GetActiveConnections() []object.Connection {
	connections := m.node.Connections()
	out := make([]object.Connection, len(connections))
	for i, conn := range m.node.Connections() {
		out[i] = object.Connection{
			Address: conn.Address(),
			State:   conn.Gnet().State().String(),
		}
	}
	return out
}

func (m *Manager) GetSavedConnections() []object.Connection {
	out := make([]object.Connection, 0)
	m.file.RangeConnections(func(address string, status bool) {
		if conn := m.node.Connection(address); conn != nil && status == true {
			out = append(out, object.Connection{
				Address: address,
				State:   conn.Gnet().State().String(),
			})
		} else {
			out = append(out, object.Connection{
				Address: address,
				State:   "CLOSED",
			})
		}
	})
	return out
}

func (m *Manager) Connect(address string) error {
	return m.file.AddConnection(address)
}

//func (m *Manager) connectNode(address string) (*gnet.Conn, error) {
//	if connection, e := m.node.Pool().Dial(address); e != nil {
//		switch e {
//		case gnet.ErrClosed, gnet.ErrConnectionsLimit:
//			return nil, boo.WrapType(e, boo.Internal)
//		case gnet.ErrAlreadyListen:
//			return nil, boo.WrapType(e, boo.AlreadyExists)
//		default:
//			return nil, boo.WrapType(e, boo.InvalidInput)
//		}
//	} else {
//		return connection, nil
//	}
//}

func (m *Manager) Disconnect(address string) error {
	m.file.RemoveConnection(address)
	if e := m.disconnectNode(address); e != nil {
		return e
	}
	return nil
}

func (m *Manager) disconnectNode(address string) error {
	conn := m.node.Pool().Connection(address)
	if conn == nil {
		return nil
	}
	if e := conn.Close(); e != nil {
		return e
	}
	return nil
}

/*
	<<< SUBSCRIPTION >>>
*/

func (m *Manager) GetSubscriptions() []cipher.PubKey {
	out, _ := m.file.GetSubKeyList()
	return out
}

func (m *Manager) SubscribeRemote(bpk cipher.PubKey) error {
	if e := m.file.AddRemoteSub(bpk); e != nil {
		return e
	}
	m.subscribeNode(bpk)
	return nil
}

func (m *Manager) SubscribeMaster(bpk cipher.PubKey, bsk cipher.SecKey) error {
	if e := m.file.AddMasterSub(bpk, bsk); e != nil {
		return e
	}
	m.subscribeNode(bpk)
	return nil
}

func (m *Manager) subscribeNode(bpk cipher.PubKey) error {
	if e := m.node.AddFeed(bpk); e != nil {
		return boo.WrapType(e, boo.Internal, "failed to add feed")
	}
	return nil
}

func (m *Manager) UnsubscribeRemote(bpk cipher.PubKey) error {
	if m.file.HasRemoteSub(bpk) {
		if e := m.file.RemoveSub(bpk); e != nil {
			return e
		}
		m.unsubscribeNode(bpk)
		return nil
	}
	return boo.Newf(boo.NotFound,
		"remote board of public key '%s' not found in cxo file", bpk.Hex()[:5]+"...")
}

func (m *Manager) UnsubscribeMaster(bpk cipher.PubKey) error {
	if m.file.HasMasterSub(bpk) {
		if e := m.file.RemoveSub(bpk); e != nil {
			return e
		}
		m.unsubscribeNode(bpk)
		return nil
	}
	return boo.Newf(boo.NotFound,
		"master board of public key '%s' not found in cxo file", bpk.Hex()[:5]+"...")
}

func (m *Manager) unsubscribeNode(bpk cipher.PubKey) {
	m.compiler.DeleteBoard(bpk)
	m.node.DelFeed(bpk)
}

/*
	<<< CONTENT >>>
*/

func (m *Manager) GetBoardInstance(bpk cipher.PubKey) (*state.BoardInstance, error) {
	return m.compiler.GetBoard(bpk)
}

func (m *Manager) GetBoards(ctx context.Context) ([]interface{}, []interface{}, error) {

	var masterOut = []interface{}{}
	m.file.RangeMasterSubs(func(pk cipher.PubKey, sk cipher.SecKey) {
		bi, e := m.compiler.GetBoard(pk)
		if e != nil {
			m.l.Println(e)
			return
		}
		bView, e := bi.Viewer().GetBoard()
		if e != nil {
			m.l.Println(e)
			return
		}
		masterOut = append(masterOut, bView)
	})

	var remoteOut = []interface{}{}
	m.file.RangeRemoteSubs(func(pk cipher.PubKey) {
		bi, e := m.compiler.GetBoard(pk)
		if e != nil {
			m.l.Println(e)
			return
		}
		bView, e := bi.Viewer().GetBoard()
		if e != nil {
			m.l.Println(e)
			return
		}
		remoteOut = append(remoteOut, bView)
	})

	return masterOut, remoteOut, nil
}

func (m *Manager) NewBoard(content *object.Content, pk cipher.PubKey, sk cipher.SecKey) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	if e := m.file.AddMasterSub(pk, sk); e != nil {
		return e
	}
	m.subscribeNode(pk)

	if r, e := setup.NewBoard(m.node, content, pk, sk); e != nil {
		return e
	} else {
		m.compiler.UpdateBoardWithContext(context.Background(), r)
		return nil
	}
}

/*
	<<< ADMIN >>>
*/

/*
	<<< IMPORT / EXPORT >>>
*/

func (m *Manager) ExportBoard(pk cipher.PubKey, path string) (*object.PagesJSON, error) {
	sk, _ := m.file.GetMasterSubSecKey(pk)
	bi, e := m.GetBoardInstance(pk)
	if e != nil {
		return nil, e
	}
	out, e := bi.Export(pk, sk)
	if e != nil {
		return nil, e
	}
	return out, nil
}

func (m *Manager) ImportBoard(ctx context.Context, in *object.PagesJSON) error {
	var (
		pk = in.GetPubKey()
		sk = in.GetSecKey()
	)
	if cipher.PubKeyFromSecKey(sk) != pk {
		return boo.New(boo.InvalidRead,
			"public key does not match secret key in exported board file")
	}
	if m.file.HasRemoteSub(pk) {
		m.unsubscribeNode(pk)
	}
	if m.file.HasMasterSub(pk) == false {
		content := new(object.Content)
		content.SetHeader(&object.ContentHeaderData{})
		content.SetBody(&object.Body{})
		if e := m.NewBoard(content, pk, sk); e != nil {
			return e
		}
	}
	bi, e := m.GetBoardInstance(pk)
	if e != nil {
		return e
	}
	goal, e := bi.Import(in)
	if e != nil {
		return e
	}
	return bi.WaitSeq(ctx, goal)
}
