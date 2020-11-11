package monitor

import "github.com/shoooooman/mg-rs/common"

// Monitor is ...
type Monitor interface {
	MonitorTx(common.Tx) bool
}
