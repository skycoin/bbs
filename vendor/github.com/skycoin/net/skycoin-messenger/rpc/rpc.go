package rpc

import (
	"net"
	"net/http"
	"net/rpc"
)

func ServeRPC(address string) error {
	rpc.Register(&Gateway{})
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", address)
	if e != nil {
		return e
	}
	return http.Serve(l, nil)
}
