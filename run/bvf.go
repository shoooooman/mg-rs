package run

import (
	"log"

	"github.com/shoooooman/mg-rs/agent"
	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/reputation"
)

// Bvf is ...
func Bvf(a *agent.Agent, n int) {
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

		var params *reputation.BvfPR
		switch m := a.Manager.(type) {
		case *reputation.Bvf:
			params = m.GetParams()[party]
		case *reputation.Bvfv:
			params = m.GetParams()[party]
		default:
			log.Fatal("GetParams can be called only with bvf or bvfv")
		}

		var fb *reputation.BvfPR
		if behavior {
			fb = &reputation.BvfPR{A: 0.0, B: 1.0}
		} else {
			fb = &reputation.BvfPR{A: 1.0, B: 0.0}
		}

		body := &reputation.BvfBody{
			TargetID: party,
			Params:   params,
			Fb:       fb,
		}
		msg := &common.Message{SenderID: a.ID, Body: body}
		a.BroadcastMessage(msg)

		peers := a.GetPeers()
		for i := 0; i < len(peers); i++ {
			a.CombineFeedback()
		}
		rlog.Printf("[ratings] %d %d %v\n", a.ID, i, a.GetRatings())
	}
	// for debugging
	switch m := a.Manager.(type) {
	case *reputation.Bvf:
		log.Printf("%d: %v\n", a.ID, m.GetParams())
		log.Printf("%d's tparams: %v\n", a.ID, m.GetTParams())
		log.Printf("%d: %v\n", a.ID, m.GetFeedbacks())
	case *reputation.Bvfv:
		log.Printf("%d: %v\n", a.ID, m.GetParams())
		log.Printf("%d's tparams: %v\n", a.ID, m.GetTParams())
		log.Printf("%d: %v\n", a.ID, m.GetFeedbacks())
	}
	log.Printf("%d: (success, failure)=(%d, %d)\n", a.ID, success, failure)
	rlog.Printf("[result] %d %d %d\n", a.ID, success, failure)
}
