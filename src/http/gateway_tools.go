package http

import (
	"github.com/skycoin/bbs/src/misc/keys"
	"net/http"
)

func RegisterToolsHandlers(mux *http.ServeMux, _ *Gateway) {

	// Generates a seed.
	mux.HandleFunc("/api/tools/new_seed",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(keys.GenerateSeed())
		})

	// Generates public/private key pair.
	mux.HandleFunc("/api/tools/new_key_pair",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(keys.GenerateKeyPair(&keys.GenerateKeyPairIn{
				Seed: r.FormValue("seed"),
			}))
		})

	// Generates a hash of given string data.
	mux.HandleFunc("/api/tools/hash_string",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(keys.SumSHA256(&keys.SumSHA256In{
				Data: r.FormValue("data"),
			}))
		})

	mux.HandleFunc("/api/tools/sign",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(keys.SignHash(&keys.SignHashIn{
				Hash:   r.FormValue("hash"),
				SecKey: r.FormValue("secret_key"),
			}))
		})
}
