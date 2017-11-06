package rpc

import (
	"encoding/json"
	"net/rpc"
	"path/filepath"
	"github.com/skycoin/bbs/src/store"
)

func Send(address string) func(method string, in interface{}) string {
	return func(method string, in interface{}) string {
		client, e := rpc.Dial("tcp", address)
		if e != nil {
			return errString(e)
		}
		defer client.Close()
		var out string
		if e := client.Call(method, in, &out); e != nil {
			return errString(e)
		} else {
			return okString(out)
		}
	}
}

func Do(out interface{}, e error) string {
	if e != nil {
		return errString(e)
	} else if data, e := json.MarshalIndent(out, "", "  "); e != nil {
		return errString(e)
	} else {
		return okString(string(data))
	}
}

/*
	<<< CONNECTIONS : MESSENGER >>>
*/

func GetMessengerConnections() (string, interface{}) {
	return method("GetMessengerConnections"), empty()
}

func NewMessengerConnection(in *store.ConnectionIn) (string, interface{}) {
	return method("NewMessengerConnection"), in
}

func DeleteMessengerConnection(in *store.ConnectionIn) (string, interface{}) {
	return method("DeleteMessengerConnection"), in
}

func Discover() (string, interface{}) {
	return method("Discover"), empty()
}

/*
	<<< CONNECTIONS >>>
*/

func GetConnections() (string, interface{}) {
	return method("GetConnections"), empty()
}

func NewConnection(in *store.ConnectionIn) (string, interface{}) {
	return method("NewConnection"), in
}

func DeleteConnection(in *store.ConnectionIn) (string, interface{}) {
	return method("DeleteConnection"), in
}

/*
	<<< SUBSCRIPTIONS >>>
*/

func GetSubscriptions() (string, interface{}) {
	return method("GetSubscriptions"), empty()
}

func NewSubscription(in *store.BoardIn) (string, interface{}) {
	return method("NewSubscription"), in
}

func DeleteSubscription(in *store.BoardIn) (string, interface{}) {
	return method("DeleteSubscription"), in
}

/*
	<<< CONTENT : ADMIN >>>
*/

func NewBoard(in *store.NewBoardIn) (string, interface{}) {
	return method("NewBoard"), in
}

func DeleteBoard(in *store.BoardIn) (string, interface{}) {
	return method("DeleteBoard"), in
}

func ExportBoard(in *store.ExportBoardIn) (string, interface{}) {
	in.FilePath, _ = filepath.Abs(in.FilePath)
	return method("ExportBoard"), in
}

func ImportBoard(in *store.ImportBoardIn) (string, interface{}) {
	in.FilePath, _ = filepath.Abs(in.FilePath)
	return method("ImportBoard"), in
}

/*
	<<< CONTENT >>>
*/

func GetBoards() (string, interface{}) {
	return method("GetBoards"), empty()
}

func GetBoard(in *store.BoardIn) (string, interface{}) {
	return method("GetBoard"), in
}

func GetBoardPage(in *store.BoardIn) (string, interface{}) {
	return method("GetBoardPage"), in
}

func GetThreadPage(in *store.ThreadIn) (string, interface{}) {
	return method("GetThreadPage"), in
}

func GetFollowPage(in *store.UserIn) (string, interface{}) {
	return method("GetFollowPage"), in
}

/*
	<<< CONTENT : SUBMISSION >>>
*/

func NewThread(in *store.NewThreadIn) (string, interface{}) {
	return method("NewThread"), in
}

func NewPost(in *store.NewPostIn) (string, interface{}) {
	return method("NewPost"), in
}

func VoteThread(in *store.VoteThreadIn) (string, interface{}) {
	return method("VoteThread"), in
}

func VotePost(in *store.VotePostIn) (string, interface{}) {
	return method("VotePost"), in
}

func VoteUser(in *store.VoteUserIn) (string, interface{}) {
	return method("VoteUser"), in
}

/*
	<<< HELPER FUNCTIONS >>>
*/

func method(v string) string {
	return "Gateway." + v
}

func okString(v string) string {
	return "[OK] " + v
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

func empty() *struct{} {
	return &struct{}{}
}
