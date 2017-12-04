package http

import (
	"bytes"
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/store"
	"log"
	"net/http"
	"os"
)

// Gateway represents what is exposed to HTTP interface.
type Gateway struct {
	l        *log.Logger
	Access   *store.Access
	QuitChan chan int
}

func (g *Gateway) host(mux *http.ServeMux) error {
	g.l = inform.NewLogger(true, os.Stdout, "GATEWAY")

	// For tools.
	RegisterToolsHandlers(mux, g)

	// For submission.
	RegisterSubmissionHandlers(mux, g)

	// Gets a list of boards; remote and master (boards that this node owns).
	mux.HandleFunc("/api/get_boards",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.GetBoards(r.Context()))
		})

	// Gets a single board.
	mux.HandleFunc("/api/get_board",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.GetBoard(r.Context(), &store.BoardIn{
				PubKeyStr:     r.FormValue("board_public_key"),
				UserPubKeyStr: r.FormValue("perspective"),
			}))
		})

	// Obtains a view of a board including it's children threads.
	mux.HandleFunc("/api/get_board_page",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.GetBoardPage(r.Context(), &store.BoardIn{
				PubKeyStr:     r.FormValue("board_public_key"),
				UserPubKeyStr: r.FormValue("perspective"),
			}))
		})

	// Gets a view of a thread including it's children posts.
	mux.HandleFunc("/api/get_thread_page",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.GetThreadPage(r.Context(), &store.ThreadIn{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				ThreadRefStr:   r.FormValue("thread_ref"),
				UserPubKeyStr:  r.FormValue("perspective"),
			}))
		})

	// Gets a view of following/avoiding of specified user.
	mux.HandleFunc("/api/get_user_profile",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.GetFollowPage(r.Context(), &store.UserIn{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				UserPubKeyStr:  r.FormValue("user_public_key"),
			}))
		})

	// Gets a view of all participating users.
	mux.HandleFunc("/api/get_participants",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.GetParticipants(r.Context(), &store.BoardIn{
				PubKeyStr: r.FormValue("board_public_key"),
			}))
		})

	// Lists boards that have been discovered, but not subscribed to.
	mux.HandleFunc("/api/get_available_boards",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.GetAvailableBoards(r.Context()))
		})

	return nil
}

/*
	<<< HELPER FUNCTIONS >>>
*/

type Error struct {
	Type    int    `json:"type"`
	Title   string `json:"title"`
	Details string `json:"details"`
}

type Response struct {
	Okay  bool        `json:"okay"`
	Data  interface{} `json:"data,omitempty"`
	Error *Error      `json:"error,omitempty"`
}

func send(w http.ResponseWriter) func(v interface{}, e error) error {
	return func(v interface{}, e error) error {
		if e != nil {
			return sendErr(w, e)
		}
		return sendOK(w, v)
	}
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
	eTitle := boo.Message(eType)

	var status int

	switch eType {
	case boo.Unknown, boo.Internal:
		status = http.StatusInternalServerError
	case boo.NotAuthorised, boo.NotMaster:
		status = http.StatusUnauthorized
	case boo.NotFound:
		status = http.StatusNotFound
	case boo.AlreadyExists:
		status = http.StatusConflict
	default:
		status = http.StatusBadRequest
	}

	d := e.Error()
	details := string(bytes.ToUpper([]byte{d[0]})) + d[1:] + "."

	response := Response{
		Okay: false,
		Error: &Error{
			Type:    eType,
			Title:   eTitle,
			Details: details,
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

func sendRaw(w http.ResponseWriter, data []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
