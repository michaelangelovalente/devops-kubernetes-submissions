package main

import (
	"common/boot"
	"log"
	"time"

	common_server "common/server"
	"todo_app/internal/server"
)

func imageRotationTask(logger *log.Logger, server *server.AppServer) {
	// Fetch the image on startup.
	if err := server.PicsumClient.FetchAndSaveImage(); err != nil {
		logger.Printf("Error fetching initial image: %v", err)
	}

	// Use a ticker to fetch the image every 1 minute.
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := server.PicsumClient.FetchAndSaveImage(); err != nil {
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

	// Start the background image rotation task.
	go imageRotationTask(server.Logger, server)

	boot.Run(server, httpServer)
}
