package run

import (
	"log"
	"os"

	"github.com/shoooooman/mg-rs/agent"
)

var (
	rlog = log.New(os.Stderr, "[ANALYSIS] ", log.LstdFlags)
)

// Run is ...
func Run(id int) {
	conf := readConfig()
	manager := conf.Manager
	gateway := conf.Gateway
	scenario := conf.Scenario
	n := conf.N
	k := conf.K

	a := agent.NewAgent(id)
	if err := a.SetGateway(gateway); err != nil {
		log.Fatal("SetGateway:", err)
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
		log.Fatal("Run: no such a manager")
	}

	for i := 0; i < k; i++ {
		f(a, n)
	}
}
