package run

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Gateway  gatewayConfig  `mapstructure:"gateway"`
	Manager  string         `mapstructure:"reputation_manager"`
	Scenario scenarioConfig `mapstructure:"scenario"`
	K        int            `mapstructure:"run_num"`
}

type gatewayConfig struct {
	Name string  `mapstructure:"name"`
	Prob float64 `mapstructure:"random_prob"`
}

type scenarioConfig struct {
	Name string `mapstructure:"name"`
	N    int    `mapstructure:"tx_num"`
}

var v = viper.New()

func readConfig() *config {
	v.SetConfigName("config")
	v.SetConfigType("json")
	v.AddConfigPath("./run")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("config file error:", err)
	}
	var c config
	err = v.Unmarshal(&c)
	if err != nil {
		log.Fatal("config unmarshal error:", err)
	}
	return &c
}
