package agent

import (
	"fmt"

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

// SetGateway is ...
func (a *Agent) SetGateway(gw string) error {
	switch gw {
	case "top":
		a.Gateway = &market.TopGateway{}
	case "toprand":
		a.Gateway = &market.TopRandGateway{}
	default:
		return fmt.Errorf("no such a gateway")
	}
	return nil
}

// SetManager is ...
func (a *Agent) SetManager(rm string) error {
	switch rm {
	case "mock":
		a.Manager = reputation.NewMockManager(a.ID)
	case "brs":
		a.Manager = reputation.NewBrs(a.ID)
	default:
		return fmt.Errorf("no such a reputation manager")
	}
	return nil
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
