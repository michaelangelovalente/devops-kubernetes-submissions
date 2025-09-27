package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"ping_pong/internal/app"
	"ping_pong/internal/routes"
)

// loadEnv searches for and loads a .env file in the current or parent directories.
func loadEnv() {
	// Start from the current directory
	dir, err := os.Getwd()
	if err != nil {
		log.Println("Error getting current directory:", err)
		return
	}

	// Look for .env file in current and parent directories
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			err := godotenv.Load(envPath)
			if err != nil {
				log.Printf("Error loading .env file from %s: %v", envPath, err)
			} else {
				log.Printf("Loaded .env file from %s", envPath)
				return
			}
		}

		// Move to parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the root directory
			log.Println("No .env file found in current or parent directories.")
			return
		}
		dir = parent
	}
}

func main() {
	loadEnv()
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	// run application and handle exit
	if err := run(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
		os.Exit(1)
	}
}

// run contains the actual application logic
func run(ctx context.Context) error {
	app, err := app.NewApplication(ctx)
	if err != nil {
		return fmt.Errorf("startup error: failed to create Application: %w", err)
	}

	// Start actual application (bg workers can be added here)
	if err := app.Start(ctx); err != nil {
		return fmt.Errorf("startup error: failed to start application")
	}

	// HTTP Server with the application
	port, _ := strconv.Atoi(os.Getenv("PING_PONG_PORT"))
	if port == 0 {
		port = 8092
		log.Printf("PING_PONG_PORT environment variable not detected, using defalt port %d\n", port)
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
		log.Printf("Starting Application HTTP server on port %s\n", srv.Addr)
		serverErrors <- srv.ListenAndServe()
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
	for range len(errs) {
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