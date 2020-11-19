package market

import (
	"reflect"
	"testing"

	"github.com/shoooooman/mg-rs/common"
)

func TestTopGateway_GetTxReq(t *testing.T) {
	reqID = 0
	tg := &TopGateway{}

	ratings := map[int]float64{
		0: 0.1,
		1: 0.3,
		2: 0.2,
		4: 0.5,
	}
	expected := []common.TxReq{
		{
			ID:      0,
			PartyID: 0,
		},
		{
			ID:      1,
			PartyID: 0,
		},
	}
	for i := 0; i < len(expected); i++ {
		req := tg.GetTxReq(ratings)
		if !reflect.DeepEqual(req, expected[i]) {
			t.Errorf("\nexpected: %v\nactual: %v\n", expected[i], req)
		}
	}
}
