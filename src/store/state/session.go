package state

import (
	"context"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/hide"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/obj"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"github.com/skycoin/skycoin/src/util/file"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"
)

const (
	usersLogPrefix     = "USERSTATE"
	usersSubDir        = "/users"
	usersFileExtension = ".usr"
	usersTimeout       = time.Second * 10
	usersSaveDuration  = time.Second * 10
	usersRetryDuration = time.Second * 5
)

var (
	ErrNotLoggedIn     = boo.New(boo.NotAuthorised, "not logged in")
	ErrAlreadyLoggedIn = boo.New(boo.NotAuthorised, "already logged in")
)

// SessionConfig configures a Session.
type SessionConfig struct {
	Master       *bool   // Whether node is master.
	TestMode     *bool   // Whether node is in test mode.
	MemoryMode   *bool   // Whether to use local storage in runtime.
	ConfigDir    *string // Configuration directory.
	CXOPort      *int    // CXO listening port.
	CXORPCEnable *bool   // Whether to enable CXO RPC.
	CXORPCPort   *int    // CXO RPC listening port.
}

// Session manages user sessions.
type Session struct {
	// Configuration.
	c *SessionConfig
	l *log.Logger

	// Session.
	cxo     *CXO
	user    *UserFile
	key     []byte
	changes bool

	// Request & Response.
	retries      chan *RetryIO      // Channel for retrying failed subscriptions/connections.
	clearRetries chan chan struct{} // Channel for clearing retry queue (i.e. logout).
	requests     chan interface{}   // Channel for user requests.
	quit         chan struct{}      // Channel for quitting Session.
	wg           sync.WaitGroup
}

// NewSession creates a new Session with configuration.
func NewSession(config *SessionConfig) (*Session, error) {
	s := &Session{
		c:            config,
		l:            inform.NewLogger(true, os.Stdout, usersLogPrefix),
		cxo:          NewCXO(config, func(_ *node.Root) {}), // TODO: Implement.
		changes:      false,
		retries:      make(chan *RetryIO, 5),
		clearRetries: make(chan chan struct{}),
		requests:     make(chan interface{}, 5),
		quit:         make(chan struct{}),
	}
	if e := os.MkdirAll(s.folderPath(), os.FileMode(0700)); e != nil {
		return nil, e
	}
	go s.service()
	go s.retryLoop()
	return s, nil
}

// Close closes the Session service.
func (s *Session) Close() {
	for {
		select {
		case s.quit <- struct{}{}:
			continue
		default:
			s.wg.Wait()
			if e := s.cxo.Close(); e != nil && e != ErrCXONotOpened {
				s.l.Println("Closing CXO Error:", e)
			}
			return
		}
	}
}

// GetCXO obtains the Inside CXO.
func (s *Session) GetCXO() *CXO { return s.cxo }

/*
	<<< HELPER FUNCTIONS >>>
*/

func (s *Session) folderPath() string {
	return path.Join(*s.c.ConfigDir, usersSubDir)
}

func (s *Session) filePath(user string) string {
	return path.Join(s.folderPath(), user+usersFileExtension)
}

func (s *Session) timedContext(ctx context.Context) context.Context {
	ctx, _ = context.WithTimeout(ctx, usersTimeout)
	return ctx
}

/*
	<<< LOOPS >>>
*/

func (s *Session) service() {
	s.wg.Add(1)
	defer s.wg.Done()

	ticker := time.NewTicker(usersSaveDuration)
	defer ticker.Stop()

	for {
		select {
		case req := <-s.requests:
			s.processRequest(req)

		case <-ticker.C:
			s.processSave()

		case <-s.quit:
			s.processSave()
			return
		}
	}
}

func (s *Session) retryLoop() {
	s.wg.Add(1)
	defer s.wg.Done()

	ticker := time.NewTicker(usersRetryDuration)
	defer ticker.Stop()

	retryQueue := new(RetryIO)

	for {
		select {
		case io := <-s.retries:
			retryQueue.Add(s.cxo.Initialize(io))

		case done := <-s.clearRetries:
			retryQueue = new(RetryIO)
			done <- struct{}{}

		case <-ticker.C:
			if retryQueue.IsEmpty() {
				continue
			}
			retryQueue = s.cxo.Initialize(retryQueue)

		case <-s.quit:
			return
		}
	}
}

/*
	<<< TRIGGERS : SAVE >>>
*/

