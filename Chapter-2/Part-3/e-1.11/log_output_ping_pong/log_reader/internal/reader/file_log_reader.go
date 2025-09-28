package reader

import (
	"strings"
	"sync"
)

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

	return fr.getLogsFromFile()
}

func (fr *FileLogReader) GetLatest(n int) ([]LogEntry, error) {
	fr.mu.RLock()
	defer fr.mu.Unlock()

	logs, err := fr.getLogsFromFile()
	return logs[:n], err
}

func (fr *FileLogReader) GetLastLog() (string, error) {
	logData, err := getLogsFromFile(fr.path)
	if err != nil {
		return "", err
	}

	logLines := strings.Split(logData, "\n")

	return logLines[len(logLines)-2], nil
}

func (fr *FileLogReader) getLogsFromFile() ([]LogEntry, error) {
	logContent, err := getLogsFromFile(fr.path)
	if err != nil {
		return nil, err
	}
	return append(fr.entries, parseLogContent(string(logContent))...), nil
}
