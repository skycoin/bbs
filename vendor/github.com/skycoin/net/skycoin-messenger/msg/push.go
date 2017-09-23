package msg

type Reg struct {
	PublicKey string
}

type PushMsg struct {
	From string
	Msg  string
}
