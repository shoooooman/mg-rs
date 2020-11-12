package run

import (
	"fmt"
	"log"
	"time"

	"github.com/shoooooman/mg-rs/agent"
	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/reputation"
)

// Brs is ...
func Brs(id int) {
	a := agent.NewAgent(id)
	a.SetGateway("top")
	a.SetManager("brs")

	go func() {
		for {
			a.CombineFeedback()
		}
	}()

	const n = 20
	for i := 0; i < n; i++ {
		req := a.GetTxReq(a.GetRatings())
		party := req.PartyID
		log.Printf("req: %d with %d\n", a.ID, party)
		tx := common.Tx{ID: req.ID, PartyID: party}
		behavior := a.MonitorTx(tx)
		var result float64
		if behavior {
			result = 1.0
		} else {
			result = -1.0
		}
		a.UpdateRating(party, result)
		var bp *reputation.BrsBP
		if behavior {
			bp = &reputation.BrsBP{R: 0.0, S: 1.0}
		} else {
			bp = &reputation.BrsBP{R: 1.0, S: 0.0}
		}
		fb := &reputation.BrsFB{
			TargetID: party,
			Bp:       bp,
		}
		msg := &common.Message{SenderID: a.ID, Body: fb}
		a.BroadcastMessage(msg)
	}
	time.Sleep(time.Second * 5)
	fmt.Println(a.ID, a.GetRatings())
}
