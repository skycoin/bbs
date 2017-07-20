package state

import (
	"context"
	"github.com/skycoin/bbs/src/boo"
	"github.com/skycoin/bbs/src/misc/hide"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store/obj"
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
	logPrefix    = "USERSTATE"
	subDir       = "/users"
	extension    = ".usr"
	timeout      = time.Second * 10
	saveDuration = time.Second * 10
)

var (
	ErrNotLoggedIn     = boo.New(boo.NotAuthorised, "not logged in")
	ErrAlreadyLoggedIn = boo.New(boo.NotAuthorised, "already logged in")
)

// SessionManagerConfig configures a SessionManager.
type SessionManagerConfig struct {
	ConfigDir *string `json:"config_dir"`
	Memory    *bool   `json:"memory"`
}

// SessionManager manages user sessions.
type SessionManager struct {
	// Configuration.
	c *SessionManagerConfig
	l *log.Logger

	// Session.
	user    *UserFile
	key     []byte
	changes bool

	// Request / response.
	requests chan interface{}
	quit     chan struct{}
	wg       sync.WaitGroup
}

// NewSessionManager creates a new SessionManager with configuration.
func NewSessionManager(config *SessionManagerConfig) (*SessionManager, error) {
	s := &SessionManager{
		c:        config,
		l:        inform.NewLogger(true, os.Stdout, logPrefix),
		changes:  false,
		requests: make(chan interface{}),
		quit:     make(chan struct{}),
	}
	e := os.MkdirAll(s.folderPath(), os.FileMode(0700))
	if e != nil {
		return nil, e
	}
	go s.service()
	return s, nil
}

// Close closes the SessionManager service.
func (s *SessionManager) Close() {
	select {
	case s.quit <- struct{}{}:
	default:
	}
	s.wg.Wait()
}

func (s *SessionManager) folderPath() string {
	return path.Join(*s.c.ConfigDir, subDir)
}

func (s *SessionManager) filePath(user string) string {
	return path.Join(s.folderPath(), user+extension)
}

func (s *SessionManager) timedContext(ctx context.Context) context.Context {
	ctx, _ = context.WithTimeout(ctx, timeout)
	return ctx
}

