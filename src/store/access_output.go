package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
)

type UsersOutput struct {
	Users []object.UserView `json:"users"`
}

func getUsers(_ context.Context, aliases []string) *UsersOutput {
	out := &UsersOutput{
		Users: make([]object.UserView, len(aliases)),
	}
	for i, alias := range aliases {
		out.Users[i] = object.UserView{
			User: object.User{Alias: alias},
		}
	}
	return out
}

type SessionOutput struct {
	LoggedIn bool                 `json:"logged_in"`
	Session  *object.UserFileView `json:"session"`
}

func getSession(_ context.Context, f *object.UserFile) *SessionOutput {
	if f == nil {
		return &SessionOutput{
			LoggedIn: false,
			Session:  nil,
		}
	} else {
		return &SessionOutput{
			LoggedIn: true,
			Session:  f.View(),
		}
	}
}

type ConnectionsOutput struct {
	Connections []object.Connection `json:"connections"`
}

func getConnections(_ context.Context, cs []object.Connection) *ConnectionsOutput {
	return &ConnectionsOutput{
		Connections: cs,
	}
}

type SubscriptionsOutput struct {
	Subscriptions []string `json:"subscriptions"`
}

func getSubscriptions(_ context.Context, ss []cipher.PubKey) *SubscriptionsOutput {
	out := &SubscriptionsOutput{
		Subscriptions: make([]string, len(ss)),
	}
	for i, s := range ss {
		out.Subscriptions[i] = s.Hex()
	}
	return out
}

type SubmissionOutput struct {
	SubAddresses []string `json:"submission_addresses"`
}

func getSubmissionOutput(addresses []string) *SubmissionOutput {
	return &SubmissionOutput{
		SubAddresses: addresses,
	}
}

type BoardsOutput struct {
	MasterBoards []interface{} `json:"master_boards"`
	RemoteBoards []interface{} `json:"remote_boards"`
}

func getBoardsOutput(_ context.Context, m, r []interface{}) *BoardsOutput {
	return &BoardsOutput{
		MasterBoards: m,
		RemoteBoards: r,
	}
}

type BoardOutput struct {
	Board interface{} `json:"board"`
}

func getBoardOutput(v interface{}) *BoardOutput {
	return &BoardOutput{
		Board: v,
	}
}

type FollowPageOutput struct {
	FollowPage interface{} `json:"follow_page"`
}

func getFollowPageOutput(v interface{}) *FollowPageOutput {
	return &FollowPageOutput{
		FollowPage: v,
	}
}
