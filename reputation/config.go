package reputation

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type config struct {
	Types   []ReportType `json:"types"`
	TypeMap map[int]RType
}

// ReportType is ...
type ReportType struct {
	ID   int    `mapstructure:"id"`
	Type string `mapstructure:"type"`
}

var (
	v            = viper.New()
	reader       = readConfig
	confFilename = "config"
)

func readConfig() config {
	v.SetConfigName(confFilename)
	v.SetConfigType("json")
	v.AddConfigPath("./reputation")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("config file error:", err)
	}
	var c config
	err = v.Unmarshal(&c)
	if err != nil {
		log.Fatal("config unmarshal error:", err)
	}
	setTypeMap(&c)
	return c
}

func setTypeMap(conf *config) {
	mp := make(map[int]RType)
	for _, t := range conf.Types {
		if ty, err := convertRType(t.Type); err != nil {
			log.Fatal(err)
		} else {
			mp[t.ID] = ty
		}
	}
	conf.TypeMap = mp
}

func convertRType(str string) (RType, error) {
	switch str {
	case "honest":
		return honest, nil
	case "reverse":
		return reverse, nil
	default:
		return -1, fmt.Errorf("no such a report type")
	}
}

func (c *config) getReportType(id int) RType {
	return c.TypeMap[id]
}
