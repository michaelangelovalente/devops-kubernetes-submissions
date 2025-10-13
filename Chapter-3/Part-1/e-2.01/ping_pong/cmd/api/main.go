package main

import (
	"common/boot"
	common_server "common/server"
	"log"

	"ping_pong/internal/app"
	"ping_pong/internal/server"
)

func main() {
	application, err := app.NewApplication()
	if err != nil {
		log.Fatalf("failed to create application")
	}

	srv := common_server.New(8092)
	srv.Handler = server.RegisterRoutes(application)

	boot.Run(application, srv)
}
