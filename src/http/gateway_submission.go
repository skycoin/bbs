package http

import (
	"github.com/skycoin/bbs/src/store"
	"net/http"
)

func RegisterSubmissionHandlers(mux *http.ServeMux, g *Gateway) {

	mux.HandleFunc("/api/submission/prepare_thread",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.PrepareThread(r.Context(), &store.PrepareThreadIn{
				OfBoardStr: r.FormValue("of_board"),
				Name:       r.FormValue("name"),
				Body:       r.FormValue("body"),
				CreatorStr: r.FormValue("creator"),
			}))
		})

	mux.HandleFunc("/api/submission/prepare_post",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.PreparePost(r.Context(), &store.PreparePostIn{
				OfBoardStr:  r.FormValue("of_board"),
				OfThreadStr: r.FormValue("of_thread"),
				OfPostStr:   r.FormValue("of_post"),
				Name:        r.FormValue("name"),
				Body:        r.FormValue("body"),
				ImagesStr:   r.FormValue("images"),
				CreatorStr:  r.FormValue("creator"),
			}))
		})

	mux.HandleFunc("/api/submission/prepare_thread_vote",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.PrepareThreadVote(r.Context(), &store.PrepareThreadVoteIn{
				OfBoardStr:  r.FormValue("of_board"),
				OfThreadStr: r.FormValue("of_thread"),
				ValueStr:    r.FormValue("value"),
				TagsStr:     r.FormValue("tags"),
				CreatorStr:  r.FormValue("creator"),
			}))
		})

	mux.HandleFunc("/api/submission/prepare_post_vote",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.PreparePostVote(r.Context(), &store.PreparePostVoteIn{
				OfBoardStr: r.FormValue("of_board"),
				OfPostStr:  r.FormValue("of_post"),
				ValueStr:   r.FormValue("value"),
				TagsStr:    r.FormValue("tags"),
				CreatorStr: r.FormValue("creator"),
			}))
		})

	mux.HandleFunc("/api/submission/prepare_user_vote",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.PrepareUserVote(r.Context(), &store.PrepareUserVoteIn{
				OfBoardStr: r.FormValue("of_board"),
				OfUserStr:  r.FormValue("of_user"),
				ValueStr:   r.FormValue("value"),
				TagsStr:    r.FormValue("tags"),
				CreatorStr: r.FormValue("creator"),
			}))
		})

	mux.HandleFunc("/api/submission/finalize",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.FinalizeSubmission(r.Context(), &store.FinalizeSubmissionIn{
				HashStr: r.FormValue("hash"),
				SigStr:  r.FormValue("sig"),
			}))
		})
}
