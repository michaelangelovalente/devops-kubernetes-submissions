package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"log_output/internal/logger"
	"log_output/internal/store"
)

type Application struct {
	Logger *logger.Logger
	wg     sync.WaitGroup
}

func NewApplication() (*Application, error) {
	path := os.Getenv("LOG_FILE_PATH")
	if path == "" {
		path = "/app/tmp/shared"
		log.Printf("No LOG_FILE_PATH env variable detected using default path: %s", path)
	}

	logMemoryStore := store.NewFileMemoryStorage(path)
	loggerConfig := logger.LoggerConfig{
		Interval:   5 * time.Second,
		TimeFormat: time.RFC3339,
	}
	logMemory := logger.NewLogger(loggerConfig, logMemoryStore)

	app := &Application{
		Logger: logMemory,
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
