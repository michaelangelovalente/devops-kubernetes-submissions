package app

import (
	"context"
	"fmt"
	"sync"
	"time"

	"log_output/internal/api"
	"log_output/internal/logger"
	"log_output/internal/store"
)

type Application struct {
	Logger           *logger.Logger
	wg               sync.WaitGroup
	LogMemoryHandler *api.LoggerEntryHandler
}

func NewApplication() (*Application, error) {
	// --- Store layer ----
	logMemoryStore := store.NewMemoryStorage()
	loggerConfig := logger.LoggerConfig{
		Interval:   5 * time.Second,
		TimeFormat: time.RFC3339,
	}
	logMemory := logger.NewLogger(loggerConfig, logMemoryStore)

	//-----------------------

	//---  Handler layer ----
	logMemoryHandler := api.NewLoggerEntryHandler(logMemoryStore, logMemory.GetNormalLogger())
	//-----------------------
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
