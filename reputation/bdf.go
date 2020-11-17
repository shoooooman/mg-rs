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

// BdfFB is ...
type BdfFB struct {
	TargetID int
	Bp       *BdfPR
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
		m.rparams[id].B += result // success
	} else {
		m.rparams[id].A += math.Abs(result) // failure
	}
	m.ratings[id] = bdfCalcExp(m.rparams[id])
}

// BroadcastMessage is ...
func (m *Bdf) BroadcastMessage(msg *common.Message) {
	m.Broadcast(msg)
}

// CombineFeedback is ...
func (m *Bdf) CombineFeedback() {
	msg := m.GetData()
	src := msg.SenderID
	fb, ok := msg.Body.(BdfFB)
	if !ok {
		log.Fatal("body type is error")
	}
	tgt := fb.TargetID
	bp := fb.Bp
	log.Printf("reputation from %d on %d: %v\n", src, tgt, *bp)

	if tgt == m.id {
		return
	}
	deviate := bdfDeviationTest(m.rparams[tgt], bp)
	if deviate {
		m.tparams[src].A++
	} else {
		m.tparams[src].B++
	}
	if bdfCalcExp(m.tparams[src]) < bdfT || !deviate {
		m.rparams[tgt].A += bdfW * bp.A
		m.rparams[tgt].B += bdfW * bp.B
		m.ratings[tgt] = bdfCalcExp(m.rparams[tgt])
	}
}

// GetParams is for debugging
func (m *Bdf) GetParams() map[int]*BdfPR {
	return m.rparams
}

// NewBdf is ...
func NewBdf(id int) *Bdf {
	gob.Register(BdfFB{})
	bdf := &Bdf{
		id:     id,
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
