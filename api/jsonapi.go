package api

import (
	json "encoding/json"
	"github.com/evanlinjin/bbs/cxo"
	"io/ioutil"
	"net/http"
)

// JsonAPI wraps cxo.Gateway.
type JsonAPI struct {
	g *cxo.Gateway
}

// New JsonAPI creates a new JsonAPI.
func NewJsonAPI(g *cxo.Gateway) *JsonAPI {
	return &JsonAPI{g}
}

// BoardHandler for /api/boards.
func (a *JsonAPI) BoardsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		reply := a.g.ViewBoards()
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

func readRequestBody(r *http.Request) (*cxo.JsonApiObj, error) {
	d, e := ioutil.ReadAll(r.Body)
	if e != nil {
		return nil, e
	}
	obj := cxo.NewJsonApiObj()
	e = json.Unmarshal(d, obj)
	return obj, e
}
