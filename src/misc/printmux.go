package misc

import (
	"log"
	"reflect"
	"runtime"
	"strings"
	"sync"
)

// Enable determines whether to enable printing.
var Enable = false

type PrintMux struct {
	fName string
	mux   sync.Mutex
}

func (m *PrintMux) Lock(function interface{}) {
	m.mux.Lock()
	if Enable {
		m.fName = runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
		m.fName = strings.Replace(m.fName, "github.com/skycoin/bbs/", "", -1)
		m.fName = strings.Replace(m.fName, "-fm", "", -1)
		log.Printf(">>> [  LOCK] %s <<<", m.fName)
	}
}

func (m *PrintMux) Unlock() {
	if Enable {
		log.Printf("<<< [UNLOCK] %s >>>", m.fName)
	}
	m.mux.Unlock()
}
