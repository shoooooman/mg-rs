package run

import (
	"reflect"
	"testing"
)

func TestReadConfig(t *testing.T) {
	expected := config{
		Gateway: gatewayConfig{
			Name: "toprand",
			Prob: 0.1,
		},
		Manager: managerConfig{
			Name:  "brs",
			Decay: 0.999,
		},
		Scenario: scenarioConfig{
			Name: "brs_simple",
			N:    1000,
		},
		K: 5,
	}

	confFilename = "config_test"
	v.AddConfigPath(".")
	result := readConfig()

	if !reflect.DeepEqual(*result, expected) {
		t.Errorf("\nexpected: %v\nactual: %v\n", expected, *result)
	}
}
