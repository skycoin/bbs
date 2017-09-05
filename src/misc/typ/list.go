package typ

import (
	"github.com/skycoin/bbs/src/misc/boo"
	"log"
)

type List struct {
	keys []interface{}
	dist map[interface{}]interface{}
}

func NewList() *List {
	return &List{
		dist: make(map[interface{}]interface{}),
	}
}

func (l *List) Len() int {
	return len(l.keys)
}

func (l *List) Append(key interface{}, value interface{}) bool {
	if _, has := l.dist[key]; has {
		return false
	}
	l.keys = append(l.keys, key)
	l.dist[key] = value
	return true
}

func (l *List) DelOfIndex(i int) bool {
	if len(l.keys) <= i {
		return false
	}
	key := l.keys[i]
	delete(l.dist, key)
	return true
}

func (l *List) DelOfKey(key interface{}) bool {
	if _, has := l.dist[key]; has == false {
		return false
	}
	for i, oldKey := range l.keys {
		if oldKey == key {
			l.keys = append(l.keys[:i], l.keys[i+1:]...)
			delete(l.dist, key)
			return true
		}
	}
	log.Panicf("key %v not in l.keys", key)
	return false
}

func (l *List) GetOfKey(key interface{}) (interface{}, bool) {
	value, has := l.dist[key]
	return value, has
}

func (l *List) GetOfIndex(index int) (interface{}, bool) {
	if len(l.keys) <= index {
		return nil, false
	}
	return l.dist[l.keys[index]], true
}

func (l *List) Keys() []interface{} {
	return l.keys
}

func (l *List) Values() []interface{} {
	out := make([]interface{}, len(l.keys))
	for i, key := range l.keys {
		out[i] = l.dist[key]
	}
	return out
}

func (l *List) HasKey(key interface{}) bool {
	_, has := l.dist[key]
	return has
}

type RangeOrder int

const (
	Ascending RangeOrder = 1 << iota
	Descending
)

type RangeFunc func(index int, key, value interface{}) (stop bool, e error)

func (l *List) Range(order RangeOrder, action RangeFunc) error {
	switch order {
	case Ascending:
		for i := 0; i < len(l.keys); i++ {
			var key = l.keys[i]
			if stop, e := action(i, key, l.dist[key]); e != nil {
				return e
			} else if stop {
				return nil
			}
		}
	case Descending:
		for i := len(l.keys) - 1; i >= 0; i-- {
			var key = l.keys[i]
			if stop, e := action(i, key, l.dist[key]); e != nil {
				return e
			} else if stop {
				return nil
			}
		}
	default:
		return boo.New(boo.Internal,
			"invalid 'RangeOrder' specified")
	}
	return nil
}
