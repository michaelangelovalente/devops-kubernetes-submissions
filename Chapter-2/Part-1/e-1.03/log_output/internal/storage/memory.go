package storage

import (
	"log_output/internal/logger"
	"sync"
	"time"
)

// Entry is astored logged entry
// type LogEntry struct {
// 	Timestamp time.Time
// 	Value     string
// }

// MemoryStore implements in-memory storage --> used for log entries
type MemoryStorage struct {
	entries []logger.LogEntry
	mu      sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		entries: make([]logger.LogEntry, 0),
	}
}

// Store used to add new entryo to memory store
func (m *MemoryStorage) Store(timestamp time.Time, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = append(m.entries, logger.LogEntry{
		Timestamp: timestamp,
		Value:     value,
	})

	return nil
}

// GetAll returns a copy of all log entries (copy used to avoid external modification)
func (m *MemoryStorage) GetAll() []logger.LogEntry {
	m.mu.Lock()
	defer m.mu.RUnlock()

	res := make([]logger.LogEntry, len(m.entries))
	copy(res, m.entries)

	return res
}

func (m *MemoryStorage) GetLatest(n int) []logger.LogEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if n > len(m.entries) {
		n = len(m.entries)
	}

	start := len(m.entries) - n
	result := make([]logger.LogEntry, n)
	copy(result, m.entries[start:])

	return result
}
