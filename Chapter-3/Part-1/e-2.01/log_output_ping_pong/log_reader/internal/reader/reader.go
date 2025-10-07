package reader

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type Reader interface {
	GetAll() ([]LogEntry, error)
	GetLatest(n int) ([]LogEntry, error)
	GetLastLog() (string, error)
}

type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Value     string    `json:"value"`
}

func parseLogContent(content string) []LogEntry {
	lines := strings.Split(content, "\n")

	var newEntries []LogEntry

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

		newEntries = append(newEntries, LogEntry{
			Timestamp: timestamp,
			Value:     value,
		})

	}

	return newEntries
}

func getLogsFromFile(path string) (string, error) {

	logContent, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("ERROR: file does not exist")
			return "", fmt.Errorf("ERROR: file does not exists: %w", err)
		}
		return "", nil
	}

	return string(logContent), nil
}
