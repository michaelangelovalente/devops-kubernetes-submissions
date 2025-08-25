# Log Output App - Exercise 1.07

A Go-based HTTP server application that generates timestamps with random strings and provides external access through Kubernetes Ingress.

## 🔗 Links

- **Docker Hub:** [michaelangelovalente/log_output:ex1.07](https://hub.docker.com/r/michaelangelovalente/log_output)
- **Source Code:** [GitHub Repository](https://github.com/your-username/devops-kubernetes-submissions)

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
docker run -p 3000:3000 <username>/log_output:ex1.07
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

# Or access directly in browser:
# http://localhost:8081/
```

**Method 2: Port Forward (Alternative)**
```bash
# Port forward to access the application (development/debugging)
kubectl port-forward service/log-output-svc 3003:2345

# Test the application
curl localhost:3003/
```

#### Expected Output

```json
{
  "message": "Hello World",
  "timestamp": "2024-01-15T10:30:45Z",
  "random_string": "a1b2c3d4-e5f6-7890-abcd-ef1234567890"
}
```

## 🛠️ Local Development

### Using Makefile Commands
```bash
make build       # Build the application
make run         # Run locally on port 3000
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
- `PORT`: Server port (default: 3000)

**Kubernetes Resources:**
- **Deployment:** Manages application pods and replica sets
- **Service Type:** ClusterIP (internal cluster communication)
- **Service Port:** 2345 (cluster-internal port)
- **Target Port:** 3000 (application port)
- **Ingress:** Provides external HTTP access on port 80

**Network Configuration:**
- **External Access:** Available via browser through Ingress controller
- **Internal Access:** Available to other cluster services via ClusterIP service on port 2345
- **Application Port:** Listens on port 3000 inside the container

## 📊 API Endpoints

- `GET /`: Returns current status with timestamp and random string in JSON format