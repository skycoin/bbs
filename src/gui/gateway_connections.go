package gui

import (
	"github.com/pkg/errors"
	"net/http"
)

// Connections represents the connections endpoint group.
type Connections struct {
	*Gateway
}

// GetAll gets all connections.
func (g *Connections) GetAll(w http.ResponseWriter, r *http.Request) {
	send(w, g.getAll(), http.StatusOK)
}

func (g *Connections) getAll() []string {
	return g.container.GetConnections()
}

// Add adds a connection.
func (g *Connections) Add(w http.ResponseWriter, r *http.Request) {
	if e := g.add(r.FormValue("address")); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *Connections) add(addr string) error {
	return g.container.Connect(addr)
}

// Remove removes a connection.
func (g *Connections) Remove(w http.ResponseWriter, r *http.Request) {
	if e := g.remove(r.FormValue("address")); e != nil {
		send(w, e.Error(), http.StatusBadRequest)
		return
	}
	send(w, true, http.StatusOK)
}

func (g *Connections) remove(addr string) error {
	boards := g.boardSaver.GetOfAddress(addr)
	if len(boards) == 0 {
		return g.container.Disconnect(addr)
	}
	return errors.Errorf("currently subscribed to %d boards under address %s",
		len(boards), addr)
}
