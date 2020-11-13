package run

import (
	"log"

	"github.com/shoooooman/mg-rs/agent"
)

// Run is ...
func Run(id int) {
	conf := readConfig()
	manager := conf.Manager
	gateway := conf.Gateway
	scenario := conf.Scenario
	n := conf.N

	a := agent.NewAgent(id)
	if err := a.SetGateway(gateway); err != nil {
		log.Fatal("SetGateway:", err)
	}
	if err := a.SetManager(manager); err != nil {
		log.Fatal("SetManager:", err)
	}

	switch scenario {
	case "mock":
		Mock(a, n)
	case "brs_simple":
		Brs(a, n)
	default:
		log.Fatal("Run: no such a manager")
	}
}
