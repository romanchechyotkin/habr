package message

type Message struct {
	Header uint8
	// Pub - Sub
	// 1 0 - 1 0
	// y n   y n
}

func New() *Message {
	msg := &Message{}
	return msg
}
