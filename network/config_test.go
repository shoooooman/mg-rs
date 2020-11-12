package network

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

// according to config.json
var expected = []Node{
	{
		ID:      0,
		Address: "127.0.0.1:10000",
		Peers:   []int{1, 2, 3},
	},
	{
		ID:      1,
		Address: "127.0.0.1:10001",
		Peers:   []int{0, 2},
	},
	{
		ID:      2,
		Address: "127.0.0.1:10002",
		Peers:   []int{0, 1, 4},
	},
	{
		ID:      3,
		Address: "127.0.0.1:10003",
		Peers:   []int{0, 4},
	},
	{
		ID:      4,
		Address: "127.0.0.1:10004",
		Peers:   []int{2, 3},
	},
}

var conf *config

func TestReadConfig(t *testing.T) {
	viper.AddConfigPath(".")
	conf = readConfig()

	nodes := conf.Nodes
	if len(nodes) != len(expected) {
		t.Fatal("node length is wrong")
	}

	for i, n := range nodes {
		if !reflect.DeepEqual(n, expected[i]) {
			t.Errorf("\nexpected: %v\nactual: %v\n", expected[i], n)
		}
	}
}

func TestConfig_getAddr(t *testing.T) {
	for _, n := range conf.Nodes {
		addr := conf.getAddr(n.ID)
		if addr != n.Address {
			t.Errorf("\nexpected: %v\nactual: %v\n", n.Address, addr)
		}
	}
}

func TestConfig_getPeers(t *testing.T) {
	for _, n := range conf.Nodes {
		peers := conf.getPeers(n.ID)
		ids := expected[n.ID].Peers
		for j, p := range peers {
			ex := expected[ids[j]]
			if p.ID != ex.ID {
				t.Errorf("\nexpected: %v\nactual: %v\n", ex.ID, p.ID)
			}
			if p.Address != ex.Address {
				t.Errorf("\nexpected: %v\nactual: %v\n", ex.Address, p.Address)
			}
			if p.Client != nil {
				t.Errorf("\nexpected: nil\nactual: %v\n", p.Client)
			}
		}
	}
}

func TestConfig_getIDs(t *testing.T) {
	result := conf.getIDs()
	ids := make([]int, len(expected))
	for i := 0; i < len(expected); i++ {
		ids[i] = expected[i].ID
	}
	if !reflect.DeepEqual(result, ids) {
		t.Errorf("\nexpected: %v\nactual: %v\n", ids, result)
	}
}