func (s *SessionManager) service() {
	s.wg.Add(1)
	defer s.wg.Done()

	ticker := time.NewTicker(saveDuration)
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

/*
	<<< TRIGGERS : SAVE >>>
*/

func (s *SessionManager) processSave() {
	s.changes = false
	if *s.c.Memory || s.user == nil {
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

func (o *output) get() (*UserFile, error) {return o.f, o.e}

type outputChan chan output

func (s *SessionManager) processRequest(r interface{}) {
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

func (s *SessionManager) processGetUsers(r *reqGetUsers) {
	files, e := ioutil.ReadDir(s.folderPath())
	if e != nil {
		r.e <- e
		return
	}
	var users []string
	for _, info := range files {
		if info.IsDir() || !strings.HasSuffix(info.Name(), extension) {
			continue
		}
		name := strings.TrimSuffix(info.Name(), extension)
		s.l.Printf("Found User: '%s'.", name)
		users = append(users, name)
	}
	r.out <- users
}

// GetUsers obtains list of available out.
func (s *SessionManager) GetUsers(ctx context.Context) ([]string, error) {
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
	in  *NewUserInput
	out outputChan
}

func (s *SessionManager) processNewUser(r *reqNewUser) {
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
func (s *SessionManager) NewUser(ctx context.Context, in *NewUserInput) (*UserFile, error) {
	ctx = s.timedContext(ctx)
	// TODO: Check inputs.

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

func (s *SessionManager) processDeleteUser(r *reqDeleteUser) {
	if s.user != nil {
		r.e <- ErrAlreadyLoggedIn
		return
	}
	r.e <- os.Remove(s.filePath(r.alias))
}

// DeleteUser deletes a user configuration user.
func (s *SessionManager) DeleteUser(ctx context.Context, alias string) error {
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
	in  *LoginInput
	out outputChan
}

func (s *SessionManager) processLogin(r *reqLogin) {
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

	s.user, s.key, s.changes = uFile, key, true
	r.out <- output{f: &(*uFile)}
}

// Login logs in a user.
func (s *SessionManager) Login(ctx context.Context, in *LoginInput) (*UserFile, error) {
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

func (s *SessionManager) processLogout(r *reqLogout) {
	s.processSave()
	s.user, s.key = nil, nil
	r.e <- nil
}

func (s *SessionManager) Logout(ctx context.Context) error {
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
	<<< REQUESTS : GET INFO >>> TODO: TEST
*/

type reqGetInfo struct {
	out outputChan
}

func (s *SessionManager) processGetInfo(r *reqGetInfo) {
	if s.user == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}
	r.out <- output{f: &(*s.user)}
}

// GetInfo obtains session information.
func (s *SessionManager) GetInfo(ctx context.Context) (*UserFile, error) {
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
	<<< REQUESTS : New SUBSCRIPTION >>> TODO: TEST
*/

type reqNewSubscription struct {
	in  *SubscriptionInput
	out outputChan
}

func (s *SessionManager) processNewSubscription(r *reqNewSubscription) {
	if s.user == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}

	if e := r.in.process(); e != nil {
		r.out <- output{e: e}
		return
	}

	for _, sub := range s.user.Subscriptions {
		if sub.PubKey == r.in.pubKey {
			r.out <- output{
				e: boo.Newf(boo.ObjectAlreadyExists, "already subscribed to board %s", r.in.pubKey.Hex())}
			return
		}
	}

	s.user.Subscriptions = append(s.user.Subscriptions, obj.Subscription{PubKey: r.in.pubKey})
	s.changes = true
	r.out <- output{f: &(*s.user)}
}

// NewSubscription creates a new subscription.
func (s *SessionManager) NewSubscription(ctx context.Context, in *SubscriptionInput) (*UserFile, error) {
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
	<<< REQUESTS : DELETE SUBSCRIPTION >>> TODO: TEST
*/

type reqDeleteSubscription struct {
	in  *SubscriptionInput
	out outputChan
}

func (s *SessionManager) processDeleteSubscription(r *reqDeleteSubscription) {
	if s.user == nil {
		r.out <- output{e: ErrNotLoggedIn}
		return
	}

	if e := r.in.process(); e != nil {
		r.out <- output{e: e}
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

	r.out <- output{f: &(*s.user)}
}

// DeleteSubscription removes a subscription.
func (s *SessionManager) DeleteSubscription(ctx context.Context, in *SubscriptionInput) (*UserFile, error) {
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
	<<< REQUESTS : NEW MASTER SUBSCRIPTION >>> TODO: TEST
*/

type reqNewMaster struct {
	in  *NewMasterInput
	out outputChan
}

func (s *SessionManager) processNewMaster(r *reqNewMaster) {
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
		Connections: r.in.connections,
	})
	s.changes = true
	r.out <- output{f: &(*s.user)}
}

// NewMaster creates a new master subscription.
func (s *SessionManager) NewMaster(ctx context.Context, in *NewMasterInput) (*UserFile, error) {
	if e := in.process(); e != nil {
		return nil, e
	}
	ctx = s.timedContext(ctx)
	req := reqNewMaster{in: in, out: make(outputChan)}
	go func() { s.requests <- &req }()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <- req.out:
		return out.get()
	}
}

/*
	<<< REQUESTS : DELETE MASTER SUBSCRIPTION >>> TODO: TEST
*/

type reqDeleteMaster struct {
	in *SubscriptionInput
	out outputChan
}

func (s *SessionManager) processDeleteMaster(r *reqDeleteMaster) {
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
func (s *SessionManager) DeleteMaster(ctx context.Context, in *SubscriptionInput) (*UserFile, error) {
	if e := in.process(); e != nil {
		return nil, e
	}
	ctx = s.timedContext(ctx)
	req := reqDeleteMaster{in: in, out: make(outputChan)}
	go func() { s.requests <- &req } ()
	select {
	case <-ctx.Done():
		return nil, boo.WrapType(ctx.Err(), boo.Internal)
	case out := <-req.out:
		return out.get()
	}
}

// TODO:
// - AddMasterConnection
// - RemoveMasterConnection
// - AddMasterSubmissionAddress
// - RemoveMasterSubmissionAddress

