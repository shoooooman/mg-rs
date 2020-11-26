package reputation

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"

	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/network"
)

// Brs represents the Beta Reputation System manager
type Brs struct {
	id      int
	decay   float64
	params  map[int]*BrsBP
	ratings map[int]float64
	network.Client
}

// BrsBP is ...
type BrsBP struct {
	R, S float64
}

// String is for debugging
func (bp *BrsBP) String() string {
	return fmt.Sprintf("(%v, %v)", bp.R, bp.S)
}

// BrsFB is ...
type BrsFB struct {
	TargetID int
	Bp       *BrsBP
}

// InitRatings is ...
func (m *Brs) InitRatings() {
	params := make(map[int]*BrsBP)
	ratings := make(map[int]float64)
	ids := m.GetIDs()
	for _, id := range ids {
		if id != m.id {
			params[id] = &BrsBP{R: 0.0, S: 0.0}
			ratings[id] = 0.5
		}
	}
	m.params = params
	m.ratings = ratings
}

// GetRatings is ...
func (m *Brs) GetRatings() map[int]float64 {
	return m.ratings
}

// UpdateRating is ...
func (m *Brs) UpdateRating(id int, result float64) {
	if result >= 0 {
		m.params[id].S = m.decay*m.params[id].S + result // success
	} else {
		m.params[id].R = m.decay*m.params[id].R + math.Abs(result) // failure
	}
	m.ratings[id] = brsCalcExp(m.params[id])
}

// BroadcastMessage is ...
func (m *Brs) BroadcastMessage(msg *common.Message) {
	m.Broadcast(msg)
}

// CombineFeedback is ...
func (m *Brs) CombineFeedback() {
	msg := m.GetData()
	src := msg.SenderID
	fb, ok := msg.Body.(BrsFB)
	if !ok {
		log.Fatal("body type is error")
	}
	tgt := fb.TargetID
	bp := fb.Bp
	log.Printf("%d: reputation from %d on %d: %v\n", m.id, src, tgt, *bp)
	if tgt != m.id {
		den := (m.params[src].S+2.0)*(bp.R+bp.S+2.0) + 2.0*m.params[src].R
		m.params[tgt].R = m.decay*m.params[tgt].R + 2.0*m.params[src].R*bp.R/den
		m.params[tgt].S = m.decay*m.params[tgt].S + 2.0*m.params[src].R*bp.S/den
		m.ratings[tgt] = brsCalcExp(m.params[tgt])
	}
}

// GetParams is for debugging
func (m *Brs) GetParams() map[int]*BrsBP {
	return m.params
}

// NewBrs is ...
func NewBrs(id int, decay float64) *Brs {
	gob.Register(BrsFB{})
	brs := &Brs{
		id:     id,
		decay:  decay,
		Client: network.NewClientImpl(id),
	}
	brs.InitRatings()
	return brs
}

func brsCalcExp(p *BrsBP) float64 {
	return (p.R + 1.0) / (p.R + p.S + 2.0)
}
