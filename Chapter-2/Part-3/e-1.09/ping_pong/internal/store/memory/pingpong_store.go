package memory

import (
	"sync"
)

type PinpongModel struct {
	Count int
}

type PingPongStore struct {
	mu            sync.RWMutex
	PingPongModel *PinpongModel
}

func NewPingPongStore() *PingPongStore {
	return &PingPongStore{
		PingPongModel: &PinpongModel{
			Count: 10,
		},
	}
}

func (ps *PingPongStore) GetCurr() (int, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	return ps.PingPongModel.Count, nil
}

func (ps *PingPongStore) Update() (int, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.PingPongModel.Count++
	return ps.PingPongModel.Count, nil
}

func (ps *PingPongStore) Reset() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.PingPongModel.Count = 0
	return nil
}
