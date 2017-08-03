package session

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/file"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

const (
	nodeLogPrefix = "SESSION"
	nodeSubDir    = "/node"
	nodeFileName  = "bbsnode.session.v2"
	timeout       = time.Second * 10
	saveInternal  = time.Second * 10
	retryInterval = time.Second * 5
)

var (
	ErrNotLoggedIn     = boo.New(boo.NotAllowed, "not logged in")
	ErrAlreadyLoggedIn = boo.New(boo.NotAllowed, "already logged in")
)

// ManagerConfig configures a Manager.
type ManagerConfig struct {
	Master       *bool   // Whether node is master.
	TestMode     *bool   // Whether node is io test mode.
	MemoryMode   *bool   // Whether to use local storage io runtime.
	ConfigDir    *string // Configuration directory.
	CXOPort      *int    // CXO listening port.
	CXORPCEnable *bool   // Whether to enable CXO RPC.
	CXORPCPort   *int    // CXO RPC listening port.
}

// Manager manages file sessions.
type Manager struct {
	// Configuration.
	c *ManagerConfig
	l *log.Logger

	// Manager.
	cxo      *CXO
	compiler *state.Compiler
	file     *File
	changes  bool

	// Request & Response.
	retries      chan *object.RetryIO // Channel for retrying failed subscriptions/connections.
	clearRetries chan chan struct{}   // Channel for clearing retry queue (i.e. logout).
	requests     chan interface{}     // Channel for file requests.
	quit         chan struct{}        // Channel for quitting Manager.
	wg           sync.WaitGroup
}

// NewManager creates a new Manager with configuration.
func NewManager(config *ManagerConfig, compilerConfig *state.CompilerConfig) (*Manager, error) {
	m := &Manager{
		c:            config,
		l:            inform.NewLogger(true, os.Stdout, nodeLogPrefix),
		cxo:          NewCXO(config),
		compiler:     state.NewCompiler(compilerConfig, state.SetV1State()),
		file:         new(File),
		changes:      false,
		retries:      make(chan *object.RetryIO, 5),
		clearRetries: make(chan chan struct{}),
		requests:     make(chan interface{}, 5),
		quit:         make(chan struct{}),
	}
	m.cxo.SetUpdater(m.compiler.Trigger)

	if e := os.MkdirAll(m.folderPath(), os.FileMode(0700)); e != nil {
		return nil, e
	}

	m.processLoad()

	failed, e := m.cxo.Open(new(object.RetryIO).Fill(
		append(m.file.Masters, m.file.Subscriptions...),
		m.file.Connections,
	))
	if e != nil {
		return nil, e
	}

	go m.service()
	go m.retryLoop()

	if !failed.IsEmpty() {
		m.retries <- failed
	}

	return m, nil
}

// Close closes the Manager service.
func (m *Manager) Close() {
	for {
		select {
		case m.quit <- struct{}{}:
			continue
		default:
			m.wg.Wait()
			m.compiler.Close()
			if e := m.cxo.Close(); e != nil && e != ErrCXONotOpened {
				m.l.Println("Closing CXO Error:", e)
			}
			return
		}
	}
}

// GetCXO obtains the internal CXO.
func (m *Manager) GetCXO() *CXO { return m.cxo }

// GetCompiler obtains the internal Compiler.
func (m *Manager) GetCompiler() *state.Compiler { return m.compiler }

/*
	<<< HELPER FUNCTIONS >>>
*/

func (m *Manager) folderPath() string {
	return path.Join(*m.c.ConfigDir, nodeSubDir)
}

func (m *Manager) filePath() string {
	return path.Join(m.folderPath(), nodeFileName)
}

func (m *Manager) timedContext(ctx context.Context) context.Context {
	ctx, _ = context.WithTimeout(ctx, timeout)
	return ctx
}

func (m *Manager) processLoad() {
	if *m.c.MemoryMode {
		return
	}
	if e := file.LoadJSON(m.filePath(), m.file); e != nil {
		// Panic if error is other than "IsNotExist".
		if !os.IsNotExist(e) {
			m.l.Panicln(e)
		}
	}
}

func (m *Manager) processSave() {
	if *m.c.MemoryMode || m.file == nil || m.changes == false {
		return
	}
	m.changes = false

	e := file.SaveJSON(m.filePath(), m.file, os.FileMode(0600))
	if e != nil {
		m.l.Printf("Error: %v", e)
		return
	}
}

/*
	<<< LOOPS >>>
*/

