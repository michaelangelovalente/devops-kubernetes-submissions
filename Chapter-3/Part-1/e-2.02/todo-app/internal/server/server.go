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
	Logger       *log.Logger
	ImageDir     string
	TodoClient   *client.TodoClient
	PicsumClient *client.PicsumClient
}

func NewServer() *AppServer {

	imageDir := os.Getenv("IMAGE_DIR")
	if imageDir == "" {
		log.Fatal("IMAGE_DIR environment variable is not set")
	}

	picsumUrl := os.Getenv("PICSUM_URL")
	if picsumUrl == "" {
		log.Fatalf("PICSUM_URL env variable is not set")
	}
	picumsClient := client.NewPicsumClient(picsumUrl, imageDir, time.Second*5)

	todoUrl := os.Getenv("TODO_SVC_URL")
	if todoUrl == "" {
		log.Fatalf("TODO_SVC_URL env variable is not set")
	}
	todoClient := client.NewTodoClient(todoUrl, time.Second*5)

	logger := log.New(os.Stdout, "[LOGGER] ", log.LstdFlags)

	NewServer := &AppServer{
		Logger:       logger,
		ImageDir:     imageDir,
		TodoClient:   todoClient,
		PicsumClient: picumsClient,
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
