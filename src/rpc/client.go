package rpc

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/store"
	"github.com/skycoin/bbs/src/store/object"
	"net/rpc"
)

type Call func(method string, in, out interface{}) error

type Action func(call Call) string

func Send(address string, action Action) (string, error) {
	client, e := rpc.Dial("tcp", address)
	if e != nil {
		return "", e
	}
	defer client.Close()
	return action(client.Call), nil
}

/*
	<<< CONNECTIONS >>>
*/

func GetConnections() Action {
	return func(call Call) string {
		out := new(store.ConnectionsOutput)
		if e := call("Gateway.GetConnections", &struct{}{}, out); e != nil {
			return errString(e)
		} else {
			return jsonString(out)
		}
	}
}

func NewConnection(in *object.ConnectionIO) Action {
	return func(call Call) string {
		out := new(store.ConnectionsOutput)
		if e := call("Gateway.NewConnection", in, out); e != nil {
			return errString(e)
		} else {
			return jsonString(out)
		}
	}
}

func DeleteConnection(in *object.ConnectionIO) Action {
	return func(call Call) string {
		out := new(store.ConnectionsOutput)
		if e := call("Gateway.DeleteConnection", in, out); e != nil {
			return errString(e)
		} else {
			return jsonString(out)
		}
	}
}

/*
	<<< SUBSCRIPTIONS >>>
*/

func GetSubscriptions() Action {
	return func(call Call) string {
		out := new(store.SubscriptionsOutput)
		if e := call("Gateway.GetSubscriptions", &struct{}{}, out); e != nil {
			return errString(e)
		} else {
			return jsonString(out)
		}
	}
}

func NewSubscription(in *object.BoardIO) Action {
	return func(call Call) string {
		out := new(store.SubscriptionsOutput)
		if e := call("Gateway.NewSubscription", in, out); e != nil {
			return errString(e)
		} else {
			return jsonString(out)
		}
	}
}

func DeleteSubscription(in *object.BoardIO) Action {
	return func(call Call) string {
		out := new(store.SubscriptionsOutput)
		if e := call("Gateway.DeleteSubscription", in, out); e != nil {
			return errString(e)
		} else {
			return jsonString(out)
		}
	}
}

/*
	<<< CONTENT : ADMIN >>>
*/

func NewBoard(in *object.NewBoardIO) Action {
	return func(call Call) string {
		out := new(store.BoardsOutput)
		if e := call("Gateway.NewBoard", in, out); e != nil {
			return errString(e)
		} else {
			return jsonString(out)
		}
	}
}

func DeleteBoard(in *object.BoardIO) Action {
	return func(call Call) string {
		out := new(store.BoardsOutput)
		if e := call("Gateway.DeleteBoard", in, out); e != nil {
			return errString(e)
		} else {
			return jsonString(out)
		}
	}
}

func ExportBoard(in *object.ExportBoardIO) Action {
	return func(call Call) string {
		out := new(store.ExportBoardOutput)
		if e := call("Gateway.ExportBoard", in, out); e != nil {
			return errString(e)
		} else {
			return jsonString(out)
		}
	}
}

func ImportBoard(in *object.ExportBoardIO) Action {
	return func(call Call) string {
		out := new(store.ExportBoardOutput)
		if e := call("Gateway.ImportBoard", in, out); e != nil {
			return errString(e)
		} else {
			return jsonString(out)
		}
	}
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func jsonString(v interface{}) string {
	data, _ := json.MarshalIndent(v, "", "  ")
	return "[OK] " + string(data)
}

func errString(e error) string {
	v := struct {
		Message string `json:"message"`
	}{
		Message: e.Error(),
	}
	data, _ := json.MarshalIndent(v, "", "  ")
	return "[ERROR] " + string(data)
}
