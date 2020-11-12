package reputation

import (
	"encoding/gob"
	"log"

	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/network"
)

// MockManager is ...
type MockManager struct {
	id      int
	ratings map[int]float64
	*network.ClientImpl
}

// Bparams is ...
type Bparams struct {
	A, B int
}

// Feedback is ...
type Feedback struct {
	TargetID int
	Bp       *Bparams
}

// GetRatings is ...
func (m *MockManager) GetRatings() map[int]float64 {
	return m.ratings
}

// UpdateRating is ...
func (m *MockManager) UpdateRating(id int, result float64) {
	m.ratings[id] += result
}

// BroadcastMessage is ...
func (m *MockManager) BroadcastMessage(msg *common.Message) {
	m.Broadcast(msg)
}

// CombineFeedback is ...
func (m *MockManager) CombineFeedback() {
	msg := m.GetData()
	s := msg.SenderID
	fb, ok := msg.Body.(Feedback)
	if !ok {
		log.Fatal("body type is error")
	}
	t := fb.TargetID
	bp := fb.Bp
	log.Printf("reputation from %d on %d: %v\n", s, t, *bp)
	if t != m.id {
		m.ratings[t] = m.ratings[t] - float64(bp.A) + float64(bp.B)
	}
}

// FIXME: mock
func initRatings(id int) map[int]float64 {
	ratings := make(map[int]float64)
	if id == 0 {
		ratings[1] = 0
	} else {
		ratings[0] = 0
	}
	return ratings
}

// NewMockManager is ...
func NewMockManager(id int) *MockManager {
	gob.Register(Feedback{})
	mock := &MockManager{
		id:         id,
		ratings:    initRatings(id),
		ClientImpl: network.NewClientImpl(id),
	}
	return mock
}
