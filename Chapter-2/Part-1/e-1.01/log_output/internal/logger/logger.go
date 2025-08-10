package logger

//TODO: --> Change struct, func, method names: remove Log or Logger
import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// ValueGenerator is a function type utility,
// it is used to generate values (used to inject different value generation strategies --- behavioral parametr)
type ValueGenerator func() string

// LogStorage defines interface for storing log entries
type LogStorage interface {
	Store(timestamp time.Time, value string) error
	GetAll() []LogEntry
}

// LogEntry represents single log entry
type LogEntry struct {
	Timestamp time.Time
	Value     string
}

// LoggerConfig holds logger configuration
type Config struct {
	Interval      time.Duration
	GenerateValue ValueGenerator
	TimeFormat    string
}

// Logger core struct, contains logging logic
type Logger struct {
	loggerConfig Config
	logStorage   LogStorage
	currentValue string
	rwMutex      sync.RWMutex // Protects "currentValue" (thread safety)
}

func NewLogger(loggerConfig Config, logStorage LogStorage) *Logger {
	return &Logger{
		loggerConfig: loggerConfig,
		logStorage:   logStorage,
		// Note --> rwMutex doesn't need initialization.
	}
}

func (l *Logger) StartLogger(ctx context.Context) error {
	// Generate initial value once
	l.rwMutex.Lock()
	l.currentValue = l.loggerConfig.GenerateValue()
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
			fmt.Println("Logger stopeed...")
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
		fmt.Printf("Error storing log: %v\n", err)
		return
	}

	// Output to console
	fmt.Printf("%s: %s\n", timestamp.Format(l.loggerConfig.TimeFormat), value)

}

func GenerateUUID() string {
	return uuid.New().String()
}
