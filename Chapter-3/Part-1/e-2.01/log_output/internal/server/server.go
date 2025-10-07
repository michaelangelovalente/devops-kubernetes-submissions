package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"log_output/internal/app"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	Port       int
	app        *app.Application
	HTTPServer *http.Server
}

func NewServer(application *app.Application) *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8091
		log.Printf("No PORT environment variable detected, starting server on default port: %d\n", port)
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      RegisterRoutes(application),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return &Server{
		Port:       port,
		HTTPServer: server,
	}
}
