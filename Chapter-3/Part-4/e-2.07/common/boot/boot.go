package boot

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

type App interface {
	Start(ctx context.Context) error
	Stop() error
}

func Run(app App, srv *http.Server) {
	appCtx, appCancel := context.WithCancel(context.Background())
	defer appCancel()

	done := make(chan bool, 1)
	go gracefulShutdown(srv, app, appCancel, done)

	if err := app.Start(appCtx); err != nil {
		log.Fatalf("failed to start application: %v", err)
	}

	log.Printf("Server started on port: %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("http server error: %s", err)
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

func gracefulShutdown(apiServer *http.Server, app App, appCancel context.CancelFunc, done chan bool) {
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
