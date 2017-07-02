package gui

import "net/http"

// Stats represents the stats endpoint group.
type Stats struct {
	*Gateway
}

// Get gets all stats.
func (g *Stats) Get(w http.ResponseWriter, r *http.Request) {
	send(w, g.get(), http.StatusOK)
}

func (g *Stats) get() *StatsView {
	return &StatsView{
		NodeIsMaster:   g.config.Master(),
		NodeCXOAddress: g.container.GetAddress(),
	}
}
