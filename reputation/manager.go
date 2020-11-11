package reputation

import (
	"github.com/shoooooman/mg-rs/common"
	"github.com/shoooooman/mg-rs/network"
)

// Algorithm is ...
type Algorithm interface {
	GetRatings() map[int]float64
	UpdateRating(int, bool)
	BroadcastMessage(*common.Message)
	TakeReputation()
}

// Manager is ...
type Manager interface {
	Algorithm
	network.Client
}
