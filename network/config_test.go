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
		Peers: []Peer{
			{
				ID:      1,
				Address: "127.0.0.1:10001",
			},
		},
	},
	{
		ID:      1,
		Address: "127.0.0.1:10001",
		Peers: []Peer{
			{
				ID:      0,
				Address: "127.0.0.1:10000",
			},
		},
	},
}

func TestReadConfig(t *testing.T) {
	viper.AddConfigPath(".")
	conf := readConfig()

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

func TestGetAddr(t *testing.T) {
	conf := readConfig()
	for i, n := range conf.Nodes {
		addr := getAddr(conf, i)
		if addr != n.Address {
			t.Errorf("\nexpected: %v\nactual: %v\n", n.Address[i], addr)
		}
	}
}

func TestGetPeerAddrs(t *testing.T) {
	conf := readConfig()
	for i, n := range conf.Nodes {
		peers := getPeerAddrs(conf, i)
		for j, p := range peers {
			if p != n.Peers[j].Address {
				t.Errorf("\nexpected: %v\nactual: %v\n", n.Peers[j].Address, p)
			}
		}
	}
}
