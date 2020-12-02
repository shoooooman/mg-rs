package reputation

import (
	"log"
)

// Bvfv is a modification of Bvf, which combines feedback from others with the variable weight
// Bvfv will use the same scenario as the Bvf's one
type Bvfv struct {
	*Bvf
}

// CombineFeedback is ...
func (m *Bvfv) CombineFeedback() {
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

// NewBvfv is ...
func NewBvfv(id int, decay float64) *Bvfv {
	bvfv := &Bvfv{
		Bvf: NewBvf(id, decay),
	}
	return bvfv
}
