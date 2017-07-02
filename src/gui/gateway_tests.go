package gui

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/skycoin/bbs/src/misc"
	"github.com/skycoin/bbs/src/store/typ"
	"net/http"
	"strconv"
)

// Tests represents the tests endpoint group.
type Tests struct {
	*Gateway
}

// AddFilledBoard creates a new board with given seed, filled with threads and posts.
func (g *Tests) AddFilledBoard(w http.ResponseWriter, r *http.Request) {
	seed := r.FormValue("seed")

	threads, e := strconv.Atoi(r.FormValue("threads"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}

	minPosts, e := strconv.Atoi(r.FormValue("min_posts"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}

	maxPosts, e := strconv.Atoi(r.FormValue("max_posts"))
	if e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}

	if e := g.addFilledBoard(seed, threads, minPosts, maxPosts); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *Gateway) addFilledBoard(seed string, threads, minPosts, maxPosts int) error {
	if threads < 0 || minPosts < 0 || maxPosts < 0 || maxPosts-minPosts < 0 {
		return errors.New("invalid inputs")
	}
	b := &typ.Board{
		Name: fmt.Sprintf("Test Board '%s'", seed),
		Desc: fmt.Sprintf("A board with '%s' as seed and %d threads.", seed, threads),
	}
	b.SetMeta(new(typ.BoardMeta))
	bi, e := g.Boards.add(b, seed)
	if e != nil {
		return e
	}
	bpk := bi.Config.GetPK()
	for i := 1; i <= threads; i++ {
		t := &typ.Thread{
			Name: fmt.Sprintf("Thread %d", i),
			Desc: fmt.Sprintf("A test thread on board with seed '%s'.", seed),
		}
		if e := g.Threads.add(bpk, t); e != nil {
			return errors.New("on creating thread " + string(i) + "; " + e.Error())
		}
		nPosts, e := misc.MakeIntBetween(minPosts, maxPosts)
		if e != nil {
			return e
		}
		for j := 1; j <= nPosts; j++ {

			p := &typ.Post{
				Title: fmt.Sprintf("Post %d", j),
				Body:  fmt.Sprintf("This is request %d on thread %d.", j, i),
			}
			if e := g.Posts.add(bpk, t.GetRef(), p); e != nil {
				return e
			}
		}
	}
	return nil
}
