package http

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/access/btp"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/store"
	"net/http"
	"time"
)

// Gateway represents what is exposed to HTTP interface.
type Gateway struct {
	BoardAccessor *btp.BoardAccessor
	CXO           *store.CXO
	Quit          chan int
}

func (g *Gateway) prepare(mux *http.ServeMux) error {
	mux.HandleFunc("/api/quit", g.quit())

	mux.HandleFunc("/api/stats/get", g.statsGet())

	mux.HandleFunc("/api/connections/get_all", g.connectionsGetAll())
	mux.HandleFunc("/api/connections/new", g.connectionsNew())
	mux.HandleFunc("/api/connections/remove", g.connectionsRemove())

	mux.HandleFunc("/api/subscriptions/get_all", g.subscriptionsGetAll())
	mux.HandleFunc("/api/subscriptions/new", g.subscriptionsNew())
	mux.HandleFunc("/api/subscriptions/remove", g.subscriptionsRemove())

	mux.HandleFunc("/api/boards/get_all", g.boardsGetAll())
	mux.HandleFunc("/api/boards/new", g.boardsNew())

	return nil
}

func (g *Gateway) quit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		timer := time.NewTimer(10 * time.Second)
		select {
		case g.Quit <- 0:
			sendOK(w, true)
		case <-timer.C:
			sendErr(w, boo.New(boo.Internal, "failed to exit node"))
		}
	}
}

func (g *Gateway) statsGet() http.HandlerFunc {
	type View struct {
		NodeIsMaster   bool   `json:"node_is_master"`
		NodeCXOAddress string `json:"node_cxo_address"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		view := &View{
			NodeIsMaster:   g.CXO.IsMaster(),
			NodeCXOAddress: g.CXO.GetAddress(),
		}
		sendOK(w, *view)
	}
}

/*
	<<< CONNECTIONS >>>
*/

func (g *Gateway) connectionsGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendOK(w, g.CXO.GetConnections())
	}
}

func (g *Gateway) connectionsNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if e := g.CXO.Connect(r.FormValue("address")); e != nil {
			sendErr(w, e)
			return
		}
		sendOK(w, true)
	}
}

func (g *Gateway) connectionsRemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if e := g.CXO.Disconnect(r.FormValue("address")); e != nil {
			sendErr(w, e)
			return
		}
		sendOK(w, true)
	}
}

/*
	<<< SUBSCRIPTIONS >>>
*/

func (g *Gateway) subscriptionsGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, e := g.BoardAccessor.GetSubscriptionsView()
		if e != nil {
			sendErr(w, e)
			return
		}
		sendRawOK(w, data)
	}
}

func (g *Gateway) subscriptionsNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := g.BoardAccessor.NewSubscription(&btp.NewSubscriptionInput{
			Address: r.FormValue("address"),
			PubKey:  r.FormValue("public_key"),
		})
		if e != nil {
			sendErr(w, e)
			return
		}
		sendOK(w, true)
	}
}

func (g *Gateway) subscriptionsRemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		e := g.BoardAccessor.RemoveSubscription(&btp.RemoveSubscriptionInput{
			PubKey: r.FormValue("public_key"),
		})
		if e != nil {
			sendErr(w, e)
			return
		}
		sendOK(w, true)
	}
}

/*
	<<< BOARDS >>>
*/

func (g *Gateway) boardsGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, e := g.BoardAccessor.GetBoards(&btp.GetBoardsInput{
			SortBy: r.FormValue("sort_by"),
		})
		if e != nil {
			sendErr(w, e)
			return
		}
		sendOK(w, *view)
	}
}

func (g *Gateway) boardsNew() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		view, e := g.BoardAccessor.NewBoard(&btp.NewBoardInput{
			Name: r.FormValue("name"),
			Desc: r.FormValue("description"),
			Seed: r.FormValue("seed"),
		})
		if e != nil {
			sendErr(w, e)
			return
		}
		sendOK(w, *view)
	}
}

/*
	<<< HELPER FUNCTIONS >>>
*/

type Error struct {
	Type    int    `json:"type"`
	Message string `json:"message"`
	Details string `json:"details"`
}

type Response struct {
	Okay  bool        `json:"okay"`
	Data  interface{} `json:"data,omitempty"`
	Error *Error      `json:"error,omitempty"`
}

func sendOK(w http.ResponseWriter, v interface{}) error {
	response := Response{Okay: true, Data: v}
	return send(w, response, http.StatusOK)
}

func sendErr(w http.ResponseWriter, e error) error {
	if e == nil {
		return sendOK(w, true)
	}
	eType := boo.Type(e)
	eMsg := boo.Message(eType)
	var status int
	switch eType {
	case boo.Unknown, boo.Internal:
		status = http.StatusInternalServerError
	case boo.NotAuthorised, boo.NotMaster:
		status = http.StatusUnauthorized
	case boo.ObjectNotFound:
		status = http.StatusNotFound
	case boo.ObjectAlreadyExists:
		status = http.StatusConflict
	default:
		status = http.StatusBadRequest
	}

	response := Response{
		Okay: false,
		Error: &Error{
			Type:    eType,
			Message: eMsg,
			Details: e.Error(),
		},
	}
	return send(w, response, status)
}

func send(w http.ResponseWriter, v interface{}, status int) error {
	data, e := json.Marshal(v)
	if e != nil {
		return e
	}
	sendRaw(w, data, status)
	return nil
}

func sendRawOK(w http.ResponseWriter, data []byte) {
	sendRaw(w, data, http.StatusOK)
}

func sendRaw(w http.ResponseWriter, data []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
