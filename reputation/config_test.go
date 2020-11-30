package reputation

import "testing"

func TestReadConfig(t *testing.T) {
	expected := []struct {
		id int
		ty RType
	}{
		{0, honest},
		{1, reverse},
	}
	confFilename = "config_test"
	v.AddConfigPath(".")
	conf := readConfig()

	types := conf.TypeMap
	if len(types) != len(expected) {
		t.Fatal("node length is wrong")
	}

	for _, exp := range expected {
		if types[exp.id] != exp.ty {
			t.Errorf("\nexpected: %v\nactual: %v\n", exp.ty, types[exp.id])
		}
	}
}
