package http

import (
	"bytes"
	"encoding/json"
	"github.com/skycoin/bbs/src/misc/boo"
	"github.com/skycoin/bbs/src/misc/inform"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
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
		<<< TOOLS >>>
	*/

	// Generates a seed.
	mux.HandleFunc("/api/tools/new_seed",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := keys.GenerateSeed()
			send(w, out, e)
		})

	// Generates public/private key pair.
	mux.HandleFunc("/api/tools/new_key_pair",
		func(w http.ResponseWriter, r *http.Request) {
			var pk cipher.PubKey
			var sk cipher.SecKey
			seed := r.FormValue("seed")
			switch seed {
			case "":
				pk, sk = cipher.GenerateKeyPair()
			default:
				pk, sk = cipher.GenerateDeterministicKeyPair([]byte(seed))
			}
			out := struct {
				PubKey string `json:"public_key"`
				SecKey string `json:"secret_key"`
			}{
				PubKey: pk.Hex(),
				SecKey: sk.Hex(),
			}
			send(w, out, nil)
		})

	/*
		<<< SESSION >>>
	*/

	// Lists all users.
	mux.HandleFunc("/api/session/users/get_all",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.GetUsers(r.Context())
			send(w, out, e)
		})

	// Creates a new user.
	mux.HandleFunc("/api/session/users/new",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.NewUser(r.Context(), &object.NewUserIO{
				Seed:  r.FormValue("seed"),
				Alias: r.FormValue("alias"),
				//Password: r.FormValue("password"), TODO: Remove.
			})
			send(w, out, e)
		})

	// Deletes a user.
	mux.HandleFunc("/api/session/users/delete",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.DeleteUser(r.Context(), r.FormValue("alias"))
			send(w, out, e)
		})

	// User login.
	mux.HandleFunc("/api/session/login",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.Login(r.Context(), &object.LoginIO{
				Alias: r.FormValue("alias"),
				//Password: r.FormValue("password"), // TODO: Remove.
			})
			send(w, out, e)
		})

	// User logout.
	mux.HandleFunc("/api/session/logout",
		func(w http.ResponseWriter, r *http.Request) {
			e := g.Access.Logout(r.Context())
			send(w, true, e)
		})

	mux.HandleFunc("/api/session/get_info",
		func(w http.ResponseWriter, r *http.Request) {
			file, e := g.Access.GetSession(r.Context())
			send(w, file, e)
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
			out, e := g.Access.NewConnection(r.Context(), &object.ConnectionIO{
				Address: r.FormValue("address"),
			})
			send(w, out, e)
		})

	// Deletes a connection.
	mux.HandleFunc("/api/connections/delete",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.DeleteConnection(r.Context(), &object.ConnectionIO{
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
			out, e := g.Access.NewSub(r.Context(), &object.BoardIO{
				PubKeyStr: r.FormValue("public_key"),
			})
			send(w, out, e)
		})

	// Deletes a subscription.
	mux.HandleFunc("/api/subscriptions/delete",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.DeleteSub(r.Context(), &object.BoardIO{
				PubKeyStr: r.FormValue("public_key"),
			})
			send(w, out, e)
		})

	/*
		<<< BOARDS >>>
	*/

	// Gets all boards (non-master and master).
	mux.HandleFunc("/api/boards/get_all",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.GetBoards(r.Context())
			send(w, out, e)
		})

	// Creates a new board.
	mux.HandleFunc("/api/boards/new",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.NewBoard(r.Context(), &object.NewBoardIO{
				Seed: r.FormValue("seed"),
				Name: r.FormValue("name"),
				Desc: r.FormValue("description"),
				SubmissionAddressesStr: r.FormValue("submission_addresses"),
				ConnectionsStr:         r.FormValue("connections"),
			})
			send(w, out, e)
		})

	// Deletes a board (That you own).
	mux.HandleFunc("/api/boards/delete",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.DeleteBoard(r.Context(), &object.BoardIO{
				PubKeyStr: r.FormValue("public_key"),
			})
			send(w, out, e)
		})

	// Adds a new submission address to specified board.
	mux.HandleFunc("/api/boards/submission_addresses/new",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.NewSubmissionAddress(r.Context(), &object.AddressIO{
				PubKeyStr: r.FormValue("public_key"),
				Address:   r.FormValue("address"),
			})
			send(w, out, e)
		})

	// Deletes a submission address from specified board.
	mux.HandleFunc("/api/boards/submission_addresses/delete",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.DeleteSubmissionAddress(r.Context(), &object.AddressIO{
				PubKeyStr: r.FormValue("public_key"),
				Address:   r.FormValue("address"),
			})
			send(w, out, e)
		})

	// Gets a view of a board (Includes board information and threads).
	mux.HandleFunc("/api/boards/page/get",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.GetBoardPage(r.Context(), &object.BoardIO{
				PubKeyStr: r.FormValue("public_key"),
			})
			send(w, out, e)
		})

	/*
		<<< THREADS >>>
	*/

	// Creates a new thread (Uses RPC for non-master).
	mux.HandleFunc("/api/threads/new",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.NewThread(r.Context(), &object.NewThreadIO{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				Title:          r.FormValue("title"),
				Body:           r.FormValue("body"),
			})
			send(w, out, e)
		})

	// Deletes a thread from specified board.
	mux.HandleFunc("/api/threads/delete",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.DeleteThread(r.Context(), &object.ThreadIO{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				ThreadRefStr:   r.FormValue("thread_reference"),
			})
			send(w, out, e)
		})

	// Adds/changes vote for a thread.
	mux.HandleFunc("/api/threads/vote",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.VoteThread(r.Context(), &object.VoteThreadIO{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				ThreadRefStr:   r.FormValue("thread_reference"),
				ModeStr:        r.FormValue("mode"),
				TagStr:         r.FormValue("tag"),
			})
			send(w, out, e)
		})

	// Gets a view of a thread (Includes board, thread information and posts).
	mux.HandleFunc("/api/threads/page/get",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.GetThreadPage(r.Context(), &object.ThreadIO{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				ThreadRefStr:   r.FormValue("thread_reference"),
			})
			send(w, out, e)
		})

	/*
		<<< POSTS >>>
	*/

	// Creates a new post on thread of board (Uses RPC for non-master).
	mux.HandleFunc("/api/posts/new",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.NewPost(r.Context(), &object.NewPostIO{
				NewThreadIO: object.NewThreadIO{
					BoardPubKeyStr: r.FormValue("board_public_key"),
					Title:          r.FormValue("title"),
					Body:           r.FormValue("body"),
				},
				ThreadRefStr: r.FormValue("thread_reference"),
			})
			send(w, out, e)
		})

	// Deletes a post on thread of board (Uses RPC for non-master).
	mux.HandleFunc("/api/posts/delete",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.DeletePost(r.Context(), &object.PostIO{
				ThreadIO: object.ThreadIO{
					BoardPubKeyStr: r.FormValue("board_public_key"),
					ThreadRefStr:   r.FormValue("thread_reference"),
				},
				PostRefStr: r.FormValue("post_reference"),
			})
			send(w, out, e)
		})

	// Adds/changes vote for post.
	mux.HandleFunc("/api/posts/vote",
		func(w http.ResponseWriter, r *http.Request) {
			out, e := g.Access.VotePost(r.Context(), &object.VotePostIO{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				PostRefStr:     r.FormValue("post_reference"),
				ModeStr:        r.FormValue("mode"),
				TagStr:         r.FormValue("tag"),
			})
			send(w, out, e)
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

//func sendRawOK(w http.ResponseWriter, data []byte) {
//	sendRaw(w, data, http.StatusOK)
//}

func sendRaw(w http.ResponseWriter, data []byte, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
