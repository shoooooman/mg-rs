package reputation

import (
	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/network"
)

// Algorithm is ...
type Algorithm interface {
	GetRatings() map[int]float64
	UpdateRating(int, float64)
	BroadcastMessage(*common.Message)
	CombineFeedback()
}

// Manager is ...
type Manager interface {
	Algorithm
	network.Client
}
