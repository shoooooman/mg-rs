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
	bdfD = 0.5
	bdfT = 0.25
	bdfW = 0.1
)

// Bdf is ...
type Bdf struct {
	id      int
	rtype   RType
	decay   float64
	rparams map[int]*BdfPR
	tparams map[int]*BdfPR
	ratings map[int]float64
	network.Client
}

// BdfPR is ...
type BdfPR struct {
	A, B float64
}

// String is for debugging
func (pr *BdfPR) String() string {
	return fmt.Sprintf("(%v, %v)", pr.A, pr.B)
}

// BdfBody is ...
type BdfBody struct {
	TargetID int
	Params   *BdfPR // The params of the target owned by the sender (for deviation test)
	Fb       *BdfPR // The feedback from the sender on the target (second-hand info)
}

// InitRatings is ...
func (m *Bdf) InitRatings() {
	rparams := make(map[int]*BdfPR)
	tparams := make(map[int]*BdfPR)
	ratings := make(map[int]float64)
	ids := m.GetIDs()
	for _, id := range ids {
		if id != m.id {
			rparams[id] = &BdfPR{A: 0.0, B: 0.0}
			tparams[id] = &BdfPR{A: 0.0, B: 0.0}
			ratings[id] = 0.5
		}
	}
	m.rparams = rparams
	m.tparams = tparams
	m.ratings = ratings
}

// GetRatings is ...
func (m *Bdf) GetRatings() map[int]float64 {
	return m.ratings
}

// UpdateRating is ...
func (m *Bdf) UpdateRating(id int, result float64) {
	if result >= 0 {
		m.rparams[id].B = m.decay*m.rparams[id].B + result // success
	} else {
		m.rparams[id].A = m.decay*m.rparams[id].A + math.Abs(result) // failure
	}
	m.ratings[id] = bdfCalcExp(m.rparams[id])
}

// BroadcastMessage is ...
func (m *Bdf) BroadcastMessage(msg *common.Message) {
	switch m.rtype {
	case honest:
		m.Broadcast(msg)
	case reverse:
		body := msg.Body.(*BdfBody)
		body.Params.A, body.Params.B = body.Params.B, body.Params.A
		body.Fb.A, body.Fb.B = body.Fb.B, body.Fb.A
		msg.Body = body
		m.Broadcast(msg)
	}
}

// CombineFeedback is ...
func (m *Bdf) CombineFeedback() {
	msg := m.GetData()
	src := msg.SenderID
	body, ok := msg.Body.(BdfBody)
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
	deviate := bdfDeviationTest(m.rparams[tgt], pr)
	if deviate {
		m.tparams[src].A++
	} else {
		m.tparams[src].B++
	}
	if bdfCalcExp(m.tparams[src]) < bdfT || !deviate {
		m.rparams[tgt].A = m.decay*m.rparams[tgt].A + bdfW*fb.A
		m.rparams[tgt].B = m.decay*m.rparams[tgt].B + bdfW*fb.B
		m.ratings[tgt] = bdfCalcExp(m.rparams[tgt])
	}
}

// GetParams returns the set of parameters
func (m *Bdf) GetParams() map[int]*BdfPR {
	return m.rparams
}

// NewBdf is ...
func NewBdf(id int, decay float64) *Bdf {
	gob.Register(BdfBody{})
	conf := readConfig()
	bdf := &Bdf{
		id:     id,
		rtype:  conf.getReportType(id),
		decay:  decay,
		Client: network.NewClientImpl(id),
	}
	bdf.InitRatings()
	return bdf
}

func bdfCalcExp(p *BdfPR) float64 {
	return (p.A + 1.0) / (p.A + p.B + 2.0)
}

func bdfDeviationTest(rparams, frparams *BdfPR) bool {
	return math.Abs(bdfCalcExp(rparams)-bdfCalcExp(frparams)) >= bdfD
}
