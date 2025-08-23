package server

import (
	"context"
	"fmt"
	"log"
	"log_output/internal/logger"
	"log_output/internal/storage"
	"log_output/internal/utils"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	*http.Server
	port       int
	logger     *logger.Logger
	logStorage logger.LogStorage
	loggerDone chan struct{}
	wg         sync.WaitGroup
}

func NewServer() *Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 3000
		log.Printf("Environment variable PORT not detected, starting application on default port %d\n", port)
	}

	// Initialize store mem
	store := storage.NewMemoryStorage()

	interval := 5 * time.Second
	// normally would get interval from env with godotenv...

	config := logger.LoggerConfig{
		Interval:      interval,
		GenerateValue: logger.GenerateUUID, // parameterized func
		TimeFormat:    time.RFC3339,        // ISO 8601 format with timezone
	}

	// Init Logger instnace
	log := logger.NewLogger(config, store)
	newServer := &Server{
		Server: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
			// Handler:      newServer.RegisterRoutes(),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
		port:       port,
		logger:     log,
		logStorage: store,
		loggerDone: make(chan struct{}),
	}

	newServer.Handler = newServer.RegisterRoutes()

	// Start bg logger
	newServer.startLogger()

	return newServer
}

// Basic health and ready handlers
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	// simple health check
	utils.WriteJson(w, http.StatusOK, utils.Envelope{"status": "ready"})
}

func (s *Server) ReadyHandler(w http.ResponseWriter, r *http.Request) {

}

// Override orig. httpServer.Shutdown() to clean up custom logger
func (s *Server) Shutdown(ctx context.Context) error {
	// Stop logger first
	close(s.loggerDone)
	s.wg.Wait()

	// Original httpServer Shutdown call
	return s.Server.Shutdown(ctx)
}

func (s *Server) startLogger() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()

		// Create cancellable ctx
		ctx, cancel := context.WithCancel(context.Background())

		// sart diff goroutine used to wait for done signal
		go func() {
			<-s.loggerDone
			cancel()
		}()

		// Start logger --> will continue running until context cancelled
		if err := s.logger.StartLogger(ctx); err != nil {
			// Logg err if start logger gives err  but continue running server
			fmt.Printf("Logger error: %v", err)
		}

	}()
}
