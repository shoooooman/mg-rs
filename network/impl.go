package network

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/shoooooman/mg-rs/common"
)

const bufsize = 1000

var done = make(chan []string)

// ClientImpl is ...
type ClientImpl struct {
	id      int
	nodeIDs []int
	peers   []*Peer
	buf     chan *common.Message
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

// GetIDs returns the all IDs of the network
func (c *ClientImpl) GetIDs() []int {
	return c.nodeIDs
}

// Register sends its address to the master for registration
func (c *ClientImpl) Register(masterAddr, addr string) {
	master, err := rpc.DialHTTP("tcp", masterAddr)
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	msg := &common.Message{SenderID: c.id, Body: addr}
	err = master.Call("RPCServer.Register", msg, &reply)
	if err != nil {
		log.Fatal("calling:", err)
	}
	log.Printf("%d: reply to registration: %v\n", c.id, reply)
}

// NewClientImpl is ...
func NewClientImpl(id int) *ClientImpl {
	client := &ClientImpl{id: id}
	client.buf = make(chan *common.Message, bufsize)

	conf := readConfig()
	client.nodeIDs = conf.getIDs()

	client.RPCServer = &RPCServer{buf: client.buf}

	addr := conf.getAddr(id)
	client.RunServer(addr)
	log.Printf("%d started serving\n", id)

	client.Register(conf.getMasterAddr(), addr)
	msg := client.GetData() // wait for receiving done message from the master
	log.Printf("%d: msg right after registration: %v\n", id, msg)
	if msg.SenderID != masterID {
		log.Fatalf("%d: registration error", id)
	}
	log.Printf("%d: initialization is done\n", id)
	peers := conf.getPeers(id)
	client.ConnectPeers(peers)

	time.Sleep(time.Second * 3)

	return client
}

// RunMasterClient starts the master, which synchronizes agents before running experiment
func RunMasterClient() {
	client := &ClientImpl{id: masterID}

	conf := readConfig()
	client.nodeIDs = conf.getIDs()

	client.RPCServer = &RPCServer{num: len(client.nodeIDs)}

	addr := conf.getAddr(masterID)
	client.RunServer(addr)
	log.Printf("%d started serving\n", masterID)

	addrs := <-done // wait for all nodes to finish registration
	log.Println("master: all nodes are registered")
	for _, addr := range addrs {
		c, err := rpc.DialHTTP("tcp", addr)
		if err != nil {
			log.Fatal("dialing:", err)
		}
		msg := &common.Message{SenderID: masterID, Body: "done"}
		var reply string
		err = c.Call("RPCServer.Receive", msg, &reply)
		if err != nil {
			log.Fatal("calling:", err)
		}
		log.Printf("%d: reply to done msg: %v\n", masterID, reply)
	}
}
