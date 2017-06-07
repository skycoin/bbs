package cxo

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/cmd/bbsnode/args"
	"github.com/skycoin/bbs/intern/typ"
	"github.com/skycoin/cxo/node"
	"github.com/skycoin/cxo/node/gnet"
	"github.com/skycoin/cxo/skyobject"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

type Container struct {
	sync.Mutex
	c      *node.Container
	client *node.Client
	config *args.Config
	msgs   chan *Msg

	cbAddFeed func(key cipher.PubKey) (bool, error)
	cbDelFeed func(key cipher.PubKey) (bool, error)
	cbClose   func() error

	tempTestFile string

	quit chan struct{}
}

func NewContainer(config *args.Config) (*Container, error) {
	c := &Container{
		config: config,
		msgs:   make(chan *Msg),
	}

	// Prepare waiter, error channel and timeout.
	var wg sync.WaitGroup
	eChan := make(chan error, 100)
	timeout := 10 * time.Second

	// run stuff.
	go c.setupCXOClient(wg, eChan, timeout)
	if config.CXOUseInternal() {
		go c.setupInternalCXODaemon(wg, eChan, timeout)
	} else {
		go c.setupCXORPCClient(wg, eChan, timeout)
	}

	// Wait.
	time.Sleep(time.Second)
	wg.Wait()

	// Check for errors.
	select {
	case e := <-eChan:
		if e != nil {
			return nil, e
		}
	default:
		break
	}

	log.Println("[CXOCONTAINER] Connection to cxo daemon established!")

	go c.service()
	return c, nil
}

