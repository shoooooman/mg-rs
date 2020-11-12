package network

import (
	"net/rpc"

	"github.com/shoooooman/mg-rs/common"
)

// Client is ...
type Client interface {
	ConnectPeers([]*Peer)
	RunServer(string)
	Broadcast(*common.Message)
	Receive(*common.Message, *string) error
	GetData() *common.Message
}

// Peer is ...
type Peer struct {
	ID      int
	Address string
	Client  *rpc.Client
}
