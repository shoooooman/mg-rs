package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/shoooooman/mg-rs/agent"
	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/reputation"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("input length is too small")
	}

	id, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal("id error")
	}

	a := agent.NewAgent(id)

	go func() {
		for {
			a.CombineFeedback()
		}
	}()

	const n = 5
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
