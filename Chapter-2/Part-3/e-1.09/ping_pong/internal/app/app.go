// package app equiv to server in other ms
package app

import (
	"context"
	"fmt"
	"log/slog"
	handler "ping_pong/internal/api"
	"ping_pong/internal/store"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Application struct {
	router *chi.Mux
	store  *store.Store
	// services *Services
	// handlers *Handlers
	PingpongHandler *handler.PingPongHandler
	// Logger *logger.Logger

	wg sync.WaitGroup

	// Lifecycle context management
	baseCtx context.Context
}

func NewApplication(ctx context.Context) (*Application, error) {
	// // Create Store Layer
	//
	// // Create services layer
	// services := &Services{
	// 	Log:  service.NewLogService(store),
	// 	Ping: service.NewPingService(store),
	// }
	//
	// // Create handlers layer
	// handlers := &Handlers{
	// 	Log:  handler.NewLogHandler(services.Log),
	// 	Ping: handler.NewPingHandler(services.Ping),
	// }
	//
	// app := &Application{
	// 	Config:   cfg,
	// 	router:   chi.NewRouter(),
	// 	store:    store,
	// 	services: services,
	// 	handlers: handlers,
	// 	workers:  []Worker{},
	// 	baseCtx:  ctx,
	// }
	//
	// // Initialize any background workers
	// app.initializeWorkers()
	//
	// // Setup HTTP routes
	// app.setupRoutes()
	return nil, nil
}

// func (a *Application) setupRoutes() {
//     r := a.router
//
//     // Global middleware
//     r.Use(middleware.RequestID)
//     r.Use(middleware.RealIP)
//     r.Use(middleware.Logger)
//     r.Use(middleware.Recoverer)
//
//     // Add context timeout middleware for all requests
//     r.Use(func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//             // Create a timeout context for each request
//             ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
//             defer cancel()
//
//             next.ServeHTTP(w, r.WithContext(ctx))
//         })
//     })
//
//     // Health check endpoints (no timeout needed)
//     r.Get("/health/live", a.handlers.Log.HealthLive)
//     r.Get("/health/ready", a.handlers.Log.HealthReady)
//
//     // API routes
//     r.Route("/api/v1", func(r chi.Router) {
//         r.Get("/logs", a.handlers.Log.GetCurrentLog)
//         r.Get("/pingpong", a.handlers.Ping.HandlePing)
//     })
// }

// func (a *Application) Router() http.Handler {
// 	return a.router
// }

func (a *Application) Start(ctx context.Context) error {
	// Start all background workers with their own context
	for _, worker := range a.workers {
		w := worker // Capture loop variable

		a.wg.Add(1)
		go func() {
			defer a.wg.Done()

			// Each worker gets a child context
			if err := w.Start(ctx); err != nil {
				slog.Error("worker failed",
					"worker", w.Name(),
					"error", err)
			}
		}()
	}

	return nil
}

func (a *Application) Stop(ctx context.Context) error {
	// Signal all workers to stop
	var stopErrors []error

	for _, worker := range a.workers {
		if err := worker.Stop(ctx); err != nil {
			stopErrors = append(stopErrors,
				fmt.Errorf("failed to stop %s: %w", worker.Name(), err))
		}
	}

	// Wait for all workers to finish with timeout
	done := make(chan struct{})
	go func() {
		a.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("all workers stopped successfully")
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout while waiting for workers")
	}

	// Close store connections
	if err := a.store.Close(ctx); err != nil {
		stopErrors = append(stopErrors,
			fmt.Errorf("failed to close store: %w", err))
	}

	if len(stopErrors) > 0 {
		return fmt.Errorf("stop errors: %v", stopErrors)
	}

	return nil
}
