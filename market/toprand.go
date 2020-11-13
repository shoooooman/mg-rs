package market

import (
	"math/rand"

	"github.com/shoooooman/mg-rs/common"
)

// TopRandGateway is ...
type TopRandGateway struct{}

// The probability that random policy is chosen
const randP = 1.0

// GetTxReq is ...
func (sm *TopRandGateway) GetTxReq(ratings map[int]float64) common.TxReq {
	r := rand.Intn(100)
	// choose the party randomly
	if r < randP*100 {
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

	max := -1.0
	maxID := -1
	for id, rating := range ratings {
		if rating > max {
			max = rating
			maxID = id
		}
	}

	req := common.TxReq{
		ID:      reqID,
		PartyID: maxID,
	}
	reqID++
	return req
}
