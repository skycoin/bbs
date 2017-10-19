package http

import (
	"github.com/skycoin/bbs/src/store/object"
	"net/http"
)

func RegisterAdminHandlers(mux *http.ServeMux, g *Gateway) {

	// Quits the node.
	mux.HandleFunc("/api/admin/quit",
		func(w http.ResponseWriter, r *http.Request) {
			g.QuitChan <- 0
			send(w)(true, nil)
		})

	// Obtains node stats. TODO
	mux.HandleFunc("/api/admin/stats",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(true, nil)
		})

	/*
		<<< CONNECTIONS >>>
		>>> Endpoints to handle the node's connections.
	*/

	// Gets all connections.
	mux.HandleFunc("/api/admin/connections/get_all",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.GetConnections(r.Context()))
		})

	// Creates a new connection.
	mux.HandleFunc("/api/admin/connections/new",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.NewConnection(r.Context(), &object.ConnectionIO{
				Address: r.FormValue("address"),
			}))
		})

	// Deletes a connection.
	mux.HandleFunc("/api/admin/connections/delete",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.DeleteConnection(r.Context(), &object.ConnectionIO{
				Address: r.FormValue("address"),
			}))
		})

	/*
		<<< SUBSCRIPTIONS >>>
		>>> Endpoints to handle the node's subscriptions.
	*/

	// Gets all subscriptions (non-master and master).
	mux.HandleFunc("/api/admin/subscriptions/get_all",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.GetSubscriptions(r.Context()))
		})

	// Creates a new subscription.
	mux.HandleFunc("/api/admin/subscriptions/new",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.NewSubscription(r.Context(), &object.BoardIO{
				PubKeyStr: r.FormValue("public_key"),
			}))
		})

	// Deletes a subscription.
	mux.HandleFunc("/api/admin/subscriptions/delete",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.DeleteSubscription(r.Context(), &object.BoardIO{
				PubKeyStr: r.FormValue("public_key"),
			}))
		})

	/*
		<<< CONTENT >>>
		>>> Endpoints to handle content hosted on this node.
	*/

	// Creates and hosts a new board on this node.
	mux.HandleFunc("/api/admin/content/new_board",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.NewBoard(r.Context(), &object.NewBoardIO{
				Seed: r.FormValue("seed"),
				Name: r.FormValue("name"),
				Body: r.FormValue("body"),
			}))
		})

	// Deletes a hosted board from this node.
	mux.HandleFunc("/api/admin/content/delete_board",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.DeleteBoard(r.Context(), &object.BoardIO{
				PubKeyStr: r.FormValue("board_public_key"),
			}))
		})

	// Exports an entire board root to file.
	mux.HandleFunc("/api/admin/content/export_board",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.ExportBoard(r.Context(), &object.ExportBoardIO{
				FilePath:  r.FormValue("file_path"),
				PubKeyStr: r.FormValue("board_public_key"),
			}))
		})

	// Imports an entire board root from file to CXO.
	mux.HandleFunc("/api/admin/content/import_board",
		func(w http.ResponseWriter, r *http.Request) {
			send(w)(g.Access.ImportBoard(r.Context(), &object.ImportBoardIO{
				FilePath:  r.FormValue("file_path"),
				SecKeyStr: r.FormValue("board_secret_key"),
			}))
		})
}
