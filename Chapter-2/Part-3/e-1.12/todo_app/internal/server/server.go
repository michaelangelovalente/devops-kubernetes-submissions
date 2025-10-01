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
	Port     int
	Logger   *log.Logger
	ImageDir string
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	if port == 0 {
		port = 3005
		logger.Printf("No port env variable detected. Running on default port %d\n", port)
	}

	imageDir := os.Getenv("IMAGE_DIR")
	if imageDir == "" {
		imageDir = "tmp/images"
	}

	NewServer := &Server{
		Port:     port,
		Logger:   logger,
		ImageDir: imageDir,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	NewServer.Server = server

	return NewServer
}
