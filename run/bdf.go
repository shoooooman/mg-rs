package run

import (
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

		var params *reputation.BdfPR
		switch m := a.Manager.(type) {
		case *reputation.Bdf:
			params = m.GetParams()[party]
		case *reputation.Bdfv:
			params = m.GetParams()[party]
		default:
			log.Fatal("GetParams can be called only with bdf or bdfv")
		}

		var fb *reputation.BdfPR
		if behavior {
			fb = &reputation.BdfPR{A: 0.0, B: 1.0}
		} else {
			fb = &reputation.BdfPR{A: 1.0, B: 0.0}
		}

		body := &reputation.BdfBody{
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
	// to use GetParams() for debugging
	switch m := a.Manager.(type) {
	case *reputation.Bdf:
		log.Printf("%d: %v\n", a.ID, m.GetParams())
	case *reputation.Bdfv:
		log.Printf("%d: %v\n", a.ID, m.GetParams())
	}
	log.Printf("%d: (success, failure)=(%d, %d)\n", a.ID, success, failure)
	rlog.Printf("[result] %d %d %d\n", a.ID, success, failure)
}
