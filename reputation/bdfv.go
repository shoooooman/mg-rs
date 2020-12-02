package reputation

import (
	"log"
)

// Bdfv is a modification of Bdf, which combines feedback from others with the variable weight
// Bdfv will use the same scenario as the Bdf's one
type Bdfv struct {
	*Bdf
}

// CombineFeedback is ...
func (m *Bdfv) CombineFeedback() {
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
	untrust := bdfCalcExp(m.tparams[src])
	if untrust < bdfT || !deviate {
		m.rparams[tgt].A = m.decay*m.rparams[tgt].A + (1.0-untrust)*fb.A // variable weight
		m.rparams[tgt].B = m.decay*m.rparams[tgt].B + (1.0-untrust)*fb.B // variable weight
		m.ratings[tgt] = bdfCalcExp(m.rparams[tgt])
	}
}

// NewBdfv is ...
func NewBdfv(id int, decay float64) *Bdfv {
	bdfv := &Bdfv{
		Bdf: NewBdf(id, decay),
	}
	return bdfv
}
