package store

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type LogStorage interface {
	Store(timestamp time.Time, value string) error
	GetLatest(n int) []LogEntry
	WriteToFile(entryLine string) error
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Value     string    `json:"value"`
}

// MemoryStore implements in-memory storage --> used for log entries
type FileMemoryStorage struct {
	entries []LogEntry // TODO: replace with generic entry?
	path    string
	mu      sync.RWMutex
}

func NewFileMemoryStorage(path string) *FileMemoryStorage {
	return &FileMemoryStorage{
		entries: make([]LogEntry, 0),
		path:    path,
	}
}

// Store used to add new entryo to memory store
func (fm *FileMemoryStorage) Store(timestamp time.Time, value string) error {
	fm.mu.Lock()
	defer fm.mu.Unlock()

	newEntry := LogEntry{
		Timestamp: timestamp,
		Value:     value,
	}

	fm.entries = append(fm.entries, newEntry)

	return nil
}

func (fm *FileMemoryStorage) WriteToFile(entryLine string) error {
	file, err := os.OpenFile(fm.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(entryLine)
	if err != nil {
		return fmt.Errorf("failed to write log on file: %w", err)
	}

	// Force immediate disk write
	if err := file.Sync(); err != nil {
		return fmt.Errorf("failed to sync log file: %w", err)
	}

	return nil
}

func (m *FileMemoryStorage) GetLatest(n int) []LogEntry {
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
