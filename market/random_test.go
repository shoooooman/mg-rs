package market

import (
	"testing"
)

func TestRandomGateway_GetTxReq(t *testing.T) {
	reqID = 0
	ratings := map[int]float64{
		0: 0.0,
		1: 2.0,
		2: 1.0,
		3: -1.0,
	}
	rg := &RandomGateway{}
	const n = 100
	for i := 0; i < n; i++ {
		req := rg.GetTxReq(ratings)
		if req.ID != reqID-1 {
			t.Errorf("\nexpected: %v\nactual: %v\n", reqID-1, req.ID)
		}
		if !(req.PartyID >= 0 && req.PartyID < len(ratings)) {
			t.Errorf("\nexpected: 0<=PartyID<=%v\nactual: %v\n", len(ratings)-1, req.PartyID)
		}
	}
}
