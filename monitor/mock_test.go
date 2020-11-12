package monitor

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/shoooooman/mg-rs/common"
	"github.com/spf13/viper"
)

var expected = []Behavior{
	{
		ID:    0,
		Kind:  "fixed",
		Prob:  0.0,
		VProb: nil,
	},
	{
		ID:    1,
		Kind:  "fixed",
		Prob:  0.1,
		VProb: nil,
	},
	{
		ID:    2,
		Kind:  "fixed",
		Prob:  0.2,
		VProb: nil,
	},
	{
		ID:    3,
		Kind:  "fixed",
		Prob:  0.3,
		VProb: nil,
	},
	{
		ID:    4,
		Kind:  "fixed",
		Prob:  0.4,
		VProb: nil,
	},
	{
		ID:   999,
		Kind: "variable",
		Prob: 0.0,
		VProb: []*VP{
			&VP{0, 10, 0.1},
			&VP{10, 20, 0.9},
		},
	},
}

func TestReadConfig(t *testing.T) {
	viper.AddConfigPath(".")
	conf := readConfig()
	if !reflect.DeepEqual(conf.Behaviors, expected) {
		t.Errorf("\nexpected: %v\nactual: %v\n", expected, conf.Behaviors)
	}
}

func TestMockMonitor_MonitorTx(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	mock := &MockMonitor{}
	probs := map[int]float64{
		0: 0.0, // always true
		1: 1.0, // always false
	}
	mock.probs = probs

	txs := []common.Tx{{ID: 0, PartyID: 0}, {ID: 0, PartyID: 1}}
	expected := []bool{true, false}
	for i, tx := range txs {
		result := mock.MonitorTx(tx)
		if result != expected[i] {
			t.Errorf("\nexpected: %v\nactual: %v\n", expected[i], result)
		}
	}
}
