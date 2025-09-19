package main

import (
	"context"
	"fmt"
	"log"
	"log_output/internal/app"
	"log_output/internal/server"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func gracefulShutdown(apiServer *http.Server, app *app.Application, appCancel context.CancelFunc, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop() // Allow Ctrl+C to force shutdown

	// Shutdown sequence: 1) Stop application services, 2) Stop HTTP server
	log.Println("stopping application services...")
	appCancel() // Cancel application context to stop logger
	if err := app.Stop(); err != nil {
		log.Printf("Application shutdown error: %v", err)
	}

	log.Println("stopping HTTP server...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

func main() {
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	// Create application context for logger
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	// Start application services (logger)
	if err = app.Start(appCtx); err != nil {
		panic(fmt.Sprintf("failed to start application: %v", err))
	}

	server := server.NewServer(app)
	log.Printf("Server 'log_writer' started on port: %d\n", server.Port)
	done := make(chan bool, 1)

	// Start graceful shutdown handler
	go gracefulShutdown(server.HttpServer, app, appCancel, done)

	log.Println("Starting HTTP server...")
	err = server.HttpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}
