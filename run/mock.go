package run

import (
	"fmt"
	"log"
	"time"

	"github.com/shoooooman/mg-rs/agent"
	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/reputation"
)

// Mock is ...
func Mock(a *agent.Agent, n int) {
	go func() {
		for {
			a.CombineFeedback()
		}
	}()

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
		fb := &reputation.Feedback{
			TargetID: party,
			Bp:       &reputation.Bparams{A: 0, B: 1},
		}
		msg := &common.Message{SenderID: a.ID, Body: fb}
		a.BroadcastMessage(msg)
	}
	time.Sleep(time.Second * 5)
	fmt.Println(a.ID, a.GetRatings())
}
