package market

import (
	"math/rand"

	"github.com/shoooooman/mg-rs/common"
)

// RandomGateway is ...
type RandomGateway struct{}

// GetTxReq is ...
func (rg *RandomGateway) GetTxReq(ratings map[int]float64) common.TxReq {
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
