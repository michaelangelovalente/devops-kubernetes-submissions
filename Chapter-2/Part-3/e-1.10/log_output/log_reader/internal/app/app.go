package app

import (
	"log"
	"os"
	"path/filepath"

	"log_reader/internal/api"
	"log_reader/internal/reader"
)

type Application struct {
	Logger *log.Logger
	// wg               sync.WaitGroup
	LogReaderHandler *api.LogReaderHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "[LOGGER] ", log.LstdFlags)
	homePath, err := os.UserHomeDir()
	if err != nil {
		logger.Printf("ERROR: failed to get Home Directory")
	}
	path := filepath.Join(homePath, "test", "tmp", "logs.txt")
	logReader := reader.NewFileLogReader(path)
	logReaderHandler := api.NewLogReaderHandler(logReader, logger)
	app := &Application{
		Logger:           logger,
		LogReaderHandler: logReaderHandler,
	}
	return app, nil
}
