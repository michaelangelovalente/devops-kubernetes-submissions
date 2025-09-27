package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"log_reader/internal/app"
)

type Server struct {
	Port       int
	app        *app.Application
	HttpServer *http.Server
}

func NewServer(app *app.Application) *Server {
	port, _ := strconv.Atoi(os.Getenv("READER_PORT"))
	if port == 0 {
		port = 8092
		log.Printf("No READER_PORT environment variable detected, starting server on default port: %d\n", port)
	}

	NewServer := &Server{
		Port: port,
		app:  app,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.Port),
		Handler:      NewServer.RegisterRoutes(app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	NewServer.HttpServer = server

	return NewServer
}
