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
		Peers:   []int{1},
	},
	{
		ID:      1,
		Address: "127.0.0.1:10001",
		Peers:   []int{0},
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

func TestConfig_getPeerIDs(t *testing.T) {
	for _, n := range conf.Nodes {
		peerIDs := conf.getPeerIDs(n.ID)
		if !reflect.DeepEqual(peerIDs, conf.NodeMap[n.ID].Peers) {
			t.Errorf("\nexpected: %v\nactual: %v\n", conf.NodeMap[n.ID].Peers, peerIDs)
		}
	}
}

func TestConfig_getPeerAddrs(t *testing.T) {
	for _, n := range conf.Nodes {
		paddrs := conf.getPeerAddrs(n.ID)
		nd := conf.NodeMap[n.ID]
		for j, paddr := range paddrs {
			peerID := nd.Peers[j]
			if paddr != conf.NodeMap[peerID].Address {
				t.Errorf("\nexpected: %v\nactual: %v\n", conf.NodeMap[peerID].Address, paddr)
			}
		}
	}
}
