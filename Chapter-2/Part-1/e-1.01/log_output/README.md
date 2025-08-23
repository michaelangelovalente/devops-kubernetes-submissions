# Log Output - Exercise 1.1

A Go-based HTTP server with background logging for the DevOps Kubernetes Course.

## 🔗 Links

- **Docker Hub:** [michaelangelovalente/log-output:ex1.01](https://hub.docker.com/layers/michaelangelovalente/log_output-app/ex1.01/images/sha256-7014f94e2d4aee68a60e4c21c19130ee5e5f192136e2bbcda71e4fa6adce177b)
- **Source Code:** [GitHub Repository](https://github.com/michaelangelovalente/devops-kubernetes-submissions/tree/main/Chapter-2/Part-1/e-1.01/log_output)

## 📋 Description

The `log_output` application is a Go-based HTTP server that generates timestamped log entries every 5 seconds. Each entry contains a random UUID string stored in memory and displayed with timestamps.

Built with modern **Go 1.24+** practices following clean architecture patterns, this application demonstrates:
- Background goroutines with graceful shutdown
- In-memory storage with interfaces
- CORS-enabled HTTP server
- Structured logging and UUID generation

> **Note:** The codebase intentionally exceeds exercise requirements to showcase clean architecture patterns using Go's standard `net/http` library.

## 🐳 Docker Commands

### Build Image
```bash
docker build -t <username>/log-output:v1.0 .
```

### Push to Docker Hub
```bash
docker push <username>/log-output:v1.0
```

## ☸️ Kubernetes Deployment

### Create k3d Cluster
```bash
k3d cluster create first-deploy-cluster -a 2
```

### Configure kubectl Context
```bash
kubectl config get-contexts
kubectl config use-context first-deploy-cluster
```

### Deploy Application
```bash
kubectl create deployment log-output-ex-1-1 --image=michaelangelovalente/log-output:v1.0
```

### Monitor Deployment
```bash
kubectl get deployments
kubectl get pods
kubectl logs -f <pod-name>
```

## 🛠️ Local Development

### Using Makefile Commands
```bash
make build       # Build the application
make run         # Run locally (default port: 3000)
make docker-run  # Start with Docker Compose
make docker-down # Stop Docker containers
make watch       # Live reload during development (uses air)
make clean       # Clean build artifacts
```

### Direct Go Commands
```bash
go run cmd/server/main.go              # Run server directly
go build -o main cmd/server/main.go    # Build manually
go test ./... -v                       # Run tests
```

---

**Environment Variables:**
- `PORT`: Server port (default: 3000)

**Features:**
- Background UUID logging every 5 seconds
- CORS-enabled HTTP endpoints
- Graceful shutdown with context cancellation
