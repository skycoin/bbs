package http

import (
	"net/http"
	"github.com/skycoin/bbs/src/store"
)

func RegisterSubmissionHandlers(mux *http.ServeMux, g *Gateway) {
	mux.HandleFunc("/api/new_submission",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.SubmitContent(r.Context(), &store.SubmissionIn{
				Body:   []byte(r.FormValue("body")),
				SigStr: r.FormValue("sig"),
			}))
		})
}