func (m *Manager) service() {
	m.wg.Add(1)
	defer m.wg.Done()

	ticker := time.NewTicker(saveInternal)
	defer ticker.Stop()

	for {
		select {
		case req := <-m.requests:
			m.processRequest(req)

		case <-ticker.C:
			m.processSave()

		case <-m.quit:
			m.processSave()
			return
		}
	}
}

func (m *Manager) retryLoop() {
	m.wg.Add(1)
	defer m.wg.Done()

	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	retryQueue := new(object.RetryIO)

	for {
		select {
		case io := <-m.retries:
			retryQueue.Add(m.cxo.Initialize(io))

		case done := <-m.clearRetries:
			retryQueue = new(object.RetryIO)
			done <- struct{}{}

		case <-ticker.C:
			if retryQueue.IsEmpty() {
				continue
			}
			retryQueue = m.cxo.Initialize(retryQueue)

		case <-m.quit:
			return
		}
	}
}

/*
	<<< REQUESTS : PROCESS >>>
*/

type output struct {
	f *File
	e error
}

func (o *output) get() (*File, error) { return o.f, o.e }

type outputChan chan output

func (m *Manager) processRequest(r interface{}) {
	switch r.(type) {

	case *reqGetInfo:
		m.processGetInfo(r.(*reqGetInfo))

	case *reqNewConnection:
		m.processNewConnection(r.(*reqNewConnection))

	case *reqDeleteConnection:
		m.processDeleteConnection(r.(*reqDeleteConnection))

	case *reqNewSubscription:
		m.processNewSubscription(r.(*reqNewSubscription))

	case *reqDeleteSubscription:
		m.processDeleteSubscription(r.(*reqDeleteSubscription))

	case *reqNewMaster:
		m.processNewMaster(r.(*reqNewMaster))

	case *reqDeleteMaster:
		m.processDeleteMaster(r.(*reqDeleteMaster))

	default:
		m.l.Printf("Unprocessed request of type '%T'", r)
	}
}

/*
	<<< REQUESTS : GET INFO >>>
*/

type reqGetInfo struct {
	out outputChan
}

func (m *Manager) processGetInfo(r *reqGetInfo) {
	if m.file == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}
	r.out <- output{f: &(*m.file)}
}

// GetInfo obtains session information.
func (m *Manager) GetInfo(ctx context.Context) (*File, error) {
	ctx = m.timedContext(ctx)
	req := reqGetInfo{out: make(outputChan)}
	go func() { m.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.get()
	}
}

/*
	<<< REQUESTS : NEW CONNECTION >>>
*/

type reqNewConnection struct {
	in  *object.ConnectionIO
	out outputChan
}

func (m *Manager) processNewConnection(r *reqNewConnection) {
	if m.file == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}

	for _, address := range m.file.Connections {
		if address == r.in.Address {
			r.out <- output{
				e: boo.Newf(boo.AlreadyExists, "already connected to %s", address)}
			return
		}
	}

	m.retries <- &object.RetryIO{Addresses: []string{r.in.Address}}

	m.file.Connections = append(m.file.Connections, r.in.Address)
	m.changes = true
	r.out <- output{f: &(*m.file)}
}

func (m *Manager) NewConnection(ctx context.Context, in *object.ConnectionIO) (*File, error) {
	ctx = m.timedContext(ctx)
	req := reqNewConnection{in: in, out: make(outputChan)}
	go func() { m.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.get()
	}
}

/*
	<<< REQUESTS : DELETE CONNECTION >>>
*/

type reqDeleteConnection struct {
	in  *object.ConnectionIO
	out outputChan
}

func (m *Manager) processDeleteConnection(r *reqDeleteConnection) {
	if m.file == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}

	for i, address := range m.file.Connections {
		if address == r.in.Address {
			m.file.Connections = append(
				m.file.Connections[:i],
				m.file.Connections[i+1:]...,
			)
			m.changes = true
		}
	}

	cleared := make(chan struct{})
	m.clearRetries <- cleared
	<-cleared
	m.cxo.Disconnect(r.in.Address)
	m.retries <- m.cxo.Initialize(new(object.RetryIO).Fill(
		append(m.file.Subscriptions, m.file.Masters...),
		m.file.Connections,
	))

	r.out <- output{f: &(*m.file)}
}

func (m *Manager) DeleteConnection(ctx context.Context, in *object.ConnectionIO) (*File, error) {
	ctx = m.timedContext(ctx)
	req := reqDeleteConnection{in: in, out: make(outputChan)}
	go func() { m.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.get()
	}
}

/*
	<<< REQUESTS : NEW SUBSCRIPTION >>>
*/

type reqNewSubscription struct {
	in  *object.BoardIO
	out outputChan
}

