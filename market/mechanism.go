package market

import "github.com/shoooooman/mg-rs/common"

// SortGateway is ...
type SortGateway struct{}

var reqID = 0

// GetTxReq is ...
func (sm *SortGateway) GetTxReq(ratings map[int]float64) common.TxReq {
	max := -1.0
	maxID := -1
	for id, rating := range ratings {
		if rating > max {
			max = rating
			maxID = id
		}
	}

	req := common.TxReq{
		ID:    reqID,
		Party: maxID,
	}
	reqID++
	return req
}