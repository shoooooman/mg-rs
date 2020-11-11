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
	probs map[int]float64
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

func readConfig() mockConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./monitor")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("config file error:", err)
	}
	var c mockConfig
	viper.Unmarshal(&c)
	return c
}

// MonitorTx is ...
func (m *MockMonitor) MonitorTx(tx common.Tx) bool {
	partyID := tx.PartyID
	r := rand.Intn(100)
	val := int(m.probs[partyID] * 100)
	log.Printf("(r, val)=(%v, %v)\n", r, val)
	if r < val {
		return false
	}
	return true
}

// NewMockMonitor is ...
func NewMockMonitor() *MockMonitor {
	rand.Seed(time.Now().UnixNano())
	mock := &MockMonitor{}

	probs := make(map[int]float64)
	conf := readConfig()
	for _, b := range conf.Behaviors {
		id := b.ID
		kind := b.Kind
		if kind == "fixed" {
			probs[id] = b.Prob
		} else if kind == "variable" {
			// TODO: variable probを実装する
		} else {
			log.Fatal("behavior kind is error")
		}
	}
	mock.probs = probs
	return mock
}