func (s *Session) processSave() {
	s.changes = false
	if *s.c.MemoryMode || s.user == nil {
		return
	}
	data, e := hide.Encrypt(s.key, encoder.Serialize(*s.user))
	if e != nil {
		s.l.Printf("Error: %v", e)
		return
	}
	e = file.SaveBinary(
		s.filePath(s.user.User.Alias),
		data, os.FileMode(0600),
	)
	if e != nil {
		s.l.Printf("Error: %v", e)
		return
	}
}

/*
	<<< REQUESTS : PROCESS >>>
*/

type output struct {
	f *UserFile
	e error
}

func (o *output) get() (*UserFile, error) { return o.f, o.e }

type outputChan chan output

func (s *Session) processRequest(r interface{}) {
	switch r.(type) {
	case *reqGetUsers:
		s.processGetUsers(r.(*reqGetUsers))

	case *reqNewUser:
		s.processNewUser(r.(*reqNewUser))

	case *reqDeleteUser:
		s.processDeleteUser(r.(*reqDeleteUser))

	case *reqLogin:
		s.processLogin(r.(*reqLogin))

	case *reqLogout:
		s.processLogout(r.(*reqLogout))

	case *reqGetInfo:
		s.processGetInfo(r.(*reqGetInfo))

	case *reqNewConnection:
		s.processNewConnection(r.(*reqNewConnection))

	case *reqDeleteConnection:
		s.processDeleteConnection(r.(*reqDeleteConnection))

	case *reqNewSubscription:
		s.processNewSubscription(r.(*reqNewSubscription))

	case *reqDeleteSubscription:
		s.processDeleteSubscription(r.(*reqDeleteSubscription))

	case *reqNewMaster:
		s.processNewMaster(r.(*reqNewMaster))

	case *reqDeleteMaster:
		s.processDeleteMaster(r.(*reqDeleteMaster))

	default:
		s.l.Printf("Unprocessed request of type '%T'", r)
	}
}

/*
	<<< REQUESTS : GET USERS >>>
*/

type reqGetUsers struct {
	out chan []string
	e   chan error
}

func (s *Session) processGetUsers(r *reqGetUsers) {
	files, e := ioutil.ReadDir(s.folderPath())
	if e != nil {
		r.e <- e
		return
	}
	var users []string
	for _, info := range files {
		if info.IsDir() || !strings.HasSuffix(info.Name(), usersFileExtension) {
			continue
		}
		name := strings.TrimSuffix(info.Name(), usersFileExtension)
		s.l.Printf("Found User: '%s'.", name)
		users = append(users, name)
	}
	r.out <- users
}

// GetUsers obtains list of available out.
func (s *Session) GetUsers(ctx context.Context) ([]string, error) {
	ctx = s.timedContext(ctx)
	req := reqGetUsers{out: make(chan []string), e: make(chan error)}
	go func() { s.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case users := <-req.out:
		return users, nil
	case e := <-req.e:
		return nil, e
	}
}

/*
	<<< REQUESTS : NEW USER >>>
*/

type reqNewUser struct {
	in  *NewUserIO
	out outputChan
}

func (s *Session) processNewUser(r *reqNewUser) {
	pk, sk := cipher.GenerateDeterministicKeyPair([]byte(r.in.Seed))
	key := []byte(r.in.Password)
	uFile := UserFile{
		User: obj.User{
			Alias:     r.in.Alias,
			PublicKey: pk,
			SecretKey: sk,
		},
	}
	data, e := hide.Encrypt(key, encoder.Serialize(uFile))
	if e != nil {
		r.out <- output{e: e}
		return
	}
	e = file.SaveBinary(s.filePath(r.in.Alias), data, os.FileMode(0600))
	if e != nil {
		r.out <- output{e: e}
		return
	}
	r.out <- output{f: &uFile}
}

// NewUser creates a new user.
func (s *Session) NewUser(ctx context.Context, in *NewUserIO) (*UserFile, error) {
	ctx = s.timedContext(ctx)
	req := reqNewUser{in: in, out: make(outputChan)}
	go func() { s.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.get()
	}
}

/*
	<<< REQUESTS : DELETE USER >>>
*/

type reqDeleteUser struct {
	alias string
	e     chan error
}

func (s *Session) processDeleteUser(r *reqDeleteUser) {
	if s.user != nil {
		r.e <- ErrAlreadyLoggedIn
		return
	}
	r.e <- os.Remove(s.filePath(r.alias))
}

// DeleteUser deletes a user configuration user.
func (s *Session) DeleteUser(ctx context.Context, alias string) error {
	ctx = s.timedContext(ctx)
	req := reqDeleteUser{alias: alias, e: make(chan error)}
	go func() { s.requests <- &req }()
	select {
	case <-ctx.Done():
		return boo.WrapType(ctx.Err(), boo.Internal)
	case e := <-req.e:
		return e
	}
}

