package factory

import (
	"time"
)

type ConnConfig struct {
	Reconnect     bool
	ReconnectWait time.Duration

	// callbacks

	FindServiceNodesByKeysCallback func(resp *QueryResp)

	FindServiceNodesByAttributesCallback func(resp *QueryByAttrsResp)

	// call after connected to server
	OnConnected func(connection *Connection)

	Creator *MessengerFactory
}
