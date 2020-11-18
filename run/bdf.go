package run

import (
	"fmt"
	"log"

	"github.com/shoooooman/mg-rs/agent"
	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/reputation"
)

// Bdf is ...
func Bdf(a *agent.Agent, n int) {
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

		var result float64
		if behavior {
			result = 1.0
			success++
		} else {
			result = -1.0
			failure++
		}
		log.Printf("result (%d with %d): %v\n", a.ID, party, result)
		a.UpdateRating(party, result)

		var bp *reputation.BdfPR
		if behavior {
			bp = &reputation.BdfPR{A: 0.0, B: 1.0}
		} else {
			bp = &reputation.BdfPR{A: 1.0, B: 0.0}
		}
		fb := &reputation.BdfFB{
			TargetID: party,
			Bp:       bp,
		}
		msg := &common.Message{SenderID: a.ID, Body: fb}
		a.BroadcastMessage(msg)

		peers := a.GetPeers()
		for i := 0; i < len(peers); i++ {
			a.CombineFeedback()
		}
		rlog.Printf("[ratings] %d %d %v\n", a.ID, i, a.GetRatings())
	}
	log.Printf("%d: %v\n", a.ID, a.Manager.(*reputation.Bdf).GetParams())
	log.Printf("%d: (success, failure)=(%d, %d)\n", a.ID, success, failure)
	rlog.Printf("[result] %d %d %d\n", a.ID, success, failure)
	fmt.Println(a.ID, a.GetRatings())
}