/*
	<<< REQUESTS : LOGIN >>>
*/

type reqLogin struct {
	in  *LoginIO
	out outputChan
}

func (s *Session) processLogin(r *reqLogin) {
	if s.user != nil {
		r.out <- output{e: ErrAlreadyLoggedIn}
		return
	}

	fName := s.filePath(r.in.Alias)
	key := []byte(r.in.Password)

	data, e := ioutil.ReadFile(fName)
	if e != nil {
		time.Sleep(5 * time.Second)
		r.out <- output{
			e: boo.WrapType(e, boo.ObjectNotFound, "user user not found")}
		return
	}

	data, e = hide.Decrypt(key, data)
	if e != nil {
		time.Sleep(5 * time.Second)
		r.out <- output{
			e: boo.WrapType(e, boo.NotAuthorised, "decryption failed")}
		return
	}

	uFile := &UserFile{}
	if e := encoder.DeserializeRaw(data, uFile); e != nil {
		time.Sleep(5 * time.Second)
		r.out <- output{
			e: boo.WrapType(e, boo.InvalidRead, "corrupt user user")}
		return
	}

	if failed, e := s.cxo.Open(r.in.Alias, new(RetryIO).Fill(uFile)); e != nil {
		r.out <- output{
			e: boo.WrapType(e, boo.Internal, "failed to open cxo")}
		return
	} else {
		if !failed.IsEmpty() {
			s.retries <- failed
		}
	}

	s.user, s.key, s.changes = uFile, key, true
	r.out <- output{f: &(*uFile)}
}

// Login logs in a user.
func (s *Session) Login(ctx context.Context, in *LoginIO) (*UserFile, error) {
	ctx = s.timedContext(ctx)
	req := reqLogin{in: in, out: make(outputChan)}
	go func() { s.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.get()
	}
}

/*
	<<< REQUESTS : LOGOUT >>>
*/

type reqLogout struct {
	e chan error
}

func (s *Session) processLogout(r *reqLogout) {
	if s.user == nil {
		r.e <- ErrNotLoggedIn
		return
	}
	s.processSave()

	// Clear retries.
	done := make(chan struct{})
	go func() { s.clearRetries <- done }()
	<-done

	// Close CXO.
	if e := s.cxo.Close(); e != nil {
		r.e <- boo.WrapType(e, boo.Internal, "failed to close cxo")
		return
	}
	s.user, s.key = nil, nil
	r.e <- nil
}

func (s *Session) Logout(ctx context.Context) error {
	ctx = s.timedContext(ctx)
	req := reqLogout{e: make(chan error)}
	go func() { s.requests <- &req }()
	select {
	case <-ctx.Done():
		return boo.WrapType(ctx.Err(), boo.Internal)
	case e := <-req.e:
		return e
	}
	return nil
}

/*
	<<< REQUESTS : GET INFO >>>
*/

type reqGetInfo struct {
	out outputChan
}

func (s *Session) processGetInfo(r *reqGetInfo) {
	if s.user == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}
	r.out <- output{f: &(*s.user)}
}

// GetInfo obtains session information.
func (s *Session) GetInfo(ctx context.Context) (*UserFile, error) {
	ctx = s.timedContext(ctx)
	req := reqGetInfo{out: make(outputChan)}
	go func() { s.requests <- &req }()
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
	in  *ConnectionIO
	out outputChan
}

func (s *Session) processNewConnection(r *reqNewConnection) {
	if s.user == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}

	for _, address := range s.user.Connections {
		if address == r.in.Address {
			r.out <- output{
				e: boo.Newf(boo.ObjectAlreadyExists, "already connected to %s", address)}
			return
		}
	}

	s.retries <- &RetryIO{addresses: []string{r.in.Address}}

	s.user.Connections = append(s.user.Connections, r.in.Address)
	s.changes = true
	r.out <- output{f: &(*s.user)}
}

func (s *Session) NewConnection(ctx context.Context, in *ConnectionIO) (*UserFile, error) {
	ctx = s.timedContext(ctx)
	req := reqNewConnection{in: in, out: make(outputChan)}
	go func() { s.requests <- &req }()
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
	in  *ConnectionIO
	out outputChan
}

