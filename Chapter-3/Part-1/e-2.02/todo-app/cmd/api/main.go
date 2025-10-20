package main

import (
	"log"
	"time"

	"common/boot"
	common_server "common/server"
	"todo_app/internal/picsum"
	"todo_app/internal/server"
)

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
	httpServer := common_server.New(3010)
	httpServer.Handler = server.RegisterRoutes()
	boot.Run(server, httpServer)

	// Start the background image rotation task.
	go imageRotationTask(server.Logger, server.ImageDir)

}
