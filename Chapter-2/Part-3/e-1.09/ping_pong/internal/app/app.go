// package app equiv to server in other ms
package app

import (
	"context"
	"fmt"
	"sync"

	handler "ping_pong/internal/api"
	"ping_pong/internal/store"
)

type Application struct {
	Store *store.Store
	// services *Services
	// handlers *Handlers
	PingpongHandler *handler.PingPongHandler // temporary solution
	// Logger *logger.Logger
	wg sync.WaitGroup
	// Lifecycle context management
	baseCtx context.Context
}

func NewApplication(ctx context.Context) (*Application, error) {
	// ------------------ Store Layer -------------
	store := store.NewStore()
	// ---------------------------------------------

	// ---------------------------------------------
	// Services layer
	// services := &Services{
	//  .....
	// }
	// -----------------------------------------------

	// ------------------ Handlers Layer -------------
	// handlers := &Handlers{
	//  .....
	// }
	pingpongHandler := handler.NewPingPongHandler(store)
	// -----------------------------------------------

	app := &Application{
		Store: store,
		// 	services: services,
		// 	handlers: handlers,
		PingpongHandler: pingpongHandler,
		baseCtx:         ctx,
	}

	return app, nil
}

func (a *Application) Start(ctx context.Context) error {
	fmt.Println("Starting applicataion services....")

	a.wg.Add(1)

	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	fmt.Println("Stopping application services...")

	a.wg.Wait()

	fmt.Println("Application stopped succesfully")
	return nil
}
