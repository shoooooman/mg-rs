package network

import (
	"log"

	"github.com/shoooooman/mg-rs/common"
)

// RPCServer is ...
type RPCServer struct {
	buf chan *common.Message
}

// Receive adds the received data to buf
// The type of the first argument must be the actual received one (cannot be interface{})
func (rc *RPCServer) Receive(msg *common.Message, reply *string) error {
	log.Println("received:", *msg)
	rc.buf <- msg
	*reply = "success"
	return nil
}
