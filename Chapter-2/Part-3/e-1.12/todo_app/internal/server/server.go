package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	*http.Server
	Port   int
	Logger *log.Logger
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	if port == 0 {
		port = 3005
		logger.Printf("No port env variable detected. Running on default port %d\n", port)
	}

	NewServer := &Server{
		Port:   port,
		Logger: logger,
		Server: &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		},
	}

	server :=
	NewServer.Server = server

	return NewServer
}
