package gui

//import (
//	"github.com/skycoin/cxo/skyobject"
//	"github.com/skycoin/bbs/src/store/cxo"
//	"github.com/pkg/errors"
//	"github.com/skycoin/skycoin/src/cipher"
//	"github.com/skycoin/bbs/src/misc"
//	"net/http"
//)
//
//func (a *API) GetThreadPageAsHex(w http.ResponseWriter, r *http.Request) {
//	// Get board public key.
//	bpk, e := misc.GetPubKey(r.FormValue("board"))
//	if e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	// Get thread reference.
//	tRef, e := misc.GetReference(r.FormValue("thread"))
//	if e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	// Get thread page as hex.
//	tph, e := a.g.GetThreadPageAsHex(bpk, tRef)
//	if e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	sendResponse(w, *tph, http.StatusOK)
//}
//
//func (a *API) GetThreadPageWithTpRefAsHex(w http.ResponseWriter, r *http.Request) {
//	// Get thread page reference.
//	tpRef, e := misc.GetReference(r.FormValue("threadpage"))
//	if e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	// Get thread page as hex.
//	tph, e := a.g.GetThreadPageWithTpRefAsHex(tpRef)
//	if e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	sendResponse(w, *tph, http.StatusOK)
//}
//
//func (a *API) NewThreadWithHex(w http.ResponseWriter, r *http.Request) {
//	// Get board public key.
//	bpk, e := misc.GetPubKey(r.FormValue("board"))
//	if e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	// Get thread data.
//	tData, e := misc.GetBytes(r.FormValue("raw_thread"))
//	if e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	// Inject.
//	if e := a.g.NewThreadWithHex(bpk, tData); e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	sendResponse(w, true, http.StatusOK)
//}
//
//func (a *API) NewPostWithHex(w http.ResponseWriter, r *http.Request) {
//	// Get board public key.
//	bpk, e := misc.GetPubKey(r.FormValue("board"))
//	if e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	// Get thread reference.
//	tRef, e := misc.GetReference(r.FormValue("thread"))
//	if e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	// Get request data.
//	pData, e := misc.GetBytes(r.FormValue("raw_post"))
//	if e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	// Inject.
//	if e := a.g.NewPostWithHex(bpk, tRef, pData); e != nil {
//		sendResponse(w, e.Error(), http.StatusBadRequest)
//		return
//	}
//	sendResponse(w, true, http.StatusOK)
//}
//
//func (g *Gateway) GetThreadPageAsHex(bpk cipher.PubKey, tRef skyobject.Reference) (*cxo.ThreadPageHex, error) {
//	return g.container.GetThreadPageAsHex(bpk, tRef)
//}
//
//func (g *Gateway) GetThreadPageWithTpRefAsHex(tpRef skyobject.Reference) (*cxo.ThreadPageHex, error) {
//	return g.container.GetThreadPageWithTpRefAsHex(tpRef)
//}
//
//func (g *Gateway) NewThreadWithHex(bpk cipher.PubKey, tData []byte) error {
//	bi, has := g.boardSaver.Get(bpk)
//	if !has {
//		return errors.New("not subscribed to board")
//	}
//	if !bi.Config.Master {
//		return errors.New("not master of board")
//	}
//	return g.container.NewThreadWithHex(bpk, bi.Config.GetSK(), tData)
//}
//
//func (g *Gateway) NewPostWithHex(bpk cipher.PubKey, tRef skyobject.Reference, pData []byte) error {
//	bi, has := g.boardSaver.Get(bpk)
//	if !has {
//		return errors.New("not subscribed to board")
//	}
//	if !bi.Config.Master {
//		return errors.New("not master of board")
//	}
//	return g.container.NewPostWithHex(bpk, bi.Config.GetSK(), tRef, pData)
//}
