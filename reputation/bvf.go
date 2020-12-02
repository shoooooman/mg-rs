package reputation

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"

	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/network"
)

const (
	bvfT = 0.25
)

// Bvf is Beta Verification Feedback
type Bvf struct {
	id        int
	rtype     RType
	decay     float64
	rparams   map[int]*BvfPR
	tparams   map[int]*BvfPR
	ratings   map[int]float64
	feedbacks map[int]map[int]float64
	network.Client
}

// BvfPR is ...
type BvfPR struct {
	A, B float64
}

// String is for debugging
func (pr *BvfPR) String() string {
	return fmt.Sprintf("(%v, %v)", pr.A, pr.B)
}

// BvfBody is ...
type BvfBody struct {
	TargetID int
	Params   *BvfPR
	Fb       *BvfPR
}

// InitRatings is ...
func (m *Bvf) InitRatings() {
	rparams := make(map[int]*BvfPR)
	tparams := make(map[int]*BvfPR)
	ratings := make(map[int]float64)
	feedbacks := make(map[int]map[int]float64)
	ids := m.GetIDs()
	for _, id := range ids {
		if id != m.id {
			rparams[id] = &BvfPR{A: 0.0, B: 0.0}
			tparams[id] = &BvfPR{A: 0.0, B: 0.0}
			ratings[id] = 0.5
			feedbacks[id] = make(map[int]float64)
		}
	}
	m.rparams = rparams
	m.tparams = tparams
	m.ratings = ratings
	m.feedbacks = feedbacks
}

// GetRatings is ...
func (m *Bvf) GetRatings() map[int]float64 {
	return m.ratings
}

// UpdateRating is ...
func (m *Bvf) UpdateRating(id int, result float64) {
	if result >= 0 {
		m.rparams[id].B = m.decay*m.rparams[id].B + result // success
	} else {
		m.rparams[id].A = m.decay*m.rparams[id].A + math.Abs(result) // failure
	}
	m.ratings[id] = bvfCalcExp(m.rparams[id])

	fbs := m.feedbacks[id]
	for src, fb := range fbs {
		if result >= 0 { // success
			if fb >= 0.5 {
				m.tparams[src].B++ // correct
			} else {
				m.tparams[src].A++ // incorrect
			}
		} else { // failure
			if fb < 0.5 {
				m.tparams[src].B++ // correct
			} else {
				m.tparams[src].A++ // incorrect
			}
		}
	}
}

// BroadcastMessage is ...
func (m *Bvf) BroadcastMessage(msg *common.Message) {
	switch m.rtype {
	case honest:
		m.Broadcast(msg)
	case reverse:
		body := msg.Body.(*BvfBody)
		body.Params.A, body.Params.B = body.Params.B, body.Params.A
		body.Fb.A, body.Fb.B = body.Fb.B, body.Fb.A
		msg.Body = body
		m.Broadcast(msg)
	}
}

// CombineFeedback is ...
func (m *Bvf) CombineFeedback() {
	msg := m.GetData()
	src := msg.SenderID
	body, ok := msg.Body.(BvfBody)
	if !ok {
		log.Fatal("body type is error")
	}
	tgt := body.TargetID
	pr := body.Params
	fb := body.Fb
	log.Printf("reputation from %d on %d: %v\n", src, tgt, *fb)

	if tgt == m.id {
		return
	}

	m.feedbacks[tgt][src] = bvfCalcExp(pr)

	untrust := bvfCalcExp(m.tparams[src])
	if untrust < bvfT {
		m.rparams[tgt].A = m.decay*m.rparams[tgt].A + (1.0-untrust)*fb.A
		m.rparams[tgt].B = m.decay*m.rparams[tgt].B + (1.0-untrust)*fb.B
		m.ratings[tgt] = bvfCalcExp(m.rparams[tgt])
	}
}

// GetParams returns the set of parameters
func (m *Bvf) GetParams() map[int]*BvfPR {
	return m.rparams
}

// GetFeedbacks is for debugging
func (m *Bvf) GetFeedbacks() map[int]map[int]float64 {
	return m.feedbacks
}

// NewBvf is ...
func NewBvf(id int, decay float64) *Bvf {
	gob.Register(BvfBody{})
	conf := readConfig()
	bvf := &Bvf{
		id:     id,
		rtype:  conf.getReportType(id),
		decay:  decay,
		Client: network.NewClientImpl(id),
	}
	bvf.InitRatings()
	return bvf
}

func bvfCalcExp(p *BvfPR) float64 {
	return (p.A + 1.0) / (p.A + p.B + 2.0)
}
