package local

type Message struct {
	Username string // empty username means server's message
	Content  string // content of the message
}
