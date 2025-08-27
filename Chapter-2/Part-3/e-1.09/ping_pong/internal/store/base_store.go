package store

import (
	"ping_pong/internal/store/memory"
)

type Store struct {
	PingPongStore PingPongStore
	// other stores...
}

// Store interfaces (for now only PingPongStore)
type PingPongStore interface {
	Update() (int, error)
	GetCurr() (int, error)
	Reset() error
}

func NewStore() *Store {
	return &Store{
		PingPongStore: memory.NewPingPongStore(),
	}
}
