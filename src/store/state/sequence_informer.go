package state

type SequenceInformer struct {
	Seq     uint64
	Trigger func()
}

func (si *SequenceInformer) Join(other *SequenceInformer) *SequenceInformer {
	si.Trigger = func() {
		si.Trigger()
		other.Trigger()
	}
	return si
}
