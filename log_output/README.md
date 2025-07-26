# 1.1 Project log_output

## Description

**Docker Hub repository:** [log-output](https://hub.docker.com/layers/michaelangelovalente/log-output/v1.0/images/sha256-d6f49f4b33d256e739cf9e8a35320af7f6a87dee84d324de739eee89ae549d54)

**Source Code:** [log output source code]()

The `log_output` application is a Go-based HTTP server that generates timestamped log entries every 5 seconds. Each log entry contains a random UUID string that is stored in memory and displayed with a timestamp.

This application was built using modern Go 1.24+ practices and serves as a foundation for backend REST applications in subsequent exercises. The codebase is "complex" based on exercise requirement, but was done so to study clean architecture patterns with the Go standard `net/http` library.

## Commands

### Docker Image Management

#### Build Docker Image
```bash
docker build -t <username>/log-output:v1.0 .
```

#### Push to Docker Hub
```bash
docker push <username>/log-output:v1.0
```

### Kubernetes Deployment

#### Create k3d Cluster
```bash
k3d cluster create log-output-cluster -a 2
```

#### Configure kubectl Context
```bash
kubectl config get-contexts
kubectl config set-context log-output-cluster
```

#### Deploy Application
```bash
kubectl create deployment log-output-ex-1-1 --image=michaelangelovalente/log-output:v1.0
```

#### Monitor Deployment
```bash
kubectl get deployments
kubectl logs -f <pod-name>
```

## Local Development

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes.

### Build and Run
```bash
# Build the application
make build

# Run the application locally
make run

# Start with Docker Compose
make docker-run

# Stop Docker containers
make docker-down

# Live reload during development
make watch

# Clean build artifacts
make clean
```

### Direct Go Commands
```bash
# Run the server directly
go run cmd/server/main.go

# Build manually
go build -o main cmd/server/main.go

# Run tests (when available)
go test ./... -v
```

<!-- Future: Run build make command with tests -->
<!-- make all -->
