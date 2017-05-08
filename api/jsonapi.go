package api

import (
	json "encoding/json"
	"github.com/evanlinjin/bbs/cxo"
	"io/ioutil"
	"net/http"
	"strings"
)

// JsonAPI wraps cxo.Gateway.
type JsonAPI struct {
	g *cxo.Gateway
}

// New JsonAPI creates a new JsonAPI.
func NewJsonAPI(g *cxo.Gateway) *JsonAPI {
	return &JsonAPI{g}
}

// BoardListHandler for /api/boards.
func (a *JsonAPI) BoardListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		reply := a.g.ListBoards()
		sendResponse(w, reply, http.StatusOK)
		return
	case "PUT":
		req, e := readRequestBody(r)
		if e != nil || req.Board == nil {
			sendResponse(w, "invalid request body", http.StatusNotAcceptable)
			return
		}
		reply := a.g.NewBoard(req.Board, req.Seed)
		sendResponse(w, reply, http.StatusOK)
		return
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// BoardHandler for /api/boards/BOARD_PUBLIC_KEY.
func (a *JsonAPI) BoardHandler(w http.ResponseWriter, r *http.Request) {
	// Obtain path.
	path := strings.Split(r.URL.EscapedPath(), "/")
	// Obtain public key.
	pkStr := path[3]
	// If it's view board, or view thread.
	switch len(path) {
	case 4:
		// View Board.
		switch r.Method {
		case "GET":
			reply := a.g.ViewBoard(pkStr)
			sendResponse(w, reply, http.StatusOK)
			return
		case "PUT":
			req, e := readRequestBody(r)
			if e != nil || req.Thread == nil {
				sendResponse(w, "invalid request body", http.StatusNotAcceptable)
				return
			}
			reply := a.g.NewThread(pkStr, req.Thread)
			sendResponse(w, reply, http.StatusOK)
			return
		}
	case 5:
		// View Thread.
		tHashStr := path[4]
		switch r.Method {
		case "GET":
			reply := a.g.ViewThread(pkStr, tHashStr)
			sendResponse(w, reply, http.StatusOK)
			return
		case "PUT":
			req, e := readRequestBody(r)
			if e != nil || req.Post == nil {
				sendResponse(w, "invalid request body", http.StatusNotAcceptable)
				return
			}
			reply := a.g.NewPost(pkStr, tHashStr, req.Post)
			sendResponse(w, reply, http.StatusOK)
			return
		}
		//sendResponse(w, tHashStr, http.StatusNotImplemented)
	}
	sendResponse(w, nil, http.StatusNotFound)
	return
}

// Helper functions.
func sendResponse(w http.ResponseWriter, v interface{}, httpStatus int) error {
	w.Header().Set("Content-Type", "application/json")
	respData, err := json.Marshal(v)
	if err != nil {
		return err
	}
	w.WriteHeader(httpStatus)
	w.Write(respData)
	return nil
}

func readRequestBody(r *http.Request) (*cxo.RepReq, error) {
	d, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, e
	}
	obj := cxo.NewRepReq()
	e = json.Unmarshal(d, obj)
	return obj, e
}
