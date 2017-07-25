package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/content"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/state"
)

type BoardsOutput struct {
	Boards       []object.BoardView `json:"boards"`
	MasterBoards []object.BoardView `json:"master_boards"`
}

func getBoards(ctx context.Context, cxo *state.CXO, file *state.UserFile) *BoardsOutput {

	masters := make([]object.BoardView, len(file.Masters))
	for i, sub := range file.Masters {
		masters[i].PublicKey = sub.PubKey.Hex()
		result, e := content.GetBoardResult(ctx, cxo, sub.PubKey)
		if e != nil {
			masters[i].Board = &object.Board{
				Name: "Unavailable Board",
				Desc: e.Error(),
			}
		} else {
			masters[i].Board = result.Board
		}
	}

	subs := make([]object.BoardView, len(file.Subscriptions))
	for i, sub := range file.Subscriptions {
		subs[i].PublicKey = sub.PubKey.Hex()
		result, e := content.GetBoardResult(ctx, cxo, sub.PubKey)
		if e != nil {
			subs[i].Board = &object.Board{
				Name: "Unavailable Board",
				Desc: e.Error(),
			}
		} else {
			subs[i].Board = result.Board
		}
	}

	return &BoardsOutput{MasterBoards: masters, Boards: subs}
}

type BoardPageOutput struct {
	Board   object.BoardView    `json:"board"`
	Threads []object.ThreadView `json:"threads,omitempty"`
}

func getBoardPage(_ context.Context, result *content.Result) *BoardPageOutput {

	out := &BoardPageOutput{
		Board: object.BoardView{
			Board:     result.Board,
			PublicKey: result.GetPK().Hex(),
		},
		Threads: make([]object.ThreadView, len(result.Threads)),
	}

	for i, thread := range result.Threads {
		out.Threads[i] = object.ThreadView{
			Thread:      thread,
			Ref:         thread.R.Hex(),
			AuthorRef:   thread.User.Hex(),
			AuthorAlias: "-", // TODO: Implement.
			Votes:       nil, // TODO: Implement.
		}
	}

	return out
}

type ThreadPageOutput struct {
	*BoardPageOutput
	Thread object.ThreadView `json:"thread"`
	Posts  []object.PostView `json:"posts"`
}

func getThreadPage(ctx context.Context, result *content.Result) *ThreadPageOutput {

	out := &ThreadPageOutput{
		BoardPageOutput: getBoardPage(ctx, result),
		Thread: object.ThreadView{
			Thread:      result.Thread,
			Ref:         result.Thread.R.Hex(),
			AuthorRef:   result.Thread.User.Hex(),
			AuthorAlias: "-",
			Votes:       nil, // TODO: Implement.
		},
		Posts: make([]object.PostView, len(result.Posts)),
	}

	for i, post := range result.Posts {
		out.Posts[i] = object.PostView{
			Post:        post,
			Ref:         post.R.Hex(),
			AuthorRef:   post.User.Hex(),
			AuthorAlias: "-", // TODO: Implement.
			Votes:       nil, // TODO: Implement.
		}
	}

	return out
}
