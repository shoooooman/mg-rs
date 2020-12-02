package run

import (
	"log"
	"os"

	"github.com/shoooooman/mg-rs/agent"
	"github.com/shoooooman/mg-rs/market"
	"github.com/shoooooman/mg-rs/network"
	"github.com/shoooooman/mg-rs/reputation"
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
		manager  = conf.Manager.Name
		decay    = conf.Manager.Decay
		scenario = conf.Scenario.Name
		n        = conf.Scenario.N
		k        = conf.K
	)

	a := agent.NewAgent(id)

	switch gateway {
	case "random":
		a.Gateway = &market.RandomGateway{}
	case "top":
		a.Gateway = &market.TopGateway{}
	case "toprand":
		a.Gateway = &market.TopRandGateway{}
		a.Gateway.(*market.TopRandGateway).SetRandProb(conf.Gateway.Prob)
	default:
		log.Fatal("no such a gateway")
	}

	switch manager {
	case "mock":
		a.Manager = reputation.NewMockManager(a.ID)
	case "brs":
		a.Manager = reputation.NewBrs(a.ID, decay)
	case "bdf":
		a.Manager = reputation.NewBdf(a.ID, decay)
	case "bdfv":
		a.Manager = reputation.NewBdfv(a.ID, decay)
	case "bvf":
		a.Manager = reputation.NewBvf(a.ID, decay)
	default:
		log.Fatal("no such a reputation manager")
	}

	var f func(*agent.Agent, int)
	switch scenario {
	case "mock":
		f = Mock
	case "brs_simple":
		f = Brs
	case "bdf_simple":
		f = Bdf
	case "bvf_simple":
		f = Bvf
	default:
		log.Fatal("Run: no such a scenario")
	}

	for i := 0; i < k; i++ {
		f(a, n)
	}
}
