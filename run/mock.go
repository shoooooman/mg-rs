package run

import (
	"fmt"
	"log"

	"github.com/shoooooman/mg-rs/agent"
	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/reputation"
)

// Mock is ...
func Mock(a *agent.Agent, n int) {
	var (
		success = 0
		failure = 0
	)
	for i := 0; i < n; i++ {
		req := a.GetTxReq(a.GetRatings())
		party := req.PartyID
		log.Printf("req: %d with %d\n", a.ID, party)

		tx := common.Tx{ID: req.ID, PartyID: party}
		behavior := a.MonitorTx(tx)

		var result float64
		if behavior {
			result = 1.0
			success++
		} else {
			result = -1.0
			failure++
		}
		a.UpdateRating(party, result)

		fb := &reputation.Feedback{
			TargetID: party,
			Bp:       &reputation.Bparams{A: 0, B: 1},
		}
		msg := &common.Message{SenderID: a.ID, Body: fb}
		a.BroadcastMessage(msg)

		peers := a.GetPeers()
		for j := 0; j < len(peers); j++ {
			a.CombineFeedback()
		}
	}
	log.Printf("%d: (success, failure)=(%d, %d)\n", a.ID, success, failure)
	fmt.Println(a.ID, a.GetRatings())
}
