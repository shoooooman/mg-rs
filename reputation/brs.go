package reputation

import (
	"encoding/gob"
	"log"
	"math"

	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/network"
)

// Brs represents the Beta Reputation System manager
type Brs struct {
	id      int
	params  map[int]*BrsBP
	ratings map[int]float64
	*network.ClientImpl
}

// BrsBP is ...
type BrsBP struct {
	R, S float64
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
	peerIDs := m.GetPeers()
	for _, p := range peerIDs {
		params[p] = &BrsBP{R: 0.0, S: 0.0}
		ratings[p] = 0.5
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
		m.params[id].R += result
	} else {
		m.params[id].S += math.Abs(result)
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
	log.Printf("reputation from %d on %d: %v\n", src, tgt, *bp)
	if tgt != m.id {
		den := (m.params[src].S+2.0)*(bp.R+bp.S+2.0) + 2.0*m.params[src].R
		m.params[tgt].R = 2.0 * m.params[src].R * bp.R / den
		m.params[tgt].S = 2.0 * m.params[src].R * bp.S / den
		m.ratings[tgt] = brsCalcExp(m.params[tgt])
	}
}

// NewBrs is ...
func NewBrs(id int) *Brs {
	gob.Register(BrsFB{})
	mock := &Brs{
		id:         id,
		ClientImpl: network.NewClientImpl(id),
	}
	mock.InitRatings()
	return mock
}

func brsCalcExp(p *BrsBP) float64 {
	return (p.R + 1.0) / (p.R + p.S + 2.0)
}