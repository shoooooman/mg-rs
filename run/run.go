package run

import (
	"log"
	"os"

	"github.com/shoooooman/mg-rs/agent"
	"github.com/shoooooman/mg-rs/market"
	"github.com/shoooooman/mg-rs/network"
)

var (
	rlog = log.New(os.Stderr, "[ANALYSIS] ", log.LstdFlags)
)

// Master is ...
func Master() {
	network.RunMasterClient()
}

// Run is ...
func Run(id int) {
	conf := readConfig()
	var (
		gateway  = conf.Gateway.Name
		manager  = conf.Manager
		scenario = conf.Scenario.Name
		n        = conf.Scenario.N
		k        = conf.K
	)

	a := agent.NewAgent(id)
	if err := a.SetGateway(gateway); err != nil {
		log.Fatal("SetGateway:", err)
	}
	if gw, ok := a.Gateway.(*market.TopRandGateway); ok {
		gw.SetRandProb(conf.Gateway.Prob)
	}

	if err := a.SetManager(manager); err != nil {
		log.Fatal("SetManager:", err)
	}

	var f func(*agent.Agent, int)
	switch scenario {
	case "mock":
		f = Mock
	case "brs_simple":
		f = Brs
	case "bdf_simple":
		f = Bdf
	default:
		log.Fatal("Run: no such a scenario")
	}

	for i := 0; i < k; i++ {
		f(a, n)
	}
}
