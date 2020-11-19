package market

import (
	"math/rand"

	"github.com/shoooooman/mg-rs/common"
)

// TopRandGateway is ...
type TopRandGateway struct {
	randP float64 // The probability that random policy is chosen
}

// GetTxReq is ...
func (gw *TopRandGateway) GetTxReq(ratings map[int]float64) common.TxReq {
	r := rand.Intn(100)
	// choose the party randomly
	if r < int(gw.randP*100) {
		rth := rand.Intn(len(ratings))
		cnt := 0
		var pid int
		// choose the r th agent in the map as a party
		for id := range ratings {
			if cnt == rth {
				pid = id
				break
			}
			cnt++
		}
		req := common.TxReq{
			ID:      reqID,
			PartyID: pid,
		}
		reqID++
		return req
	}

	min := 1.1
	minID := -1
	for id, rating := range ratings {
		if rating < min {
			min = rating
			minID = id
		}
	}

	req := common.TxReq{
		ID:      reqID,
		PartyID: minID,
	}
	reqID++
	return req
}

// SetRandProb is ...
func (gw *TopRandGateway) SetRandProb(prob float64) {
	gw.randP = prob
}
