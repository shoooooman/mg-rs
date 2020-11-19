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
	trust := bdfCalcExp(m.tparams[src])
	if trust < bdfT || !deviate {
		m.rparams[tgt].A += trust * bp.A // variable weight
		m.rparams[tgt].B += trust * bp.B // variable weight
		m.ratings[tgt] = bdfCalcExp(m.rparams[tgt])
	}
}

// NewBdfv is ...
func NewBdfv(id int) *Bdfv {
	gob.Register(BdfFB{})
	bdfv := &Bdfv{
		Bdf: NewBdf(id),
	}
	bdfv.InitRatings()
	return bdfv
}
