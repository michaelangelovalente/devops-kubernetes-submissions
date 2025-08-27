package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"ping_pong/internal/app"
	"ping_pong/internal/routes"
)

func main() {
	// TODO set up logger slog here or use default logger and app logger??

	// Root context, the parent context of all contexts
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// run application and handle exit
	if err := run(ctx); err != nil {
		// log error with slog?
		log.Printf("Server forced to shutdown with error: %v", err)
		os.Exit(1)
	}
}

// run contains the actual application logic
func run(ctx context.Context) error {
	// TODO: replace app with server
	app, err := app.NewApplication(ctx)
	if err != nil {
		return fmt.Errorf("startup error: failed to create Application: %w", err)
	}

	// Start actual application (bg workers can be added here)
	if err := app.Start(ctx); err != nil {
		return fmt.Errorf("startup error: failed to start application")
	}

	// HTTP Server with the application
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8090
		log.Print("PORT environment variable not detected, using defalt port %d\n", port)
	}
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      routes.RegisterRoutes(app),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Channel specific for server errors
	serverErrors := make(chan error, 1)

	// Start Http server in a go routine
	go func() {
		// slog log...
		log.Printf("Starting Application HTTP server on port %d\n", srv.Addr)
	}()

	// Wait for either:
	// - ctx cancellation
	// - server error
	select {
	case <-ctx.Done():
		log.Printf("shutdown signal received")
		return shutdown(srv, app)
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("server failed: %w", err)
	}
}

// shutdown handles graceful shutdown of all components
func shutdown(srv *http.Server, app *app.Application) error {
	// Create ctx specific for timeout shutdown operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Create error channel to collect shutdown errors
	shutdownErrors := make(chan error, 2)

	// Shutdown HTTP server
	go func() {
		log.Printf("shutting down Application HTTP server")
		shutdownErrors <- srv.Shutdown(ctx)
	}()

	// Stop Application Services
	go func() {
		log.Printf("stopping application services")
		shutdownErrors <- app.Stop(ctx)
	}()

	// Wait for both shutdowns or timeouts
	var errs []error
	for i := 0; i < 2; i++ {
		if err := <-shutdownErrors; err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("shutdown errors: %v", errs)
	}

	log.Printf("graceful shutdown cmplete")
	return nil
}
