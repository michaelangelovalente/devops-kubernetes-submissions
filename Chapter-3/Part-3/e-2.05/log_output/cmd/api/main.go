package main

import (
	"common/boot"
	common_server "common/server"
	"log"

	"log_output/internal/app"
	"log_output/internal/server"
)

func main() {
	application, err := app.NewApplication()
	if err != nil {
		log.Fatalf("failed to create application")
	}

	srv := common_server.New(8091)
	srv.Handler = server.RegisterRoutes(application)

	boot.Run(application, srv)
}
