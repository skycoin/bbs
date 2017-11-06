package http

import (
	"net/http"
	"github.com/skycoin/bbs/src/misc/tag"
)

func RegisterToolsHandlers(mux *http.ServeMux, _ *Gateway) {

	// Generates a seed.
	mux.HandleFunc("/api/tools/new_seed",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(tag.GenerateSeed())
		})

	// Generates public/private key pair.
	mux.HandleFunc("/api/tools/new_key_pair",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(tag.GenerateKeyPair(&tag.GenerateKeyPairIn{
				Seed: r.FormValue("seed"),
			}))
		})

	// Generates a hash of given string data.
	mux.HandleFunc("/api/tools/hash_string",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(tag.SumSHA256(&tag.SumSHA256In{
				Data: r.FormValue("data"),
			}))
		})

	mux.HandleFunc("/api/tools/sign",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(tag.SignHash(&tag.SignHashIn{
				Hash:   r.FormValue("hash"),
				SecKey: r.FormValue("secret_key"),
			}))
		})
}
