package network

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Nodes []Node `json:"nodes"`
}

// Node is ...
type Node struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
	Peers   []Peer `json:"peers"`
}

// Peer is ...
type Peer struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
}

func readConfig() config {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./network")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("config file error:", err)
	}
	var c config
	viper.Unmarshal(&c)
	return c
}

func getAddr(conf config, id int) string {
	if len(conf.Nodes)-1 < id {
		log.Fatal("wrong id error")
	}
	return conf.Nodes[id].Address
}

func getPeerAddrs(conf config, id int) []string {
	if len(conf.Nodes)-1 < id {
		log.Fatal("wrong id error")
	}
	node := conf.Nodes[id]
	peers := node.Peers
	addrs := make([]string, len(peers))
	for i, p := range peers {
		addrs[i] = p.Address
	}
	return addrs
}
