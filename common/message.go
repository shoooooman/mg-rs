package common

// Message is a type sent through network pkg
type Message struct {
	SenderID int
	Body     interface{}
}
