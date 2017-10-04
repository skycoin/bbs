package http

import (
	"net/http"
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/skycoin/src/cipher"
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
			var pk cipher.PubKey
			var sk cipher.SecKey
			seed := r.FormValue("seed")
			switch seed {
			case "":
				pk, sk = cipher.GenerateKeyPair()
			default:
				pk, sk = cipher.GenerateDeterministicKeyPair([]byte(seed))
			}
			send(w)(
				struct {
					PubKey string `json:"public_key"`
					SecKey string `json:"secret_key"`
				}{
					PubKey: pk.Hex(),
					SecKey: sk.Hex(),
				},
				nil,
			)
		})
}