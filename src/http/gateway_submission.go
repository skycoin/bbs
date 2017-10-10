package http

import (
	"github.com/skycoin/bbs/src/store/object"
	"net/http"
)

func RegisterSubmissionHandlers(mux *http.ServeMux, g *Gateway) {
	mux.HandleFunc("/api/new_submission",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.SubmitContent(r.Context(), &object.SubmissionIO{
				Body:   []byte(r.FormValue("body")),
				SigStr: r.FormValue("sig"),
			}))
		})
}

func RegisterLegacySubmissionsHandlers(mux *http.ServeMux, g *Gateway) {
	// Submits a new thread on specified board.
	mux.HandleFunc("/api/submission/new_thread",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.NewThread(r.Context(), &object.NewThreadIO{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				Name:           r.FormValue("name"),
				Body:           r.FormValue("body"),
			}))
		})

	// Adds a new text post on specified thread.
	mux.HandleFunc("/api/submission/new_post",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.NewPost(r.Context(), &object.NewPostIO{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				ThreadRefStr:   r.FormValue("thread_ref"),
				PostRefStr:     r.FormValue("post_ref"), // Optional.
				Name:           r.FormValue("name"),
				Body:           r.FormValue("body"),
				ImagesStr:      r.FormValue("images"), // Optional.
			}))
		})

	// Votes on a specified user.
	mux.HandleFunc("/api/submission/vote_user",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.VoteUser(r.Context(), &object.UserVoteIO{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				UserPubKeyStr:  r.FormValue("user_public_key"),
				ModeStr:        r.FormValue("mode"),
				TagStr:         r.FormValue("tag"),
			}))
		})

	// Votes on a specified thread.
	mux.HandleFunc("/api/submission/vote_thread",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.VoteThread(r.Context(), &object.ThreadVoteIO{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				ThreadRefStr:   r.FormValue("thread_ref"),
				ModeStr:        r.FormValue("mode"),
				TagStr:         r.FormValue("tag"),
			}))
		})

	// Votes on a specified post.
	mux.HandleFunc("/api/submission/vote_post",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.VotePost(r.Context(), &object.PostVoteIO{
				BoardPubKeyStr: r.FormValue("board_public_key"),
				PostRefStr:     r.FormValue("post_ref"),
				ModeStr:        r.FormValue("mode"),
				TagStr:         r.FormValue("tag"),
			}))
		})
}
