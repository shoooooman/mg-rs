package market

import "github.com/shoooooman/mg-rs/common"

// Gateway is ...
type Gateway interface {
	GetTxReq(map[int]float64) common.TxReq
}
