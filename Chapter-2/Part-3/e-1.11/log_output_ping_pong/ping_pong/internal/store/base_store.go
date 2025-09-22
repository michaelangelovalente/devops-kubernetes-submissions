package store

import "ping_pong/internal/store/memory"

type Store struct {
	PingPongStore *memory.PingPongStore
	// other stores...
}

type PingPongStore interface {
	Update() (*int, error)
	GetCurr() (*int, error)
	Reset() error
}

func NewStore(pingpongStore *memory.PingPongStore) *Store {
	return &Store{
		PingPongStore: pingpongStore,
	}
}
