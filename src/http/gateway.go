package http

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/state"
	"log"
	"net/http"
	"os"
)

// Gateway represents what is exposed to HTTP interface.
type Gateway struct {
	l      *log.Logger
	Access *store.Access
	Quit   chan int
}

func (g *Gateway) prepare(mux *http.ServeMux) error {
	g.l = inform.NewLogger(true, os.Stdout, "")

	/*
		<<< NODE >>>
	*/

	// Quits the node.
	mux.HandleFunc("/api/node/quit",
		func(w http.ResponseWriter, r *http.Request) {
			g.Quit <- 0
			send(w, true, nil)
		})

	// Obtains node states.
	mux.HandleFunc("/api/node/stats",
		func(w http.ResponseWriter, r *http.Request) {
			view := struct {
				NodeIsMaster bool `json:"node_is_master"`
			}{
				NodeIsMaster: g.Access.Session.GetCXO().IsMaster(),
			}
			send(w, view, nil)
		})

	/*
		<<< CONNECTIONS >>>
	*/

	// Gets all connections.
	mux.HandleFunc("/api/connections/get_all",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.GetConnections(r.Context())
			send(w, out, e)
		})

	// Creates a new connection.
	mux.HandleFunc("/api/connections/new",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.NewConnection(r.Context(), &state.ConnectionIO{
				Address: r.FormValue("address"),
			})
			send(w, out, e)
		})

	// Deletes a connection.
	mux.HandleFunc("/api/connections/delete",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.DeleteConnection(r.Context(), &state.ConnectionIO{
				Address: r.FormValue("address"),
			})
			send(w, out, e)
		})

	/*
		<<< SUBSCRIPTIONS >>>
	*/

	// Gets all subscriptions (non-master and master).
	mux.HandleFunc("/api/subscriptions/get_all",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.GetSubs(r.Context())
			send(w, out, e)
		})

	// Creates a new subscription.
	mux.HandleFunc("/api/subscriptions/new",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.NewSub(r.Context(), &state.SubscriptionIO{
				PubKey: r.FormValue("public_key"),
			})
			send(w, out, e)
		})

	// Deletes a subscription.
	mux.HandleFunc("/api/subscriptions/remove",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.DeleteSub(r.Context(), &state.SubscriptionIO{
				PubKey: r.FormValue("public_key"),
			})
			send(w, out, e)
		})

	/*
		<<< BOARDS >>>
	*/

	// Gets all boards (non-master and master). TODO
	mux.HandleFunc("/api/boards/get_all",
		func(w http.ResponseWriter, r *http.Request) {
			send(w, true, nil)
		})

	// Creates a new board. TODO
	mux.HandleFunc("/api/boards/new",
		func(w http.ResponseWriter, r *http.Request) {
			send(w, true, nil)
		})

	return nil
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

func send(w http.ResponseWriter, v interface{}, e error) error {
	if e != nil {
		return sendErr(w, e)
	}
	return sendOK(w, v)
}

func sendOK(w http.ResponseWriter, v interface{}) error {
	response := Response{Okay: true, Data: v}
	return sendStatus(w, response, http.StatusOK)
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
	return sendStatus(w, response, status)
}

func sendStatus(w http.ResponseWriter, v interface{}, status int) error {
	data, e := json.Marshal(v)
	if e != nil {
		return e
	}
	sendRaw(w, data, status)
	return nil
}

//func sendRawOK(w http.ResponseWriter, data []byte) {
//	sendRaw(w, data, http.StatusOK)
//}

func sendRaw(w http.ResponseWriter, data []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
