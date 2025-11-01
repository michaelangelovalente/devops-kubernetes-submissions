package main

import (
	"common/boot"
	common_server "common/server"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"todo-backend/internal/server"
)

func main() {

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	healthCheck := flag.Bool("health-check", false, "Run health check and exit")
	flag.Parse()

	if *healthCheck {
		log.Println("Performing health check")
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/health", port))
		if err != nil || resp.StatusCode != 200 {
			log.Println("Health Check failed")
			os.Exit(1)
		}

		log.Println("Health Check passed")
		os.Exit(0)
	}
	appServer := server.NewServer()

	httpServer := common_server.New(8088)
	httpServer.Handler = appServer.RegisterRoutes()

	boot.Run(appServer, httpServer)

}
