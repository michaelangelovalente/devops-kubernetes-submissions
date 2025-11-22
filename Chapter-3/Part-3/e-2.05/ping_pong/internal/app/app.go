package app

import (
	"context"
	"fmt"
	"sync"

	handler "ping_pong/internal/api"
	"ping_pong/internal/store"
)

type Application struct {
	Store           *store.Store
	PingpongHandler *handler.PingPongHandler
	wg              sync.WaitGroup
}

func NewApplication() (*Application, error) {
	store := store.NewStore()
	pingpongHandler := handler.NewPingPongHandler(store)

	app := &Application{
		Store:           store,
		PingpongHandler: pingpongHandler,
	}

	return app, nil
}

func (a *Application) Start(ctx context.Context) error {
	fmt.Println("Starting application services....")

	a.wg.Add(1)

	return nil
}

func (a *Application) Stop() error {
	fmt.Println("Stopping application services...")

	a.wg.Wait()

	fmt.Println("Application stopped succesfully")
	return nil
}
