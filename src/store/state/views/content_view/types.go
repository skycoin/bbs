package content_view

import (
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
	"log"
	"sync"
)

type BoardRep struct {
	sync.Mutex
	PubKey       string          `json:"public_key"`
	Name         string          `json:"name"`
	Body         string          `json:"body"`
	Created      int64           `json:"created"`
	SubAddresses []string        `json:"submission_addresses"`
	Threads      []cipher.SHA256 `json:"-"`
}

func (r *BoardRep) Fill(pk cipher.PubKey, board *object.Board) *BoardRep {
	r.Lock()
	defer r.Unlock()

	data := object.GetData(board)
	r.PubKey = pk.Hex()
	r.Name = data.Name
	r.Body = data.Body
	r.Created = board.Created
	r.SubAddresses = data.SubAddresses
	return r
}

type ThreadRep struct {
	sync.Mutex
	Ref     string          `json:"ref"`
	Name    string          `json:"name"`
	Body    string          `json:"body"`
	Created int64           `json:"created"`
	Creator string          `json:"creator"`
	Posts   []cipher.SHA256 `json:"-"`
}

func (r *ThreadRep) FillThread(thread *object.Thread, mux *sync.Mutex) *ThreadRep {
	r.Lock()
	defer r.Unlock()

	data := object.GetData(thread)
	r.Ref = thread.R.Hex()
	r.Name = data.Name
	r.Body = data.Body
	r.Created = thread.Created
	r.Creator = thread.Creator.Hex()
	return nil
}

func (r *ThreadRep) Fill(tPage *object.ThreadPage, mux *sync.Mutex) *ThreadRep {
	r.Lock()
	defer r.Unlock()

	t, e := tPage.GetThread(mux)
	if e != nil {
		log.Println("ThreadRep.Fill() Error:", e)
		return nil
	}
	data := object.GetData(t)
	r.Ref = t.R.Hex()
	r.Name = data.Name
	r.Body = data.Body
	r.Created = t.Created
	r.Creator = t.Creator.Hex()
	return nil
}

type PostRep struct {
	sync.Mutex
	Ref     string `json:"ref"`
	Name    string `json:"name"`
	Body    string `json:"body"`
	Created int64  `json:"created"`
	Creator string `json:"creator"`
}

func (r *PostRep) Fill(post *object.Post, mux *sync.Mutex) *PostRep {
	r.Lock()
	defer r.Unlock()

	data := object.GetData(post)
	r.Ref = post.R.Hex()
	r.Name = data.Name
	r.Body = data.Body
	r.Created = post.Created
	r.Creator = post.Creator.Hex()
	return nil
}
