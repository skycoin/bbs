package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
)

type MessengersOutput struct {
	Connections []*object.MessengerConnection `json:"connections"`
}

func getMessengers(_ context.Context, cs []*object.MessengerConnection) *MessengersOutput {
	return &MessengersOutput{
		Connections: cs,
	}
}

type AvailableBoardsOutput struct {
	Boards []string `json:"boards"`
}

func getAvailableBoards(pks []cipher.PubKey) *AvailableBoardsOutput {
	out := &AvailableBoardsOutput{
		Boards: make([]string, len(pks)),
	}
	for i, pk := range pks {
		out.Boards[i] = pk.Hex()
	}
	return out
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

type ExportBoardOutput struct {
	FilePath string             `json:"file_path"`
	Board    *object.ContentRep `json:"board"`
}

func getExportBoardOutput(path string, pages *object.PagesJSON) *ExportBoardOutput {
	return &ExportBoardOutput{
		FilePath: path,
		Board:    pages.BoardPage.Board.ToRep(),
	}
}
