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
		// 	workers:  []Worker{},
		baseCtx: ctx,
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

// func (a *Application) Start(ctx context.Context) error {
// 	// Start all background workers with their own context
// 	for _, worker := range a.workers {
// 		w := worker // Capture loop variable
//
// 		a.wg.Add(1)
// 		go func() {
// 			defer a.wg.Done()
//
// 			// Each worker gets a child context
// 			if err := w.Start(ctx); err != nil {
// 				slog.Error("worker failed",
// 					"worker", w.Name(),
// 					"error", err)
// 			}
// 		}()
// 	}
//
// 	return nil
// }

// func (a *Application) Stop(ctx context.Context) error {
// 	// Signal all workers to stop
// 	var stopErrors []error
//
// 	for _, worker := range a.workers {
// 		if err := worker.Stop(ctx); err != nil {
// 			stopErrors = append(stopErrors,
// 				fmt.Errorf("failed to stop %s: %w", worker.Name(), err))
// 		}
// 	}
//
// 	// Wait for all workers to finish with timeout
// 	done := make(chan struct{})
// 	go func() {
// 		a.wg.Wait()
// 		close(done)
// 	}()
//
// 	select {
// 	case <-done:
// 		slog.Info("all workers stopped successfully")
// 	case <-ctx.Done():
// 		return fmt.Errorf("shutdown timeout while waiting for workers")
// 	}
//
// 	// Close store connections
// 	if err := a.store.Close(ctx); err != nil {
// 		stopErrors = append(stopErrors,
// 			fmt.Errorf("failed to close store: %w", err))
// 	}
//
// 	if len(stopErrors) > 0 {
// 		return fmt.Errorf("stop errors: %v", stopErrors)
// 	}
//
// 	return nil
// }
