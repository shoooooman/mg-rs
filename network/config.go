package network

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Nodes   []Node `mapstructure:"nodes"`
	NodeMap map[int]*Node
}

// Node is a struct for viper unmarshal
type Node struct {
	ID      int    `mapstructure:"id"`
	Address string `mapstructure:"address"`
	Peers   []int  `mapstructure:"peers"`
}

func readConfig() *config {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./network")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("config file error:", err)
	}
	var c config
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Fatal("config unmarshal error:", err)
	}
	setNodeMap(&c)
	return &c
}

// setNodeMap converts a slice of Node to map[ID]*Node and set it to NodeMap
func setNodeMap(conf *config) {
	mp := make(map[int]*Node)
	for _, n := range conf.Nodes {
		cp := n
		mp[cp.ID] = &cp // &nとするとループの度に参照先が変更されてしまう
	}
	conf.NodeMap = mp
}

func (c *config) getAddr(id int) string {
	if _, ok := c.NodeMap[id]; !ok {
		log.Fatal("wrong id error")
	}
	return c.NodeMap[id].Address
}

func (c *config) getPeers(id int) []*Peer {
	if _, ok := c.NodeMap[id]; !ok {
		log.Fatal("wrong id error")
	}
	peerIDs := c.NodeMap[id].Peers
	peers := make([]*Peer, len(peerIDs))
	for i, pid := range peerIDs {
		p := c.NodeMap[pid]
		peers[i] = &Peer{ID: p.ID, Address: p.Address, Client: nil}
	}
	return peers
}

func (c *config) getIDs() []int {
	ids := make([]int, len(c.Nodes))
	for i, n := range c.Nodes {
		ids[i] = n.ID
	}
	return ids
}
