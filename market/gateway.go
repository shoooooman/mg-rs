package market

import "github.com/shoooooman/mg-rs/common"

var reqID = 0

// Gateway is ...
type Gateway interface {
	GetTxReq(map[int]float64) common.TxReq
}
