package reader

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type LogReader interface {
	GetAll() ([]LogEntry, error)
	// GetLatest(n int) ([]LogEntry, error)
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Value     string    `json:"value"`
}

type FileLogReader struct {
	entries []LogEntry
	path    string
	mu      sync.RWMutex
}

func NewFileLogReader(path string) *FileLogReader {
	return &FileLogReader{
		entries: make([]LogEntry, 0),
		path:    path,
	}
}

func (fr *FileLogReader) GetAll() ([]LogEntry, error) {
	fr.mu.RLock()
	defer fr.mu.RUnlock()

	logContent, err := os.ReadFile(fr.path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("ERROR: file does not exist")
			return nil, fmt.Errorf("ERROR: file does not exists: %w", err)
		}
		return []LogEntry{}, nil
	}

	return fr.parseContent(string(logContent)), nil
}

func (fr *FileLogReader) parseContent(content string) []LogEntry {
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		separatorIdx := strings.Index(line, ": ")
		if separatorIdx == -1 {
			continue
		}

		timestampStr := line[:separatorIdx]
		value := line[separatorIdx+2:]

		timestamp, err := time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			continue // malformed lines
		}

		fr.entries = append(fr.entries, LogEntry{
			Timestamp: timestamp,
			Value:     value,
		})

	}

	return fr.entries
}
