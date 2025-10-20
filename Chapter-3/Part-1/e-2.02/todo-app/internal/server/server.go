package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"todo_app/internal/client"

	_ "github.com/joho/godotenv/autoload"
)

type AppServer struct {
	Logger     *log.Logger
	ImageDir   string
	TodoClient *client.TodoClient
}

func NewServer() *AppServer {

	imageDir := os.Getenv("IMAGE_DIR")
	if imageDir == "" {
		log.Fatal("IMAGE_DIR environment variable is not set")
	}

	logger := log.New(os.Stdout, "[LOGGER] ", log.LstdFlags)

	todoUrl := "http://localhost:8080"
	todoClient := client.NewClient(todoUrl, time.Second*5)
	NewServer := &AppServer{
		Logger:     logger,
		ImageDir:   imageDir,
		TodoClient: todoClient,
	}

	return NewServer
}

func (as *AppServer) Start(ctx context.Context) error {
	fmt.Println("Starting application services....")
	return nil
}

func (as *AppServer) Stop() error {
	fmt.Println("Stopping application services...")
	return nil
}
