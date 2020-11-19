package market

import "github.com/shoooooman/mg-rs/common"

// TopGateway selects a party with the highest rating
type TopGateway struct{}

// GetTxReq is ...
func (sm *TopGateway) GetTxReq(ratings map[int]float64) common.TxReq {
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
