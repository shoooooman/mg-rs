package monitor

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/shoooooman/mg-rs/common"
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
	confFilename = "config_test"
	v.AddConfigPath(".")
	conf := readConfig()
	if !reflect.DeepEqual(conf.Behaviors, expected) {
		t.Errorf("\nexpected: %v\nactual: %v\n", expected, conf.Behaviors)
	}
}

func TestMockMonitor_MonitorTx(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	mock := &MockMonitor{}
	probs := map[int]func(int) float64{
		0: func(t int) float64 { return 0.0 }, // always true
		1: func(t int) float64 { return 1.0 }, // always false
	}
	mock.probs = probs

	txs := []common.Tx{{ID: 0, Time: 0, PartyID: 0}, {ID: 0, Time: 0, PartyID: 1}}
	expected := []bool{true, false}
	for i, tx := range txs {
		result := mock.MonitorTx(tx)
		if result != expected[i] {
			t.Errorf("\nexpected: %v\nactual: %v\n", expected[i], result)
		}
	}
}

func TestNewMockMonitor(t *testing.T) {
	// mock reading config
	reader = func() mockConfig {
		return mockConfig{expected}
	}
	result := NewMockMonitor()
	testcases := map[int]map[int]float64{
		0:   {1: 0.0, 10: 0.0},
		1:   {1: 0.1, 21: 0.1},
		2:   {1: 0.2, 32: 0.2},
		3:   {1: 0.3, 43: 0.3},
		4:   {1: 0.4, 54: 0.4},
		999: {0: 0.1, 9: 0.1, 10: 0.9, 15: 0.9},
	}

	for id, f := range testcases {
		for time, expected := range f {
			actual := result.probs[id](time)
			if actual != expected {
				t.Errorf("\nexpected: %v\nactual: %v\n", expected, actual)
			}
		}
	}
}
