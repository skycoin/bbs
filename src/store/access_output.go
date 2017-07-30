package store

import (
	"context"
	"github.com/skycoin/bbs/src/store/content"
	"github.com/skycoin/bbs/src/store/object"
	"github.com/skycoin/bbs/src/store/session"
	"github.com/skycoin/bbs/src/store/state"
	"github.com/skycoin/cxo/skyobject"
)

type UsersOutput struct {
	Users []object.UserView `json:"users"`
}

func getUsers(ctx context.Context, aliases []string) *UsersOutput {
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

type ConnectionsOutput struct {
	Connections []object.ConnectionView `json:"connections"`
}

func getConnections(cxo *session.CXO, file *session.UserFile) (*ConnectionsOutput, error) {
	actives, e := cxo.GetConnections()
	if e != nil {
		return nil, e
	}
	activeMap := make(map[string]bool)
	for _, address := range actives {
		activeMap[address] = true
	}

	out := new(ConnectionsOutput)
	for _, address := range file.Connections {
		out.Connections = append(out.Connections, object.ConnectionView{
			Address: address,
			Active:  activeMap[address],
		})
	}

	return out, nil
}

type SubsOutput struct {
	Subscriptions       []object.SubscriptionView `json:"subscriptions"`
	MasterSubscriptions []object.SubscriptionView `json:"master_subscriptions"`
}

func getSubs(_ context.Context, cxo *session.CXO, file *session.UserFile) *SubsOutput {
	view := file.GenerateView(cxo)
	return &SubsOutput{
		Subscriptions:       view.Subscriptions,
		MasterSubscriptions: view.Masters,
	}
}

type BoardsOutput struct {
	Boards       []object.BoardView `json:"boards"`
	MasterBoards []object.BoardView `json:"master_boards"`
}

func getBoards(ctx context.Context, cxo *session.CXO, file *session.UserFile) *BoardsOutput {

	masters := make([]object.BoardView, len(file.Masters))
	for i, sub := range file.Masters {
		masters[i].PublicKey = sub.PubKey.Hex()
		root, _ := cxo.GetRoot(sub.PubKey)
		result, e := content.GetBoardResult(ctx, root)
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
		root, _ := cxo.GetRoot(sub.PubKey)
		result, e := content.GetBoardResult(ctx, root)
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

func getBoardPage(_ context.Context, compiler *state.Compiler, result *content.Result) *BoardPageOutput {

	out := &BoardPageOutput{
		Board: object.BoardView{
			Board:     result.Board,
			PublicKey: result.GetPK().Hex(),
		},
		Threads: make([]object.ThreadView, len(result.Threads)),
	}

	bState := compiler.GetBoard(result.GetPK())

	for i, thread := range result.Threads {
		out.Threads[i] = object.ThreadView{
			Thread:      thread,
			Ref:         thread.R.Hex(),
			AuthorRef:   thread.User.Hex(),
			AuthorAlias: "-", // TODO: Implement.
			Votes:       bState.GetThreadVotes(skyobject.Reference(thread.R)),
		}
	}

	return out
}

type ThreadPageOutput struct {
	*BoardPageOutput
	Thread object.ThreadView `json:"thread"`
	Posts  []object.PostView `json:"posts"`
}

func getThreadPage(ctx context.Context, compiler *state.Compiler, result *content.Result) *ThreadPageOutput {

	out := &ThreadPageOutput{
		BoardPageOutput: getBoardPage(ctx, compiler, result),
		Thread: object.ThreadView{
			Thread:      result.Thread,
			Ref:         result.Thread.R.Hex(),
			AuthorRef:   result.Thread.User.Hex(),
			AuthorAlias: "-",
			Votes:       nil, // TODO: Implement.
		},
		Posts: make([]object.PostView, len(result.Posts)),
	}

	bState := compiler.GetBoard(result.GetPK())

	for i, post := range result.Posts {
		out.Posts[i] = object.PostView{
			Post:        post,
			Ref:         post.R.Hex(),
			AuthorRef:   post.User.Hex(),
			AuthorAlias: "-", // TODO: Implement.
			Votes:       bState.GetPostVotes(skyobject.Reference(post.R)),
		}
	}

	return out
}

type VotesOutput struct {
	Reference string              `json:"reference"`
	Votes     *object.VoteSummary `json:"votes"`
}

func getThreadVotes(
	ctx context.Context, compiler *state.Compiler, result *content.Result, tRef skyobject.Reference,
) *VotesOutput {
	return &VotesOutput{
		Reference: tRef.String(),
		Votes: compiler.
			GetBoard(result.GetPK()).
			GetThreadVotesSeq(ctx, tRef, result.GetSeq()),
	}
}

func getPostVotes(
	ctx context.Context, compiler *state.Compiler, result *content.Result, pRef skyobject.Reference,
) *VotesOutput {
	return &VotesOutput{
		Reference: pRef.String(),
		Votes: compiler.
			GetBoard(result.GetPK()).
			GetPostVotesSeq(ctx, pRef, result.GetSeq()),
	}
}
