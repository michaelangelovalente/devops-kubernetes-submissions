package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"todo-backend/internal/api"

	_ "github.com/joho/godotenv/autoload"
)

type AppServer struct {
	logger      *log.Logger
	todoHandler *api.TodoHandler
}

func NewServer() *AppServer {
	logger := log.New(os.Stdout, "[LOGGER] ", log.LstdFlags)
	todoHander := api.NewTodoHandler(logger)

	appServer := &AppServer{
		logger:      logger,
		todoHandler: todoHander,
	}

	return appServer
}

func (appS *AppServer) Start(ctx context.Context) error {
	fmt.Println("Starting application services....")
	return nil
}

func (appS *AppServer) Stop() error {
	fmt.Println("Stopping application services...")
	return nil
}
