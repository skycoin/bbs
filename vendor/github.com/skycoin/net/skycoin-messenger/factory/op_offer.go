package factory

import (
	"encoding/json"
	"strings"
	"sync"
)

func init() {
	ops[OP_OFFER_SERVICE] = &sync.Pool{
		New: func() interface{} {
			return new(offer)
		},
	}
}

type offer struct {
	Services *NodeServices
}

func (offer *offer) UnmarshalJSON(data []byte) (err error) {
	ss := &NodeServices{}
	err = json.Unmarshal(data, ss)
	if err != nil {
		return
	}
	offer.Services = ss
	return
}

func (offer *offer) Execute(f *MessengerFactory, conn *Connection) (r resp, err error) {
	if len(offer.Services.ServiceAddress) > 0 {
		remote := conn.GetRemoteAddr().String()
		addr := remote[:strings.LastIndex(remote, ":")]
		lastIndex := strings.LastIndex(offer.Services.ServiceAddress, ":")
		if lastIndex < 0 {
			return
		}
		addr += offer.Services.ServiceAddress[lastIndex:]
		offer.Services.ServiceAddress = addr
	}
	f.discoveryRegister(conn, offer.Services)
	return
}
