package network

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/shoooooman/mg-rs/common"
)

const bufsize = 5

// ClientImpl is ...
type ClientImpl struct {
	peers []*Peer
	buf   chan *common.Message
	*RPCServer
}

// ConnectPeers is ...
func (c *ClientImpl) ConnectPeers(peers []*Peer) {
	for _, p := range peers {
		client, err := rpc.DialHTTP("tcp", p.Address)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		p.Client = client
	}
	c.peers = peers
}

// RunServer is ...
func (c *ClientImpl) RunServer(addr string) {
	err := rpc.Register(c.RPCServer)
	if err != nil {
		log.Fatal("rpc registration error:", err)
	}
	rpc.HandleHTTP()
	l, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go func() {
		err := http.Serve(l, nil)
		if err != nil {
			log.Fatal("serve error", err)
		}
	}()
}

// Broadcast calls Receive methods of all peers
func (c *ClientImpl) Broadcast(msg *common.Message) {
	for _, p := range c.peers {
		var reply string
		err := p.Client.Call("RPCServer.Receive", msg, &reply)
		if err != nil {
			log.Fatal("calling:", err)
		}
		log.Println("reply:", reply)
	}
}

// GetData receives a message from buf (if empty, waits for receiving)
func (c *ClientImpl) GetData() *common.Message {
	return <-c.buf
}

// GetPeers returns IDs of the peers
func (c *ClientImpl) GetPeers() []int {
	ids := make([]int, len(c.peers))
	for i, p := range c.peers {
		ids[i] = p.ID
	}
	return ids
}

// NewClientImpl is ...
func NewClientImpl(id int) *ClientImpl {
	client := &ClientImpl{}
	client.buf = make(chan *common.Message, bufsize)
	client.RPCServer = &RPCServer{client.buf}

	conf := readConfig()

	addr := conf.getAddr(id)
	client.RunServer(addr)

	time.Sleep(time.Second * 5)

	peers := conf.getPeers(id)
	client.ConnectPeers(peers)

	return client
}
