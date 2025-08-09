package app

import (
	"log_output/internal/logger"
	"log_output/internal/store"
	"time"
)

type Application struct {
	Logger *logger.Logger
	// .. Handlers
}

func NewApplication() (*Application, error) {
	logMemoryStore := store.NewMemoryStorage()

	loggerConfig := logger.LoggerConfig{
		Interval:   5 * time.Second, // replace with env var
		TimeFormat: time.RFC3339,
	}

	logger := logger.NewLogger(loggerConfig, logMemoryStore)

	//---  Handler layer ----

	//-----------------------
	app := &Application{
		Logger: logger,
	}
	return app, nil
}
