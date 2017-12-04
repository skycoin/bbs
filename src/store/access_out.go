package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/skycoin/src/cipher"
)

type SubmissionOut struct {
	NewSubmission   *object.ContentRep `json:"new_submission"`
	NewVotesSummary *state.VoteRepView `json:"new_votes_summary"`
}

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
	ActiveConnections []object.Connection `json:"connections"`
	SavedConnections  []object.Connection `json:"saved_connections"`
}

func getConnectionsOut(_ context.Context, active, saved []object.Connection) *ConnectionsOut {
	return &ConnectionsOut{
		ActiveConnections: active,
		SavedConnections:  saved,
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
