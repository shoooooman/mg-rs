package monitor

import (
	"log"
	"math/rand"
	"time"

	"github.com/shoooooman/mg-rs/common"
	"github.com/spf13/viper"
)

// MockMonitor is ...
type MockMonitor struct {
	probs map[int]func(int) float64 // each probability func mapped with id
}

type mockConfig struct {
	Behaviors []Behavior `json:"behaviors"`
}

// Behavior is ...
type Behavior struct {
	ID    int     `mapstructure:"id"`
	Kind  string  `mapstructure:"kind"`
	Prob  float64 `mapstructure:"probability"`
	VProb []*VP   `mapstructure:"var_probs"`
}

// VP is a probility function
type VP struct {
	L    int     `mapstructure:"left"`
	R    int     `mapstructure:"right"`
	Prob float64 `mapstructure:"value"`
}

var (
	v            = viper.New()
	reader       = readConfig
	confFilename = "config"
)

func readConfig() mockConfig {
	v.SetConfigName(confFilename)
	v.SetConfigType("json")
	v.AddConfigPath("./monitor")
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("config file error:", err)
	}
	var c mockConfig
	err = v.Unmarshal(&c)
	if err != nil {
		log.Fatal("config unmarshal error:", err)
	}
	return c
}

// MonitorTx is ...
func (m *MockMonitor) MonitorTx(tx common.Tx) bool {
	time := tx.Time
	partyID := tx.PartyID
	f := m.probs[partyID]
	val := int(f(time) * 100)
	r := rand.Intn(100)
	return r >= val
}

// NewMockMonitor is ...
func NewMockMonitor() *MockMonitor {
	rand.Seed(time.Now().UnixNano())
	mock := &MockMonitor{}

	conf := reader()
	probs := make(map[int]func(int) float64)
	for _, b := range conf.Behaviors {
		id := b.ID
		kind := b.Kind
		prob := b.Prob
		if kind == "fixed" {
			probs[id] = func(t int) float64 {
				return prob
			}
		} else if kind == "variable" {
			probs[id] = func(t int) float64 {
				for _, vprob := range b.VProb {
					if t >= vprob.L && t < vprob.R {
						return vprob.Prob
					}
				}
				log.Fatal("variable probability is not set for this time")
				return 0.0
			}
		} else {
			log.Fatal("behavior kind is error")
		}
	}
	mock.probs = probs
	return mock
}
