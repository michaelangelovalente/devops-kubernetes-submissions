# Log Output App - Exercise 1.07

A Go-based HTTP server application that generates timestamps with random strings and provides external access through Kubernetes Ingress.

## 🔗 Links

- **Docker Hub:** [michaelangelovalente/log_output:ex1.07](https://hub.docker.com/layers/michaelangelovalente/log_output-img/ex1.07/images/sha256-9a16e8050cfc228ff4bdfc0ec7a5faed689ee7420bda17a941a1287b4c155057)
- **Source Code:** [GitHub Repository](https://github.com/michaelangelovalente/devops-kubernetes-submissions/tree/main/Chapter-2/Part-3/e-1.07/log_output)

## 📋 Description

- **Backend:** Go HTTP server with structured logging and graceful shutdown
- **Background Logger:** Continuously generates timestamps and random UUID strings every 5 seconds
- **Status Endpoint:** Provides `/` endpoint to retrieve current status (timestamp + random string)
- **Infrastructure:** Docker containerization and Kubernetes deployment with Ingress for external access

This application demonstrates how to use Kubernetes Ingress to enable external HTTP access to services running in a cluster, fulfilling exercise 1.07 requirements.

## 🐳 Docker Commands

### Build Image
```bash
docker build -t <username>/log_output:ex1.07 .
```

### Push to Docker Hub
```bash
docker push <username>/log_output:ex1.07
```

### Run Locally with Docker
```bash
docker run -p 8080:8080 <username>/log_output:ex1.07
```

## ☸️ Kubernetes Deployment

### Create k3d Cluster
```bash
k3d cluster create --port 8082:30080@agent:0 -p 8081:80@loadbalancer --agents 2
```

### Configure kubectl Context
```bash
kubectl config get-contexts
kubectl config use-context k3d-k3s-default
```

### Deploy Application
```bash
# Apply all manifests (deployment, service, ingress)
kubectl apply -f manifests/
```

### Verify Deployment
```bash
# Check deployment status
kubectl get deployments

# Check pods
kubectl get pods

# Check services and ingress
kubectl get svc,ing

# Check logs
kubectl logs -f deployment/log-output-dep
```

### Access Application

**Method 1: Ingress (Primary for Exercise 1.07)**
```bash
# Access via Ingress through load balancer
# If using k3d with port mapping 8081:80@loadbalancer
curl localhost:8081/
curl localhost:8081/status
curl localhost:8081/status?n=5

# Or access directly in browser:
# http://localhost:8081/
# http://localhost:8081/status
```

**Method 2: Port Forward (Alternative)**
```bash
# Port forward to access the application (development/debugging)
kubectl port-forward service/log-output-svc 3003:2345

# Test the application
curl localhost:3003/
curl localhost:3003/status
curl localhost:3003/status?n=5
```

#### Expected Output

**GET /** - Hello World endpoint:
```json
{
  "message": "Hello World"
}
```

**GET /status** - Status endpoint with logs:
```json
{
  "status": "ready",
  "logs": [
    {
      "timestamp": "2025-08-25T11:36:58.538054382Z",
      "value": "2ae58e2b-effe-401a-b39d-f9b4be711aab"
    }
  ]
}
```

## 🛠️ Local Development

### Using Makefile Commands
```bash
make build       # Build the application
make run         # Run locally on port 8080
make docker-run  # Start with Docker Compose
make docker-down # Stop Docker containers
make watch       # Live reload during development (uses Air)
make clean       # Clean build artifacts
```

### Direct Go Commands
```bash
go run cmd/server/main.go              # Run server directly
go build -o main cmd/server/main.go    # Build manually
go test ./... -v                       # Run tests
```

### Development Features
- **Background Logging:** Automatic timestamp and UUID generation every 5 seconds
- **In-Memory Storage:** Random string persists until process restart
- **Graceful Shutdown:** Clean server termination handling
- **Structured Logging:** JSON formatted logs with timestamps

---

## ⚙️ Configuration

**Environment Variables:**
- `PORT`: Server port (default: 8080)

**Kubernetes Resources:**
- **Deployment:** Manages application pods and replica sets
- **Service Type:** ClusterIP (internal cluster communication)
- **Service Port:** 2345 (cluster-internal port)
- **Target Port:** 8080 (application port)
- **Ingress:** Provides external HTTP access on port 80

**Network Configuration:**
- **External Access:** Available via browser through Ingress controller
- **Internal Access:** Available to other cluster services via ClusterIP service on port 2345
- **Application Port:** Listens on port 8080 inside the container

## 📊 API Endpoints

- `GET /`: Returns "Hello World" message in JSON format
- `GET /status`: Returns application status with stored log entries
- `GET /status?n=5`: Returns application status with last 5 log entries (default: 10)
