package network

import (
	"fmt"
	"log"
	"sync"

	"github.com/shoooooman/mg-rs/common"
)

var (
	mu sync.Mutex
)

// RPCServer is ...
type RPCServer struct {
	buf   *chan *common.Message
	addrs []string
	num   int
}

// Receive adds the received data to buf
// The type of the first argument must be the actual received one (cannot be interface{})
func (rc *RPCServer) Receive(msg *common.Message, reply *string) error {
	log.Println("received:", *msg)
	*rc.buf <- msg
	*reply = "success"
	return nil
}

// Register is ...
func (rc *RPCServer) Register(msg *common.Message, reply *string) error {
	addr, ok := msg.Body.(string)
	if !ok {
		return fmt.Errorf("Register: the message body should be an address")
	}
	mu.Lock()
	rc.addrs = append(rc.addrs, addr)
	mu.Unlock()
	if len(rc.addrs) == rc.num {
		done <- rc.addrs
	}
	*reply = fmt.Sprintf("registration of %s is complete", addr)
	return nil
}
