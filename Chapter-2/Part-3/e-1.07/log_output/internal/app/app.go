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
	Logger        *logger.Logger
	wg            sync.WaitGroup
	LoggerHandler *api.LoggerEntryHandler
	// .. Handlers
}

func NewApplication() (*Application, error) {
	logMemoryStore := store.NewMemoryStorage()

	loggerConfig := logger.LoggerConfig{
		Interval:   5 * time.Second,
		TimeFormat: time.RFC3339,
	}

	log := logger.NewLogger(loggerConfig, logMemoryStore)

	//---  Handler layer ----
	loggerHandler := api.NewLoggerEntryHandler(logMemoryStore)

	//-----------------------
	app := &Application{
		Logger:        log,
		LoggerHandler: loggerHandler,
	}
	return app, nil
}

func (a *Application) Start(ctx context.Context) error {
	fmt.Println("Starting applicataion services....")

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
