# Integrating `log_output` with `ping_pong`: A Guide

This guide details the process of integrating the `log_output` service with the `ping_pong` service, following idiomatic Go best practices.
The goal is for `log_output` to fetch the current count from `ping_pong` via a REST API call and include this count in its log entries.

## Core Principles

We will adhere to the following software design principles:

1.  **Separation of Concerns**: The logic for communicating with the `ping_pong` service will be isolated into its own package.
    This makes the code easier to understand, maintain, and test.
2.  **Dependency Injection (DI)**: Instead of creating dependencies inside the components that need them, we will provide them from the outside.
    The `log_output` application will receive a *client* to communicate with `ping_pong`.
3.  **Programming to an Interface**: We will define an interface for the `ping_pong` client.
    The rest of the application will depend on this interface, not on a concrete implementation.
    This allows us to easily swap out the client implementation (e.g., for a mock client in tests) without changing the business logic.
4.  **Configuration over Hardcoding**: The address of the `ping_pong` service will be managed via environment variables, not hardcoded in the source code.
    This is essential for running the application in different environments (development, staging, production).

---

## Step-by-Step Implementation

### Step 1: Update Configuration

We need a way to tell the `log_output` service where to find the `ping_pong` service. We'll use an environment variable for this.

**Why?** Hardcoding URLs makes your application brittle. Using environment variables allows you to reconfigure your application for different environments without changing a single line of code.

1.  **Update `.env`**: Add the URL for the `ping_pong` service.

    ```dotenv
    # .env
    LOG_OUTPUT_PORT=8091
    PING_PONG_PORT=8092
    PING_PONG_SVC_URL=http://localhost:8092
    ```

2.  **Update `docker-compose.yml`**: Ensure the `log_output` service can communicate with `ping_pong` and pass the environment variable.

    ```yaml
    # docker-compose.yml
    services:
      log-output:
        image: log_output_img
        container_name: log_output_ctr
        build:
          context: .
          dockerfile: ./log_output/Dockerfile
        restart: unless-stopped
        ports:
          - "${LOG_OUTPUT_PORT}:${LOG_OUTPUT_PORT}"
        environment:
          - PORT=${LOG_OUTPUT_PORT}
          - PING_PONG_SVC_URL=http://ping-pong:${PING_PONG_PORT} # Use service name for discovery

      ping-pong:
        image: ping_pong_img
        container_name: ping_pong_ctr
        build:
          context: .
          dockerfile: ./ping_pong/Dockerfile
        restart: unless-stopped
        ports:
          - "${PING_PONG_PORT}:${PING_PONG_PORT}"
        environment:
          - PORT=${PING_PONG_PORT}
    ```
    *Note*: In Docker Compose, services can reach each other using their service name (`ping-pong`) as the hostname.

### Step 2: Create the PingPong Client

We will create a dedicated client inside the `log_output` service to handle all communication with the `ping_pong` API.

**Why?** This encapsulates the logic for making HTTP requests, handling errors, and decoding responses related to the `ping_pong` service. The rest of our application doesn't need to know these details; it just needs a count.

1.  **Create a new directory**: `log_output/internal/client/pingpong`
2.  **Create `log_output/internal/client/pingpong/client.go`**:

    ```go
    package pingpong

    import (
    	"context"
    	"encoding/json"
    	"fmt"
    	"net/http"
    	"time"
    )

    // Client defines the interface for interacting with the ping_pong service.
    type Client interface {
    	// GetCount fetches the current count from the ping_pong service.
    	GetCount(ctx context.Context) (int, error)
    }

    // httpCient implements the Client interface using HTTP.
    type httpClient struct {
    	client  *http.Client
    	baseURL string
    }

    // NewClient creates a new configured HTTP client for the ping_pong service.
    func NewClient(baseURL string, timeout time.Duration) Client {
    	return &httpClient{
    		client: &http.Client{
    			Timeout: timeout,
    		},
    		baseURL: baseURL,
    	}
    }

    // GetCount performs the HTTP GET request to the /pingpong endpoint.
    func (c *httpClient) GetCount(ctx context.Context) (int, error) {
    	url := fmt.Sprintf("%s/pingpong", c.baseURL)
    	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    	if err != nil {
    		return 0, fmt.Errorf("failed to create request: %w", err)
    	}

    	resp, err := c.client.Do(req)
    	if err != nil {
    		return 0, fmt.Errorf("failed to execute request: %w", err)
    	}
    	defer resp.Body.Close()

    	if resp.StatusCode != http.StatusOK {
    		return 0, fmt.Errorf("received non-OK status: %s", resp.Status)
    	}

    	var payload struct {
    		Count int `json:"count"`
    	}
    	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
    		return 0, fmt.Errorf("failed to decode response body: %w", err)
    	}

    	return payload.Count, nil
    }
    ```

### Step 3: Update the Data Model

The log entry must now include the ping-pong count.

**Why?** The data model (struct) must reflect the data we intend to store and display.

