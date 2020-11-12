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
	peers []*rpc.Client
	buf   chan *common.Message
	*RPCServer
}

// ConnectPeers is ...
func (c *ClientImpl) ConnectPeers(peerAddrs []string) {
	peers := make([]*rpc.Client, len(peerAddrs))
	for i, addr := range peerAddrs {
		peer, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		peers[i] = peer
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
		err := p.Call("RPCServer.Receive", msg, &reply)
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

// NewClientImpl is ...
func NewClientImpl(id int) *ClientImpl {
	client := &ClientImpl{}
	client.buf = make(chan *common.Message, bufsize)
	client.RPCServer = &RPCServer{client.buf}

	conf := readConfig()

	addr := conf.getAddr(id)
	client.RunServer(addr)

	time.Sleep(time.Second * 5)

	peerAddrs := conf.getPeerAddrs(id)
	client.ConnectPeers(peerAddrs)

	return client
}
