package run

import (
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
	a.InitRatings()
	for i := 0; i < n; i++ {
		req := a.GetTxReq(a.GetRatings())
		party := req.PartyID

		tx := common.Tx{ID: req.ID, Time: i, PartyID: party}
		behavior := a.MonitorTx(tx)
		log.Printf("%d with %d: %v\n", a.ID, party, behavior)

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
		log.Printf("%d: %dth tx finished\n", a.ID, i)
	}
	log.Printf("%d: (success, failure)=(%d, %d)\n", a.ID, success, failure)
	rlog.Printf("[result] %d %d %d\n", a.ID, success, failure)
}