func (c *Container) setupCXOClient(wg sync.WaitGroup, eChan chan error, timeout time.Duration) {
	wg.Add(1)
	defer wg.Done()
	timer := time.NewTimer(timeout)

	// Setup CXO Registry.
	r := skyobject.NewRegistry()
	r.Register("Board", typ.Board{})
	r.Register("Thread", typ.Thread{})
	r.Register("Post", typ.Post{})
	r.Register("ThreadPage", typ.ThreadPage{})
	r.Register("BoardContainer", typ.BoardContainer{})

	r.Register("Vote", typ.Vote{})
	r.Register("ThreadVotes", typ.ThreadVotes{})
	r.Register("PostVotes", typ.PostVotes{})
	r.Register("ThreadVotesContainer", typ.ThreadVotesContainer{})
	r.Register("PostVotesContainer", typ.PostVotesContainer{})
	r.Done()

	// Setup CXO Configuration.
	cc := node.NewClientConfig()
	cc.MaxMessageSize = 0 /* TODO: Should really look into adjusting this in the future. */
	cc.InMemoryDB = c.config.CXOUseMemory()
	cc.DataDir = c.config.CXOCDir()
	cc.DBPath = filepath.Join(cc.DataDir, "client.db")

	// Setup CXO Callbacks.
	if c.config.Master() {
		cc.OnRootFilled = c.rootFilledCallBack
	}
	cc.OnAddFeed = c.feedAddedInternalCB
	cc.OnDelFeed = c.feedDeletedInternalCB

	// Change some configurations if test mode.
	if c.config.TestMode() {

		// Make temp file.
		tempDir, e := ioutil.TempDir("", "skybbs")
		if e != nil {
			eChan <- errors.Wrap(e, "unable to create temp dir")
			return
		}
		tempFile, e := ioutil.TempFile(tempDir, "")
		if e != nil {
			eChan <- errors.Wrap(e, "unable to create temp file")
			return
		}
		tempFileName := tempFile.Name()
		tempFile.Close()

		// Change stuff.
		c.tempTestFile = tempFileName
		cc.DataDir = tempDir
		cc.DBPath = tempFileName
		cc.InMemoryDB = false
	}

	// Attempt to set up CXO Client.
	var e error
	for {
		c.client, e = node.NewClient(cc, r)
		if e == nil {
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.Wrap(e, "timeout before cxo client initiation")
			return
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	addr := "[::]:" + strconv.Itoa(c.config.CXOPort())
	for {
		e = c.client.Start(addr)
		if e == nil {
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.Wrap(e, "timeout before cxo client start")
			return
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	for {
		if c.client.Conn().State() == gnet.ConnStateConnected {
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.New("timeout before cxo client-daemon connection")
			return
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	// Setup internal container.
	c.c = c.client.Container()
	return
}

func (c *Container) setupCXORPCClient(wg sync.WaitGroup, eChan chan error, timeout time.Duration) {
	wg.Add(1)
	defer wg.Done()
	timer := time.NewTimer(timeout)

	// Attempt to set up CXO RPC.
	addr := "[::]:" + strconv.Itoa(c.config.CXORPCPort())
	for {
		rpc, e := node.NewRPCClient(addr)
		if e == nil {
			c.cbAddFeed = rpc.AddFeed
			c.cbDelFeed = rpc.DelFeed
			c.cbClose = rpc.Close
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.Wrap(e, "timeout before rpc connected")
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	return
}

func (c *Container) setupInternalCXODaemon(wg sync.WaitGroup, eChan chan error, timeout time.Duration) {
	wg.Add(1)
	defer wg.Done()
	timer := time.NewTimer(timeout)

	// Setup CXO Server Configuration.
	sc := node.NewServerConfig()
	sc.DataDir = c.config.CXOCDir()
	sc.DBPath = filepath.Join(sc.DataDir, "server.db")
	sc.InMemoryDB = c.config.CXOUseMemory()
	sc.MaxMessageSize = 0
	sc.Listen = "[::]:" + strconv.Itoa(c.config.CXOPort())
	sc.RPCAddress = "[::]:" + strconv.Itoa(c.config.CXORPCPort())

	// Change some configurations if test mode.
	if c.config.TestMode() {

		// Make temp file.
		tempDir, e := ioutil.TempDir("", "skybbs")
		if e != nil {
			eChan <- errors.Wrap(e, "unable to create temp dir")
			return
		}
		tempFile, e := ioutil.TempFile(tempDir, "")
		if e != nil {
			eChan <- errors.Wrap(e, "unable to create temp file")
			return
		}
		tempFileName := tempFile.Name()
		tempFile.Close()

		// Change stuff.
		sc.DataDir = tempDir
		sc.DBPath = tempFileName
		sc.InMemoryDB = false
	}

	// Attempt to run server.
	var cxoServer *node.Server
	var e error
	for {
		cxoServer, e = node.NewServer(sc)
		if e == nil {
			c.cbAddFeed = func(key cipher.PubKey) (bool, error) {
				return cxoServer.AddFeed(key), nil
			}
			c.cbDelFeed = func(key cipher.PubKey) (bool, error) {
				return cxoServer.DelFeed(key), nil
			}
			c.cbClose = cxoServer.Close
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.Wrap(e, "timeout before cxo server created")
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	for {
		e := cxoServer.Start()
		if e == nil {
			break
		}
		select {
		case <-timer.C:
			eChan <- errors.Wrap(e, "timeout before cxo server started")
		default:
			time.Sleep(time.Second)
			continue
		}
	}
	return
}

func (c *Container) Close() error {
	for {
		select {
		case c.quit <- struct{}{}:
		default:
			goto ServiceFinished
		}
	}
ServiceFinished:
	if c.config.TestMode() {
		defer c.deleteTemp(c.tempTestFile)
	}
	c.cbClose()
	return c.client.Close()
}

func (c *Container) deleteTemp(cxoDir string) {
	// Get home dir.
	if e := os.Remove(cxoDir); e != nil {
		fmt.Println("REMOVING (ERROR):", e.Error())
	}
}

func (c *Container) service() {
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-c.quit:
			return
		case <-ticker.C:
			c.Lock()
			c.c.GC(false)
			c.Unlock()
		}
	}
}

func (c *Container) Feeds() []cipher.PubKey {
	return c.client.Feeds()
}

func (c *Container) Subscribe(pk cipher.PubKey) (bool, error) {
	if _, e := c.cbAddFeed(pk); e != nil {
		return false, e
	}
	return c.client.Subscribe(pk), nil
}

func (c *Container) Unsubscribe(pk cipher.PubKey) (bool, error) {
	if _, e := c.cbDelFeed(pk); e != nil {
		return false, e
	}
	return c.client.Unsubscribe(pk, false), nil
}

// ChangeBoardURL changes the board's URL of given public key.
func (c *Container) ChangeBoardURL(bpk cipher.PubKey, bsk cipher.SecKey, url string) error {
	c.Lock()
	defer c.Unlock()

	r := c.c.LastRootSk(bpk, bsk)
	w := r.Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(r)); e != nil {
		return e
	}
	b := &typ.Board{}
	if e := w.AdvanceFromRefField("Board", b); e != nil {
		return e
	}
	b.URL = url
	w.Retreat()
	_, e := w.ReplaceInRefField("Board", *b)
	return e
}

// GetBoard attempts to obtain the board of a given public key.
func (c *Container) GetBoard(bpk cipher.PubKey) (*typ.Board, error) {
	c.Lock()
	defer c.Unlock()

	w := c.c.LastFullRoot(bpk).Walker()
	bc := typ.BoardContainer{}
	if e := w.AdvanceFromRoot(&bc, makeBoardContainerFinder(w.Root())); e != nil {
		return nil, e
	}
	b := &typ.Board{}
	e := w.AdvanceFromRefField("Board", b)
	return b, e
}

// GetBoards attempts to obtain a list of boards from the given public keys.
func (c *Container) GetBoards(bpks ...cipher.PubKey) []*typ.Board {
	c.Lock()
	defer c.Unlock()

	boards := make([]*typ.Board, len(bpks))
	for i, bpk := range bpks {
		w := c.c.LastRoot(bpk).Walker()
		bc, b := typ.BoardContainer{}, typ.Board{}
		if e := w.AdvanceFromRoot(&bc, makeBoardContainerFinder(w.Root())); e != nil {
			continue
		}
		w.AdvanceFromRefField("Board", &b)
		boards[i] = &b
	}
	return boards
}

// NewBoard attempts to create a new board from a given board and seed.
func (c *Container) NewBoard(board *typ.Board, pk cipher.PubKey, sk cipher.SecKey) error {
	c.Lock()
	defer c.Unlock()

	r, e := c.c.NewRoot(pk, sk)
	if e != nil {
		return e
	}
	bRef := r.Save(*board)
	// Prepare board container.
	bCont := typ.BoardContainer{Board: bRef}
	if _, _, e = r.Inject("BoardContainer", bCont); e != nil {
		return e
	}
	// Prepare thread vote container.
	tvCont := typ.ThreadVotesContainer{}
	if _, _, e := r.Inject("ThreadVotesContainer", tvCont); e != nil {
		return e
	}
	// Prepare post vote container.
	pvCont := typ.PostVotesContainer{}
	if _, _, e := r.Inject("PostVotesContainer", pvCont); e != nil {
		return e
	}
	return nil
}

// RemoveBoard attempts to remove a board by a given public key.
func (c *Container) RemoveBoard(bpk cipher.PubKey, bsk cipher.SecKey) error {
	c.Lock()
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	fmt.Println("Removing board:", bpk.Hex())
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	return w.RemoveCurrent()
}

// GetThread obtains a single thread via reference.
func (c *Container) GetThread(tRef skyobject.Reference) (*typ.Thread, error) {
	c.Lock()
	defer c.Unlock()

	tData, has := c.c.Get(tRef)
	if !has {
		return nil, errors.New("thread not found")
	}
	thread := &typ.Thread{}
	if e := encoder.DeserializeRaw(tData, thread); e != nil {
		return nil, e
	}
	thread.Ref = tRef.String()
	return thread, nil
}

// GetThreads attempts to obtain a list of threads from a board of public key.
func (c *Container) GetThreads(bpk cipher.PubKey) ([]*typ.Thread, error) {
	c.Lock()
	defer c.Unlock()

	w := c.c.LastFullRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return nil, e
	}
	threads := make([]*typ.Thread, len(bc.Threads))
	for i, tRef := range bc.Threads {
		tData, has := c.c.Get(tRef)
		if has == false {
			continue
		}
		threads[i] = new(typ.Thread)
		if e := encoder.DeserializeRaw(tData, threads[i]); e != nil {
			return nil, e
		}
		threads[i].Ref = cipher.SHA256(tRef).Hex()
	}
	return threads, nil
}

// NewThread attempts to create a new thread from a board of given public key.
func (c *Container) NewThread(bpk cipher.PubKey, bsk cipher.SecKey, thread *typ.Thread) error {
	c.Lock()
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	thread.MasterBoard = bpk.Hex()
	tRef, e := w.AppendToRefsField("Threads", *thread)
	if e != nil {
		return e
	}
	tp := typ.ThreadPage{Thread: tRef}
	if _, e := w.AppendToRefsField("ThreadPages", tp); e != nil {
		return e
	}
	thread.Ref = cipher.SHA256(tRef).Hex()
	// Prepare thread vote container.
	w.Clear()
	tvc := &typ.ThreadVotesContainer{}
	if e := w.AdvanceFromRoot(tvc, makeThreadVotesContainerFinder(w.Root())); e != nil {
		return e
	}
	tvc.AddThread(tRef)
	if e := w.ReplaceCurrent(*tvc); e != nil {
		return e
	}
	return nil
}

// RemoveThread attempts to remove a thread from a board of given public key.
func (c *Container) RemoveThread(bpk cipher.PubKey, bsk cipher.SecKey, tRef skyobject.Reference) error {
	c.Lock()
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	if e := w.RemoveInRefsByRef("Threads", tRef); e != nil {
		return errors.Wrap(e, "remove thread failed")
	}
	if e := w.RemoveInRefsField("ThreadPages", makeThreadPageFinder(w, tRef)); e != nil {
		return errors.Wrap(e, "remove thread page failed")
	}

	// remove thread votes.
	w.Clear()
	tvc := &typ.ThreadVotesContainer{}
	if e := w.AdvanceFromRoot(tvc, makeThreadVotesContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "obtaining thread vote container failed")
	}
	tvc.RemoveThread(tRef)
	if e := w.ReplaceCurrent(*tvc); e != nil {
		return errors.Wrap(e, "swapping thread vote container failed")
	}
	return nil
}

// GetThreadPage requests a page from a thread
func (c *Container) GetThreadPage(bpk cipher.PubKey, tRef skyobject.Reference) (*typ.Thread, []*typ.Post, error) {
	c.Lock()
	defer c.Unlock()

	w := c.c.LastRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return nil, nil, e
	}
	// Get thread.
	tData, has := c.c.Get(tRef)
	if has == false {
		return nil, nil, errors.New("unable to obtain thread")
	}
	thread := new(typ.Thread)
	if e := encoder.DeserializeRaw(tData, thread); e != nil {
		return nil, nil, e
	}
	thread.Ref = cipher.SHA256(tRef).Hex()
	// Get posts.
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return nil, nil, e
	}
	posts := make([]*typ.Post, len(tp.Posts))
	for i, pRef := range tp.Posts {
		pData, has := c.c.Get(pRef)
		if has == false {
			continue
		}
		posts[i] = new(typ.Post)
		if e := encoder.DeserializeRaw(pData, posts[i]); e != nil {
			return nil, nil, e
		}
		posts[i].Ref = cipher.SHA256(pRef).Hex()
	}
	return thread, posts, nil
}

// GetPosts attempts to obtain posts from a specified board and thread.
func (c *Container) GetPosts(bpk cipher.PubKey, tRef skyobject.Reference) ([]*typ.Post, error) {
	c.Lock()
	defer c.Unlock()

	w := c.c.LastFullRoot(bpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return nil, e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return nil, e
	}
	posts := make([]*typ.Post, len(tp.Posts))
	for i, pRef := range tp.Posts {
		pData, has := c.c.Get(pRef)
		if has == false {
			continue
		}
		posts[i] = new(typ.Post)
		if e := encoder.DeserializeRaw(pData, posts[i]); e != nil {
			return nil, e
		}
		posts[i].Ref = cipher.SHA256(pRef).Hex()
	}
	return posts, nil
}

// NewPost attempts to create a new post in a given board and thread.
func (c *Container) NewPost(bpk cipher.PubKey, bsk cipher.SecKey, tRef skyobject.Reference, post *typ.Post) error {
	c.Lock()
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return e
	}
	t := &typ.Thread{}
	if e := w.GetFromRefField("Thread", t); e != nil {
		return e
	}
	if t.MasterBoard != bpk.Hex() {
		return errors.New("this board is not master of this thread")
	}
	var pRef skyobject.Reference
	var e error
	if pRef, e = w.AppendToRefsField("Posts", *post); e != nil {
		return e
	}
	post.Ref = cipher.SHA256(pRef).Hex()

	// Prepare post vote container.
	w.Clear()
	pvc := &typ.PostVotesContainer{}
	if e := w.AdvanceFromRoot(pvc, makePostVotesContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "unable to obtain post vote container")
	}
	pvc.AddPost(pRef)
	if e := w.ReplaceCurrent(*pvc); e != nil {
		return e
	}
	return nil
}

// RemovePost attempts to remove a post in a given board and thread.
func (c *Container) RemovePost(bpk cipher.PubKey, bsk cipher.SecKey, tRef, pRef skyobject.Reference) error {
	c.Lock()
	defer c.Unlock()

	w := c.c.LastRootSk(bpk, bsk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return e
	}
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return e
	}
	if e := w.RemoveInRefsByRef("Posts", pRef); e != nil {
		return errors.Wrap(e, "post removal failed")
	}

	// Remove post votes.
	w.Clear()
	pvc := &typ.PostVotesContainer{}
	if e := w.AdvanceFromRoot(pvc, makePostVotesContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "unable to obtain post vote container")
	}
	pvc.RemovePost(pRef)
	if e := w.ReplaceCurrent(*pvc); e != nil {
		return errors.Wrap(e, "unable to replace post vote container")
	}
	return nil
}

