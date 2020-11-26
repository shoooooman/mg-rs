package reputation

import (
	"encoding/gob"
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
	trust := bdfCalcExp(m.tparams[src])
	if trust < bdfT || !deviate {
		m.rparams[tgt].A = m.decay*m.rparams[tgt].A + trust*fb.A // variable weight
		m.rparams[tgt].B = m.decay*m.rparams[tgt].B + trust*fb.B // variable weight
		m.ratings[tgt] = bdfCalcExp(m.rparams[tgt])
	}
}

// NewBdfv is ...
func NewBdfv(id int, decay float64) *Bdfv {
	gob.Register(BdfBody{})
	bdfv := &Bdfv{
		Bdf: NewBdf(id, decay),
	}
	bdfv.InitRatings()
	return bdfv
}
