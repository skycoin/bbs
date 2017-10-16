package rpc

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/store/object"
	"net/rpc"
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
	<<< CONNECTIONS >>>
*/

func GetConnections() (string, interface{}) {
	return method("GetConnections"), empty()
}

func NewConnection(in *object.ConnectionIO) (string, interface{}) {
	return method("NewConnection"), in
}

func DeleteConnection(in *object.ConnectionIO) (string, interface{}) {
	return method("DeleteConnection"), in
}

/*
	<<< SUBSCRIPTIONS >>>
*/

func GetSubscriptions() (string, interface{}) {
	return method("GetSubscriptions"), empty()
}

func NewSubscription(in *object.BoardIO) (string, interface{}) {
	return method("NewSubscription"), in
}

func DeleteSubscription(in *object.BoardIO) (string, interface{}) {
	return method("DeleteSubscription"), in
}

/*
	<<< CONTENT : ADMIN >>>
*/

func NewBoard(in *object.NewBoardIO) (string, interface{}) {
	return method("NewBoard"), in
}

func DeleteBoard(in *object.BoardIO) (string, interface{}) {
	return method("DeleteBoard"), in
}

func ExportBoard(in *object.ExportBoardIO) (string, interface{}) {
	return method("ExportBoard"), in
}

func ImportBoard(in *object.ExportBoardIO) (string, interface{}) {
	return method("ImportBoard"), in
}

/*
	<<< CONTENT >>>
*/

func GetBoards() (string, interface{}) {
	return method("GetBoards"), empty()
}

func GetBoard(in *object.BoardIO) (string, interface{}) {
	return method("GetBoard"), in
}

func GetBoardPage(in *object.BoardIO) (string, interface{}) {
	return method("GetBoardPage"), in
}

func GetThreadPage(in *object.ThreadIO) (string, interface{}) {
	return method("GetThreadPage"), in
}

func GetFollowPage(in *object.UserIO) (string, interface{}) {
	return method("GetFollowPage"), in
}

/*
	<<< CONTENT : SUBMISSION >>>
*/

func NewThread(in *object.NewThreadIO) (string, interface{}) {
	return method("NewThread"), in
}

func NewPost(in *object.NewPostIO) (string, interface{}) {
	return method("NewPost"), in
}

func VoteThread(in *object.ThreadVoteIO) (string, interface{}) {
	return method("VoteThread"), in
}

func VotePost(in *object.PostVoteIO) (string, interface{}) {
	return method("VotePost"), in
}

func VoteUser(in *object.UserVoteIO) (string, interface{}) {
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
