package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"todo_app/internal/picsum"
	"todo_app/internal/server"
)

// const imageDir = "tmp/images"

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")

	done <- true
}

func imageRotationTask(logger *log.Logger, imageDir string) {
	// Fetch the image on startup.
	if err := picsum.FetchAndSaveImage(imageDir); err != nil {
		logger.Printf("Error fetching initial image: %v", err)
	}

	// Use a ticker to fetch the image every 1 minute.
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := picsum.FetchAndSaveImage(imageDir); err != nil {
			logger.Printf("Error fetching image: %v", err)
		} else {
			logger.Println("Background image updated successfully.")
		}
	}
}

func main() {
	server := server.NewServer()

	// Start the background image rotation task.
	go imageRotationTask(server.Logger, server.ImageDir)

	done := make(chan bool, 1)

	go gracefulShutdown(server.Server, done)

	server.Logger.Printf("Server started on port %d\n", server.Port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	server.Logger.Println("Graceful shutdown complete")
}