func (s *Session) processDeleteConnection(r *reqDeleteConnection) {
	if s.user == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}

	for i, address := range s.user.Connections {
		if address == r.in.Address {
			s.user.Connections = append(
				s.user.Connections[:i],
				s.user.Connections[i+1:]...,
			)
			s.changes = true
		}
	}

	cleared := make(chan struct{})
	s.clearRetries <- cleared
	<-cleared
	s.cxo.Disconnect(r.in.Address)
	s.retries <- s.cxo.Initialize(new(RetryIO).Fill(s.user))

	r.out <- output{f: &(*s.user)}
}

func (s *Session) DeleteConnection(ctx context.Context, in *ConnectionIO) (*UserFile, error) {
	ctx = s.timedContext(ctx)
	req := reqDeleteConnection{in: in, out: make(outputChan)}
	go func() { s.requests <- &req }()
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
	in  *SubscriptionIO
	out outputChan
}

func (s *Session) processNewSubscription(r *reqNewSubscription) {
	if s.user == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}

	for _, sub := range s.user.Subscriptions {
		if sub.PubKey == r.in.pubKey {
			r.out <- output{
				e: boo.Newf(boo.ObjectAlreadyExists, "already subscribed to board %s", r.in.pubKey.Hex())}
			return
		}
	}

	s.retries <- &RetryIO{pks: []cipher.PubKey{r.in.pubKey}}

	s.user.Subscriptions = append(s.user.Subscriptions, obj.Subscription{PubKey: r.in.pubKey})
	s.changes = true
	r.out <- output{f: &(*s.user)}
}

// NewSub creates a new subscription.
func (s *Session) NewSubscription(ctx context.Context, in *SubscriptionIO) (*UserFile, error) {
	ctx = s.timedContext(ctx)
	req := reqNewSubscription{in: in, out: make(outputChan)}
	go func() { s.requests <- &req }()
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
	in  *SubscriptionIO
	out outputChan
}

func (s *Session) processDeleteSubscription(r *reqDeleteSubscription) {
	if s.user == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}

	for i, sub := range s.user.Subscriptions {
		if sub.PubKey == r.in.pubKey {
			s.user.Subscriptions = append(
				s.user.Subscriptions[:i],
				s.user.Subscriptions[i+1:]...,
			)
			s.changes = true
		}
	}

	cleared := make(chan struct{})
	s.clearRetries <- cleared
	<-cleared
	s.cxo.Unsubscribe("", r.in.pubKey)
	for _, address := range s.user.Connections {
		s.cxo.Unsubscribe(address, r.in.pubKey)
	}
	s.retries <- s.cxo.Initialize(new(RetryIO).Fill(s.user))

	r.out <- output{f: &(*s.user)}
}

// DeleteSub removes a subscription.
func (s *Session) DeleteSubscription(ctx context.Context, in *SubscriptionIO) (*UserFile, error) {
	ctx = s.timedContext(ctx)
	req := reqDeleteSubscription{in: in, out: make(outputChan)}
	go func() { s.requests <- &req }()
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
	in  *NewMasterIO
	out outputChan
}

func (s *Session) processNewMaster(r *reqNewMaster) {
	if s.user == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}
	for _, sub := range s.user.Masters {
		if sub.PubKey == r.in.pubKey {
			r.out <- output{e: boo.Newf(boo.ObjectAlreadyExists,
				"you already own a board with public key %s", r.in.pubKey.Hex())}
			return
		}
	}
	s.user.Masters = append(s.user.Masters, obj.Subscription{
		PubKey: r.in.pubKey,
		SecKey: r.in.secKey,
	})

	s.changes = true
	r.out <- output{f: &(*s.user)}
}

// NewMaster creates a new master subscription.
func (s *Session) NewMaster(ctx context.Context, in *NewMasterIO) (*UserFile, error) {
	ctx = s.timedContext(ctx)
	req := reqNewMaster{in: in, out: make(outputChan)}
	go func() { s.requests <- &req }()
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
	in  *SubscriptionIO
	out outputChan
}

func (s *Session) processDeleteMaster(r *reqDeleteMaster) {
	if s.user == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}
	for i, sub := range s.user.Masters {
		if sub.PubKey == r.in.pubKey {
			s.user.Masters = append(
				s.user.Masters[:i],
				s.user.Masters[i+1:]...,
			)
			s.changes = true
		}
	}
	r.out <- output{f: &(*s.user)}
}

// DeleteMaster removes a master subscription.
func (s *Session) DeleteMaster(ctx context.Context, in *SubscriptionIO) (*UserFile, error) {
	ctx = s.timedContext(ctx)
	req := reqDeleteMaster{in: in, out: make(outputChan)}
	go func() { s.requests <- &req }()
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

// - AddMasterSubmissionAddress
// - RemoveMasterSubmissionAddress
