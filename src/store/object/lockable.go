package object

type Lockable interface {
	Lock()
	Unlock()
}

func Lock(obj Lockable) func() {
	obj.Lock()
	return obj.Unlock
}
