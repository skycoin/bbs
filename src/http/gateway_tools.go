package http

import (
	"github.com/skycoin/bbs/src/misc/keys"
	"github.com/skycoin/skycoin/src/cipher"
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

	// Generates a hash of given string data.
	mux.HandleFunc("/api/tools/hash_string",
		func(w http.ResponseWriter, r *http.Request) {
			data := r.FormValue("data")
			hash := cipher.SumSHA256([]byte(data))
			send(w)(
				struct {
					Data string `json:"data"`
					Hash string `json:"hash"`
				}{
					Data: data,
					Hash: hash.Hex(),
				},
				nil,
			)
		})

	mux.HandleFunc("/api/tools/sign",
		func(w http.ResponseWriter, r *http.Request) {
			hashStr := r.FormValue("hash")
			skStr := r.FormValue("secret_key")

			// Check secret key and hash.
			hash, e := keys.GetHash(hashStr)
			if e != nil {
				send(w)(nil, e)
			}
			sk, e := keys.GetSecKey(skStr)
			if e != nil {
				send(w)(nil, e)
			}

			// Out.
			send(w)(
				struct {
					Hash string `json:"hash"`
					SK   string `json:"secret_key"`
					Sig  string `json:"sig"`
				}{
					Hash: hash.Hex(),
					SK:   sk.Hex(),
					Sig:  cipher.SignHash(hash, sk).Hex(),
				},
				nil,
			)
		})
}
