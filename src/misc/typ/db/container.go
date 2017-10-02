package db

import (
	"sync"
)

type Ref struct {
	Type string
	Key  interface{}
}

type Container struct {
	mux sync.Mutex
}