// ImportThread imports a thread from a board to another board (which this node owns). If already imported replaces it.
func (c *Container) ImportThread(fromBpk, toBpk cipher.PubKey, toBsk cipher.SecKey, tRef skyobject.Reference) error {
	c.Lock()
	defer c.Unlock()

	log.Printf("[CONTAINER] Syncing thread '%s' from board '%s' to board '%s'.",
		tRef.String(), fromBpk.Hex(), toBpk.Hex())

	// Get from 'from' Board.
	w := c.c.LastFullRoot(fromBpk).Walker()
	bc := &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "import thread failed: failed to obtain board "+fromBpk.Hex())
	}

	// Obtain thread and thread page.
	tp := &typ.ThreadPage{}
	if e := w.AdvanceFromRefsField("ThreadPages", tp, makeThreadPageFinder(w, tRef)); e != nil {
		return errors.Wrap(e, "import thread failed: failed to obtain thread page for board "+fromBpk.Hex())
	}
	t := &typ.Thread{}
	if e := w.GetFromRefField("Thread", t); e != nil {
		return errors.Wrap(e, "import thread failed: failed to obtain thread for board "+fromBpk.Hex())
	}

	// Get from 'to' Board.
	w = c.c.LastRootSk(toBpk, toBsk).Walker()
	bc = &typ.BoardContainer{}
	if e := w.AdvanceFromRoot(bc, makeBoardContainerFinder(w.Root())); e != nil {
		return errors.Wrap(e, "import thread failed: failed to obtain board "+toBpk.Hex())
	}
	if e := w.ReplaceInRefsField("ThreadPages", *tp, makeThreadPageFinder(w, tRef)); e != nil {
		/* THREAD DOES NOT EXIST */
		// Append thread and thread page.
		if _, e := w.AppendToRefsField("Threads", *t); e != nil {
			return e
		}
		if _, e := w.AppendToRefsField("ThreadPages", *tp); e != nil {
			return e
		}
		return nil
	}
	return nil
}
