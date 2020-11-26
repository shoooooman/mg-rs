package run

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Gateway  gatewayConfig  `mapstructure:"gateway"`
	Manager  managerConfig  `mapstructure:"reputation_manager"`
	Scenario scenarioConfig `mapstructure:"scenario"`
	K        int            `mapstructure:"run_num"`
}

type gatewayConfig struct {
	Name string  `mapstructure:"name"`
	Prob float64 `mapstructure:"random_prob"`
}

type managerConfig struct {
	Name  string  `mapstructure:"name"`
	Decay float64 `mapstructure:"decay_factor"`
}

type scenarioConfig struct {
	Name string `mapstructure:"name"`
	N    int    `mapstructure:"tx_num"`
}

var (
	v            = viper.New()
	confFilename = "config"
)

const (
	decayDefault = 1.0
)

func readConfig() *config {
	v.SetConfigName(confFilename)
	v.SetConfigType("json")
	v.AddConfigPath("./run")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("config file error:", err)
	}
	var c config
	c.Manager.Decay = decayDefault
	err = v.Unmarshal(&c)
	if err != nil {
		log.Fatal("config unmarshal error:", err)
	}
	return &c
}
