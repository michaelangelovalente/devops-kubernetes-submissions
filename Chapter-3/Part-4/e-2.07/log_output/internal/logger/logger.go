package logger

import (
	"context"
	"fmt"
	"log"
	"log_output/internal/store"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ValueGenerator is a function type utility,
// it is used to generate values (used to inject different value generation strategies --- behavioral parametr)
type ValueGenerator func() string

// LoggerConfig holds logger configuration
type LoggerConfig struct {
	Interval   time.Duration
	TimeFormat string
}

// Logger core struct, contains logging logic
type Logger struct {
	loggerConfig LoggerConfig
	logStorage   store.LogStorage
	currentValue string
	rwMutex      sync.RWMutex // Protects "currentValue" (thread safety)
	normalLogger *log.Logger  // Standard Go logger for normal logging (not stored)
}

func NewLogger(loggerConfig LoggerConfig, logStorage store.LogStorage) *Logger {
	return &Logger{
		loggerConfig: loggerConfig,
		logStorage:   logStorage,
		normalLogger: log.New(os.Stdout, "[LOGGER] ", log.LstdFlags),
		// Note --> rwMutex doesn't need initialization.
	}
}

func (l *Logger) StartLogger(ctx context.Context) error {
	// Generate initial value once
	l.rwMutex.Lock()
	l.currentValue = generateUUID()
	l.rwMutex.Unlock()

	// Create ticker --> periodic logging
	ticker := time.NewTicker(l.loggerConfig.Interval)
	defer ticker.Stop()

	// immediate log on Start
	l.logCurrent()

	// log on time intervals till conext cancellation
	for {
		select {
		case <-ctx.Done():
			l.normalLogger.Println("Logger stopped...")
			return nil
		case <-ticker.C:
			l.logCurrent()
		}
	}
}

func (l *Logger) logCurrent() {
	l.rwMutex.RLock()
	value := l.currentValue
	l.rwMutex.RUnlock()

	timestamp := time.Now()

	if err := l.logStorage.Store(timestamp, value); err != nil {
		l.normalLogger.Printf("Error storing log: %v", err)
		return
	}

	// Output generated log value to console (stored in memory)
	fmt.Printf("%s: %s\n", timestamp.Format(l.loggerConfig.TimeFormat), value)
}

func generateUUID() string {
	return uuid.New().String()
}

// GetNormalLogger returns the normal logger instance for external use
func (l *Logger) GetNormalLogger() *log.Logger {
	return l.normalLogger
}
