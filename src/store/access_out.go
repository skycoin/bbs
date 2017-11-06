package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/skycoin/src/cipher"
)

type MessengersOut struct {
	Connections []*object.MessengerConnection `json:"connections"`
}

func getMessengersOut(_ context.Context, cs []*object.MessengerConnection) *MessengersOut {
	return &MessengersOut{
		Connections: cs,
	}
}

type AvailableBoardsOut struct {
	Boards []string `json:"boards"`
}

func getAvailableBoardsOut(pks []cipher.PubKey) *AvailableBoardsOut {
	out := &AvailableBoardsOut{
		Boards: make([]string, len(pks)),
	}
	for i, pk := range pks {
		out.Boards[i] = pk.Hex()
	}
	return out
}

type ConnectionsOut struct {
	Connections []object.Connection `json:"connections"`
}

func getConnectionsOut(_ context.Context, cs []object.Connection) *ConnectionsOut {
	return &ConnectionsOut{
		Connections: cs,
	}
}

type SubscriptionsOut struct {
	Subscriptions []string `json:"subscriptions"`
}

func getSubscriptionsOut(_ context.Context, ss []cipher.PubKey) *SubscriptionsOut {
	out := &SubscriptionsOut{
		Subscriptions: make([]string, len(ss)),
	}
	for i, s := range ss {
		out.Subscriptions[i] = s.Hex()
	}
	return out
}

type BoardsOut struct {
	MasterBoards []interface{} `json:"master_boards"`
	RemoteBoards []interface{} `json:"remote_boards"`
}

func getBoardsOut(_ context.Context, m, r []interface{}) *BoardsOut {
	return &BoardsOut{
		MasterBoards: m,
		RemoteBoards: r,
	}
}

type BoardOut struct {
	Board interface{} `json:"board"`
}

func getBoardOut(v interface{}) *BoardOut {
	return &BoardOut{
		Board: v,
	}
}

type FollowPageOut struct {
	FollowPage interface{} `json:"follow_page"`
}

func getFollowPageOut(v interface{}) *FollowPageOut {
	return &FollowPageOut{
		FollowPage: v,
	}
}

type ExportBoardOut struct {
	FilePath string             `json:"file_path"`
	Board    *object.ContentRep `json:"board"`
}

func getExportBoardOut(path string, pages *object.PagesJSON) *ExportBoardOut {
	return &ExportBoardOut{
		FilePath: path,
		Board:    pages.BoardPage.Board.ToRep(),
	}
}
