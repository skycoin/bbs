package typ

import "sync"

type Bool struct {
	value bool
	mux   sync.RWMutex
}

func (b *Bool) Set() {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.value = true
}

func (b *Bool) Clear() {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.value = false
}

func (b *Bool) Value() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.value
}
