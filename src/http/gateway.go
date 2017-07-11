package http

import (
	"encoding/json"
	"github.com/skycoin/bbs/src/access/btp"
	"github.com/skycoin/bbs/src/boo"
	"net/http"
)

// Gateway represents what is exposed to HTTP interface.
type Gateway struct {
	boardAccessor *btp.BoardAccessor
}

// NewGateway creates a new Gateway.
func NewGateway(boardAccessor *btp.BoardAccessor) *Gateway {
	return &Gateway{
		boardAccessor: boardAccessor,
	}
}

func (g *Gateway) prepare(mux *http.ServeMux) error {
	mux.HandleFunc("/api/boards/new", g.BoardsNew)
	return nil
}

// NewBoard creates a new board.
func (g *Gateway) BoardsNew(w http.ResponseWriter, r *http.Request) {
	view, e := g.boardAccessor.NewBoard(&btp.NewBoardInput{
		Name: r.FormValue("name"),
		Desc: r.FormValue("description"),
		Seed: r.FormValue("seed"),
	})
	if e != nil {
		sendErr(w, e)
		return
	}
	sendOK(w, *view)
}

/*
	<<< HELPER FUNCTIONS >>>
*/

type Error struct {
	Type    boo.Type `json:"type"`
	Message string   `json:"message"`
	Details string   `json:"details"`
}

type Response struct {
	Okay  bool        `json:"okay"`
	Data  interface{} `json:"data,omitempty"`
	Error *Error      `json:"error,omitempty"`
}

func sendOK(w http.ResponseWriter, v interface{}) error {
	response := Response{Okay: true, Data: v}
	return send(w, response, http.StatusOK)
}

func sendErr(w http.ResponseWriter, e error) error {
	eType := boo.What(e)
	eMsg := boo.Message(eType)
	var status int
	switch eType {
	case boo.Unknown, boo.Internal:
		status = http.StatusInternalServerError
	case boo.NotAuthorised, boo.NotMaster:
		status = http.StatusUnauthorized
	case boo.ObjectNotFound:
		status = http.StatusNotFound
	case boo.ObjectAlreadyExists:
		status = http.StatusConflict
	default:
		status = http.StatusBadRequest
	}

	response := Response{
		Okay: false,
		Error: &Error{
			Type:    eType,
			Message: eMsg,
			Details: e.Error(),
		},
	}
	return send(w, response, status)
}

func send(w http.ResponseWriter, v interface{}, httpStatus int) error {
	w.Header().Set("Content-Type", "application/json")
	respData, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.WriteHeader(httpStatus)
	w.Write(respData)
	return nil
}
