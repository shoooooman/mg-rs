package network

import (
	"github.com/shoooooman/mg-rs/common"
)

// Client is ...
type Client interface {
	ConnectPeers([]string)
	RunServer(string)
	Broadcast(*common.Message)
	Receive(*common.Message, *string) error
	GetData() *common.Message
}
