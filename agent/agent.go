package agent

import (
	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/market"
	"github.com/shoooooman/mg-rs/monitor"
	"github.com/shoooooman/mg-rs/reputation"
)

// Agent is ...
type Agent struct {
	ID int
	market.Gateway
	monitor.Monitor
	reputation.Manager
	TxReqs    map[int]common.TxReq
	TxHistory map[int]common.Tx
}

// NewAgent is ...
func NewAgent(id int) *Agent {
	a := &Agent{
		ID:        id,
		Gateway:   nil,
		Monitor:   monitor.NewMockMonitor(),
		Manager:   nil,
		TxReqs:    make(map[int]common.TxReq),
		TxHistory: make(map[int]common.Tx),
	}

	return a
}