1.  **Modify `log_output/internal/store/memory_store.go`**:

    ```go
    package store

    import (
    	"sync"
    	"time"
    )

    type LogStorage interface {
        // Store now accepts a pingPongCount
    	Store(timestamp time.Time, value string, pingPongCount int) error
    	GetAll() []LogEntry
    	GetLatest(n int) []LogEntry
    }

    // LogEntry now includes the PingPongCount
    type LogEntry struct {
    	Timestamp     time.Time `json:"timestamp"`
    	Value         string    `json:"value"`
    	PingPongCount int       `json:"pingPongCount"`
    }

    type MemoryStorage struct {
    	entries []LogEntry
    	mu      sync.RWMutex
    }

    func NewMemoryStorage() *MemoryStorage {
        return &MemoryStorage{
            entries: make([]LogEntry, 0),
        }
    }

    // Store method updated to save the new field
    func (m *MemoryStorage) Store(timestamp time.Time, value string, pingPongCount int) error {
    	m.mu.Lock()
    	defer m.mu.Unlock()

    	m.entries = append(m.entries, LogEntry{
    		Timestamp:     timestamp,
    		Value:         value,
    		PingPongCount: pingPongCount,
    	})

    	return nil
    }
    // ... (GetAll and GetLatest remain the same)
    ```

### Step 4: Integrate the Client into the Application

Now we wire the new client into the `log_output` application logic.

**Why?** This is where Dependency Injection happens. We connect the components, providing the logger with the client it needs to do its job.

1.  **Modify `log_output/internal/logger/logger.go`**:

    ```go
    package logger

    import (
    	"context"
    	"fmt"
    	"log"
    	"os"
    	"sync"
    	"time"

    	"github.com/google/uuid"
    	"log_output/internal/client/pingpong"
    	"log_output/internal/store"
    )

    // ... (ValueGenerator and LoggerConfig are unchanged)

    type Logger struct {
    	loggerConfig  LoggerConfig
    	logStorage    store.LogStorage
    	pingpongClient pingpong.Client // Add the client interface
    	currentValue  string
    	rwMutex       sync.RWMutex
    	normalLogger  *log.Logger
    }

    // NewLogger now accepts the pingpong.Client
    func NewLogger(cfg LoggerConfig, storage store.LogStorage, ppClient pingpong.Client) *Logger {
    	return &Logger{
    		loggerConfig:   cfg,
    		logStorage:     storage,
    		pingpongClient: ppClient, // Store the client
    		normalLogger:   log.New(os.Stdout, "[LOGGER] ", log.LstdFlags),
    	}
    }

    // ... (StartLogger is mostly the same)

    func (l *Logger) logCurrent(ctx context.Context) {
    	l.rwMutex.RLock()
    	value := l.currentValue
    	l.rwMutex.RUnlock()

        // Fetch the count from the ping_pong service
    	count, err := l.pingpongClient.GetCount(ctx)
    	if err != nil {
    		l.normalLogger.Printf("Error fetching ping-pong count: %v", err)
            // Decide on error handling: we can log with a zero value or skip logging
            // For this example, we'll log with a -1 to indicate an error
            count = -1
    	}

    	timestamp := time.Now()

        // Store the log with the count
    	if err := l.logStorage.Store(timestamp, value, count); err != nil {
    		l.normalLogger.Printf("Error storing log: %v", err)
    		return
    	}

    	fmt.Printf("%s: %s, PingPongCount: %d", timestamp.Format(l.loggerConfig.TimeFormat), value, count)
    }

    // Update the loop in StartLogger to pass the context
    func (l *Logger) StartLogger(ctx context.Context) error {
        // ... (setup is the same)

        l.logCurrent(ctx) // Initial log

        for {
            select {
            case <-ctx.Done():
                l.normalLogger.Println("Logger stopped...")
                return nil
            case <-ticker.C:
                l.logCurrent(ctx) // Log on interval
            }
        }
    }
    ```

2.  **Modify `log_output/internal/app/app.go`** to create and inject the client:

    ```go
    package app

    import (
    	"context"
    	"fmt"
    	"os"
    	"sync"
    	"time"

    	"log_output/internal/api"
    	"log_output/internal/client/pingpong"
    	"log_output/internal/logger"
    	"log_output/internal/store"
    )

    type Application struct {
    	Logger           *logger.Logger
    	wg               sync.WaitGroup
    	LogMemoryHandler *api.LoggerEntryHandler
        PingPongClient   pingpong.Client // Add the client
    }

    func NewApplication() (*Application, error) {
        // --- Configuration ---
        pingPongURL := os.Getenv("PING_PONG_SVC_URL")
        if pingPongURL == "" {
            return nil, fmt.Errorf("PING_PONG_SVC_URL environment variable not set")
        }

    	// --- Client layer ---
        pingPongClient := pingpong.NewClient(pingPongURL, 5*time.Second)

    	// --- Store layer ----
    	logMemoryStore := store.NewMemoryStorage()
    	loggerConfig := logger.LoggerConfig{
    		Interval:   5 * time.Second,
    		TimeFormat: time.RFC3339,
    	}
        // Inject the client into the logger
    	logMemory := logger.NewLogger(loggerConfig, logMemoryStore, pingPongClient)

    	//---  Handler layer ----
    	logMemoryHandler := api.NewLoggerEntryHandler(logMemoryStore, logMemory.GetNormalLogger())

    	app := &Application{
    		Logger:           logMemory,
    		LogMemoryHandler: logMemoryHandler,
            PingPongClient:   pingPongClient,
    	}
    	return app, nil
    }

    // ... (Start and Stop methods are unchanged)
    ```

### Summary

By following these steps, we have successfully integrated the two services in a way that is robust, testable, and scalable. The `log_output` service now depends on an abstract `pingpong.Client` interface,
receives its dependencies via DI, and is configured through the environment. This clean architecture makes the system easier to reason about and evolve over time.
