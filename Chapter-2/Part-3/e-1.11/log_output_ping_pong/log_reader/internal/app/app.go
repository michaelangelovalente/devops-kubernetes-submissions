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
	pingpongPath := os.Getenv("PING_PONG_FILE_PATH")
	if path == "" {
		path = "/app/tmp/shared/logs.txt"
		log.Printf("No LOG_FILE_PATH env variable detected using default path: %s", path)
	}
	if path == "" {
		path = "/app/tmp/shared/pingpong.txt"
		log.Printf("No PING_PONG_FILE_PATH env variable detected using default path: %s", path)
	}

	logReader := reader.NewFileLogReader(path)
	pingpongReader := reader.NewPingPongReader(pingpongPath)
	logReaderHandler := api.NewLogReaderHandler(logReader, *pingpongReader, logger)
	app := &Application{
		Logger:           logger,
		LogReaderHandler: logReaderHandler,
	}
	return app, nil
}
