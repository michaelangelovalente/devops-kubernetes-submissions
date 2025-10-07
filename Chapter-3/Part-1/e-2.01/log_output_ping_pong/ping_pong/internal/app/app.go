// package app equiv to server in other ms
package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	handler "ping_pong/internal/api"
	"ping_pong/internal/store"
	"ping_pong/internal/store/memory"
)

type Application struct {
	Store           *store.Store
	PingpongHandler *handler.PingPongHandler
	wg              sync.WaitGroup
	// Lifecycle context management
	baseCtx context.Context
}

func NewApplication(ctx context.Context) (*Application, error) {
	path := os.Getenv("PING_PONG_FILE_PATH")
	if path == "" {
		path := "/app/tmp/shared"
		log.Printf("No LOG_FILE_PATH env variable detected using default path: %s", path)
	}
	pingpongStore := memory.NewPingPongStore(path)
	store := store.NewStore(pingpongStore)

	pingpongHandler := handler.NewPingPongHandler(store)

	app := &Application{
		Store:           store,
		PingpongHandler: pingpongHandler,
		baseCtx:         ctx,
	}

	return app, nil
}

func (a *Application) Start(ctx context.Context) error {
	fmt.Println("Starting application services....")

	a.wg.Add(1)

	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	fmt.Println("Stopping application services...")

	a.wg.Wait()

	fmt.Println("Application stopped succesfully")
	return nil
}
