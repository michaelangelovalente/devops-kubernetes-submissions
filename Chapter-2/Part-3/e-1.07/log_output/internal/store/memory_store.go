package store

import (
	"sync"
	"time"
)

type LogStorage interface {
	Store(timestamp time.Time, value string) error
	GetAll() []LogEntry
	GetLatest(n int) []LogEntry
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Value     string    `json:"value"`
}

// MemoryStore implements in-memory storage --> used for log entries
type MemoryStorage struct {
	entries []LogEntry // TODO: replace with generic entry?
	mu      sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		entries: make([]LogEntry, 0),
	}
}

// Store used to add new entryo to memory store
func (m *MemoryStorage) Store(timestamp time.Time, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = append(m.entries, LogEntry{
		Timestamp: timestamp,
		Value:     value,
	})

	return nil
}

// GetAll returns a copy of all log entries (copy used to avoid external modification)
func (m *MemoryStorage) GetAll() []LogEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	res := make([]LogEntry, len(m.entries))
	copy(res, m.entries)

	return res
}

func (m *MemoryStorage) GetLatest(n int) []LogEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if n > len(m.entries) {
		n = len(m.entries)
	}

	start := len(m.entries) - n
	result := make([]LogEntry, n)
	copy(result, m.entries[start:])

	return result
}
