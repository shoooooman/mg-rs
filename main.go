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
			a.TakeReputation()
		}
	}()

	const n = 5
	for i := 0; i < n; i++ {
		req := a.GetTxReq(a.GetRatings())
		party := req.Party
		log.Printf("req: %d with %d\n", a.ID, party)
		behavior := a.MonitorTx(party)
		a.UpdateRating(party, behavior)
		fb := &reputation.Feedback{
			TargetID: party,
			Bp:       &reputation.Bparams{A: 0, B: 1},
		}
		msg := common.Message{SenderID: a.ID, Body: fb}
		a.BroadcastMessage(msg)
	}
	time.Sleep(time.Second * 5)
	fmt.Println(a.ID, a.GetRatings())
}
