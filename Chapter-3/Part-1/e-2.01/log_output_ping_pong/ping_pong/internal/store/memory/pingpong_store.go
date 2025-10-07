package memory

import (
	"fmt"
	"os"
	"sync"
)

type PingPongMem struct {
	Count int
}

type PingPongStore struct {
	mu          sync.RWMutex
	PingPongMem *PingPongMem
	path        string
}

func NewPingPongStore(path string) *PingPongStore {
	return &PingPongStore{
		path: path,
		PingPongMem: &PingPongMem{
			Count: 0,
		},
	}
}

func (ps *PingPongStore) GetCurr() (*int, error) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	return &ps.PingPongMem.Count, nil
}

func (ps *PingPongStore) Update() (*int, error) {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.PingPongMem.Count++
	if err := ps.writeToFile(ps.PingPongMem); err != nil {
		return nil, err
	}

	return &ps.PingPongMem.Count, nil
}

func (ps *PingPongStore) Reset() error {
	ps.mu.Lock()
	defer ps.mu.Unlock()
	ps.PingPongMem.Count = 0
	return nil
}

func (ps *PingPongStore) writeToFile(pm *PingPongMem) error {
	file, err := os.OpenFile(ps.path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "Ping / Pongs: %d\n", pm.Count)
	if err != nil {
		return fmt.Errorf("failed to write on log file: %w", err)
	}

	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed sync log file: %w", err)
	}

	return nil
}
