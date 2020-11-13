package run

import (
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Manager  string `mapstructure:"reputation_manager"`
	Gateway  string `mapstructure:"gateway"`
	Scenario string `mapstructure:"scenario"`
	N        int    `mapstructure:"tx_num"`
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
