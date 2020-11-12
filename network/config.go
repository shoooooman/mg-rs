package network

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Nodes   []Node `mapstructure:"nodes"`
	NodeMap map[int]*Node
}

// Node is ...
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

func (c *config) getPeerIDs(id int) []int {
	if _, ok := c.NodeMap[id]; !ok {
		log.Fatal("wrong id error")
	}
	return c.NodeMap[id].Peers
}

func (c *config) getPeerAddrs(id int) []string {
	if _, ok := c.NodeMap[id]; !ok {
		log.Fatal("wrong id error")
	}
	node := c.NodeMap[id]
	peers := node.Peers
	addrs := make([]string, len(peers))
	for i, p := range peers {
		addrs[i] = c.getAddr(p)
	}
	return addrs
}