func (m *Manager) processNewSubscription(r *reqNewSubscription) {
	if m.file == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}

	for _, sub := range m.file.Subscriptions {
		if sub.PubKey == r.in.PubKey {
			r.out <- output{
				e: boo.Newf(boo.AlreadyExists, "already subscribed to board %s", r.in.PubKeyStr)}
			return
		}
	}

	m.retries <- &object.RetryIO{PublicKeys: []cipher.PubKey{r.in.PubKey}}

	m.file.Subscriptions = append(m.file.Subscriptions, object.Subscription{PubKey: r.in.PubKey})
	m.changes = true
	r.out <- output{f: &(*m.file)}
}

// NewSub creates a new subscription.
func (m *Manager) NewSubscription(ctx context.Context, in *object.BoardIO) (*File, error) {
	ctx = m.timedContext(ctx)
	req := reqNewSubscription{in: in, out: make(outputChan)}
	go func() { m.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.get()
	}
}

/*
	<<< REQUESTS : DELETE SUBSCRIPTION >>>
*/

type reqDeleteSubscription struct {
	in  *object.BoardIO
	out outputChan
}

func (m *Manager) processDeleteSubscription(r *reqDeleteSubscription) {
	if m.file == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}

	for i, sub := range m.file.Subscriptions {
		if sub.PubKey == r.in.PubKey {
			m.file.Subscriptions = append(
				m.file.Subscriptions[:i],
				m.file.Subscriptions[i+1:]...,
			)
			m.changes = true
		}
	}

	cleared := make(chan struct{})
	m.clearRetries <- cleared
	<-cleared
	m.cxo.Unsubscribe("", r.in.PubKey)
	for _, address := range m.file.Connections {
		m.cxo.Unsubscribe(address, r.in.PubKey)
	}
	m.retries <- m.cxo.Initialize(new(object.RetryIO).Fill(
		append(m.file.Subscriptions, m.file.Masters...),
		m.file.Connections,
	))

	r.out <- output{f: &(*m.file)}
}

// DeleteSub removes a subscription.
func (m *Manager) DeleteSubscription(ctx context.Context, in *object.BoardIO) (*File, error) {
	ctx = m.timedContext(ctx)
	req := reqDeleteSubscription{in: in, out: make(outputChan)}
	go func() { m.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.get()
	}
}

/*
	<<< REQUESTS : NEW MASTER SUBSCRIPTION >>>
*/

type reqNewMaster struct {
	in  *object.NewBoardIO
	out outputChan
}

func (m *Manager) processNewMaster(r *reqNewMaster) {
	if m.file == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}
	for _, sub := range m.file.Masters {
		if sub.PubKey == r.in.BoardPubKey {
			r.out <- output{e: boo.Newf(boo.AlreadyExists,
				"you already own a board with public key %s", r.in.BoardPubKey.Hex())}
			return
		}
	}
	m.file.Masters = append(m.file.Masters, object.Subscription{
		PubKey: r.in.BoardPubKey,
		SecKey: r.in.BoardSecKey,
	})

	m.changes = true
	r.out <- output{f: &(*m.file)}
}

// NewMaster creates a new master subscription.
func (m *Manager) NewMaster(ctx context.Context, in *object.NewBoardIO) (*File, error) {
	ctx = m.timedContext(ctx)
	req := reqNewMaster{in: in, out: make(outputChan)}
	go func() { m.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.get()
	}
}

/*
	<<< REQUESTS : DELETE MASTER SUBSCRIPTION >>>
*/

type reqDeleteMaster struct {
	in  *object.BoardIO
	out outputChan
}

func (m *Manager) processDeleteMaster(r *reqDeleteMaster) {
	if m.file == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}
	for i, sub := range m.file.Masters {
		if sub.PubKey == r.in.PubKey {
			r.in.SecKey = sub.SecKey // Get secret key.
			m.file.Masters = append(
				m.file.Masters[:i],
				m.file.Masters[i+1:]...,
			)
			m.changes = true
			r.out <- output{f: &(*m.file)}
			return
		}
	}
	r.out <- output{e: boo.Newf(boo.NotFound,
		"master board of public key %s not found", r.in.PubKeyStr)}
}

// DeleteMaster removes a master subscription.
func (m *Manager) DeleteMaster(ctx context.Context, in *object.BoardIO) (*File, error) {
	ctx = m.timedContext(ctx)
	req := reqDeleteMaster{in: in, out: make(outputChan)}
	go func() { m.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.get()
	}
}

// TODO:
// - NewConnection
// - DeleteConnection
