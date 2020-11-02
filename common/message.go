package common

// TxReq is ...
type TxReq struct {
	ID    int
	Party int
}

// Tx is ...
type Tx struct {
	ID    int
	Party int
}

// Message is a type sent through network pkg
type Message struct {
	SenderID int
	Body     interface{}
}
