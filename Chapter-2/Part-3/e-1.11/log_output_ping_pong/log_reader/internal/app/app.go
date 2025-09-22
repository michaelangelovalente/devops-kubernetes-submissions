package app

import (
	"log"
	"os"

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
	path := os.Getenv("LOG_FILE_PATH")
	if path == "" {
		path = "/app/tmp/shared"
		log.Printf("No LOG_FILE_PATH env variable detected using default path: %s", path)
	}
	logReader := reader.NewFileLogReader(path)
	logReaderHandler := api.NewLogReaderHandler(logReader, logger)
	app := &Application{
		Logger:           logger,
		LogReaderHandler: logReaderHandler,
	}
	return app, nil
}
