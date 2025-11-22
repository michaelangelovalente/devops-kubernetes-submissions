package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"log_output/internal/api"
	client "log_output/internal/client/pingpong"
	"log_output/internal/logger"
	"log_output/internal/store"
)

type Application struct {
	Logger           *logger.Logger
	wg               sync.WaitGroup
	LogMemoryHandler *api.LoggerEntryHandler
}

func NewApplication() (*Application, error) {
	logMemoryStore := store.NewMemoryStorage()
	loggerConfig := logger.LoggerConfig{
		Interval:   5 * time.Second,
		TimeFormat: time.RFC3339,
	}

	fileInfoDir := os.Getenv("FILE_INFO_TXT_PATH")
	if fileInfoDir == "" {
		log.Fatal("FILE_INFO_TXT_PATH was not set.")
	}
	pingPongURL := os.Getenv("PING_PONG_SVC_URL")
	if pingPongURL == "" {
		return nil, fmt.Errorf("PING_PONG_SVC_URL environment variable not set")
	}
	pingpongClient := client.NewClient(pingPongURL, 5*time.Second)

	logMemory := logger.NewLogger(loggerConfig, logMemoryStore)
	logMemoryHandler := api.NewLoggerEntryHandler(logMemoryStore, logMemory.GetNormalLogger(), pingpongClient, fileInfoDir)
	app := &Application{
		Logger:           logMemory,
		LogMemoryHandler: logMemoryHandler,
	}
	return app, nil
}

func (a *Application) Start(ctx context.Context) error {
	fmt.Println("Starting application services....")

	a.wg.Add(1)

	go func() {
		defer a.wg.Done()

		if err := a.Logger.StartLogger(ctx); err != nil {
			fmt.Printf("Logger start error: %v\n", err)
		}
	}()

	return nil
}

func (a *Application) Stop() error {
	fmt.Println("Stopping application services...")

	a.wg.Wait()

	fmt.Println("Application stopped succesfully")
	return nil
}
