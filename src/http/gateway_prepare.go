package http

import (
	"github.com/skycoin/bbs/src/store"
	"net/http"
)

func RegisterSubmissionHandlers(mux *http.ServeMux, g *Gateway) {

	mux.HandleFunc("/api/submission/prepare_thread",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.PrepareThread(r.Context(), &store.PrepareThreadIn{
				BoardPubKeyStr:   r.FormValue("board_public_key"),
				Name:             r.FormValue("name"),
				Body:             r.FormValue("body"),
				CreatorPubKeyStr: r.FormValue("creator"),
			}))
		})

	mux.HandleFunc("/api/submission/prepare_post",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.PreparePost(r.Context(), &store.PreparePostIn{
				BoardPubKeyStr:   r.FormValue("board_public_key"),
				ThreadHashStr:    r.FormValue("thread_hash"),
				PostHashStr:      r.FormValue("post_hash"),
				Name:             r.FormValue("name"),
				Body:             r.FormValue("body"),
				ImagesStr:        r.FormValue("images"),
				CreatorPubKeyStr: r.FormValue("creator"),
			}))
		})

	mux.HandleFunc("/api/submission/prepare_thread_vote",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.PrepareThreadVote(r.Context(), &store.PrepareThreadVoteIn{
				BoardPubKeyStr:   r.FormValue("board_public_key"),
				ThreadHashStr:    r.FormValue("thread_hash"),
				ValueStr:         r.FormValue("value"),
				TagsStr:          r.FormValue("tags"),
				CreatorPubKeyStr: r.FormValue("creator"),
			}))
		})

	mux.HandleFunc("/api/submission/prepare_post_vote",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.PreparePostVote(r.Context(), &store.PreparePostVoteIn{
				BoardPubKeyStr:   r.FormValue("board_public_key"),
				PostHashStr:      r.FormValue("post_hash"),
				ValueStr:         r.FormValue("value"),
				TagsStr:          r.FormValue("tags"),
				CreatorPubKeyStr: r.FormValue("creator"),
			}))
		})

	mux.HandleFunc("/api/submission/prepare_user_vote",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.PrepareUserVote(r.Context(), &store.PrepareUserVoteIn{
				BoardPubKeyStr:   r.FormValue("board_public_key"),
				UserPubKeyStr:    r.FormValue("user_public_key"),
				ValueStr:         r.FormValue("value"),
				TagsStr:          r.FormValue("tags"),
				CreatorPubKeyStr: r.FormValue("creator"),
			}))
		})

	mux.HandleFunc("/api/submission/finalize",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.FinalizeSubmission(r.Context(), &store.FinalizeSubmissionIn{
				HashStr: r.FormValue("hash"),
				SigStr:  r.FormValue("sig"),
			}))
		})

	// LEGACY. TODO: (Get rid of it)
	mux.HandleFunc("/api/new_submission",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.SubmitContent(r.Context(), &store.SubmissionIn{
				Body:   []byte(r.FormValue("body")),
				SigStr: r.FormValue("sig"),
			}))
		})
}
