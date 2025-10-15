# Exercise 1.09 - More Services
Two Go applications with shared Kubernetes Ingress for the DevOps with Kubernetes Course.

## 🔗 Links

- **Log Output Docker Hub:** [michaelangelovalente/log_output-img:ex1.09](https://hub.docker.com/repository/docker/michaelangelovalente/log_output-img/tags/ex1.09/sha256-65344f2974d17f6a47bc80ada080d7df3e813834380d3252484236752a568153)
- **Ping Pong Docker Hub:** [michaelangelovalente/ping_pong-img:ex1.09](https://hub.docker.com/layers/michaelangelovalente/ping_pong-img/ex1.09/images/sha256-c0d4eec5bccff9f1efbc4e91a8fca6f2b5e5df75e4ec247f1ec0d3b48c19e5da)
- **Source Code:** [GitHub Repository](https://github.com/michaelangelovalente/devops-kubernetes-submissions/tree/main/Chapter-2/Part-3/e-1.09)

## 📋 Description
 This exercise demonstrates path-based routing with multiple services sharing a single Ingress resource as per requirement for exercise 1.09

### Log Output Application
- **Available Paths:** `/logs`, `/status`
- **Backend:** Go HTTP server with Chi router and structured logging
- **Logging:** Background UUID generation every 5 seconds with RFC3339 timestamps
- **Storage:** Thread-safe in-memory log storage
- **Response:** Returns structured JSON with logs and status information

### Ping Pong Application
- **Path:** `/pingpong`
- **Backend:** Go HTTP server with request counter
- **Functionality:** Responds with JSON containing current counter and increments counter
- **Storage:** In-memory counter (resets on pod restart)
- **Response:** JSON response with incremented counter value

## 🚀 Features

### Shared Infrastructure
- **Shared Ingress:** Single Ingress resource routing to both applications
- **Path-based Routing:** `/` → log_output, `/pingpong` → ping_pong
- **Kubernetes Services:** Individual ClusterIP services for each app
- **Docker Images:** Separate containerized applications

### Individual Application Features
**Log Output:**
- Background UUID logging with configurable intervals
- Thread-safe storage with RWMutex protection

**Ping Pong:**
- Simple counter-based responses
- Memory-based state management

## 🐳 Docker Commands

### Build Images
```bash
# Build log_output image
cd log_output
docker build -t <username>/log_output-img:ex1.09 .

# Build ping_pong image
cd ../ping_pong
docker build -t <username>/ping_pong-img:ex1.09 .
```

### Push to Docker Hub
```bash
docker push <username>/log_output-img:ex1.09
docker push <username>/ping_pong-img:ex1.09
```

### Run Locally with Docker
```bash
# Run log_output locally
docker run -p 8091:8091 -e PORT=8091 <username>/log_output-img:ex1.09

# Run ping_pong locally (in separate terminal)
docker run -p 8092:8092 -e PORT=8092 <username>/ping_pong-img:ex1.09
```

## ☸️ Kubernetes Deployment

### Create k3d Cluster
```bash
k3d cluster create --port 8081:80@loadbalancer --agents 2
```

### Configure kubectl Context
```bash
kubectl config get-contexts
kubectl config use-context k3d-k3s-default
```

### Deploy Applications
```bash
# Deploy both applications and shared ingress
kubectl apply -f log_output/manifests/
kubectl apply -f ping_pong/manifests/
kubectl apply -f shared/manifests/
```

### Verify Deployment
```bash
# Check all deployments
kubectl get deployments

# Check all pods
kubectl get pods

# Check all services and ingress
kubectl get svc,ing

# Check logs from both applications
kubectl logs -f deployment/log-output-dep
kubectl logs -f deployment/ping-pong-dep
```

### Access Applications

**Primary Method: Shared Ingress**
```bash
# Access log_output endpoints
curl localhost:8081/logs      # Get all logs
curl localhost:8081/status    # Get status with last 10 logs
curl localhost:8081/status?n=5 # Get status with last 5 logs

# Access ping_pong (pingpong path)
curl localhost:8081/pingpong

```

**Alternative Method: Port Forward**
```bash
# Port forward to log_output service
kubectl port-forward service/log-output-svc 8091:2345

# Port forward to ping_pong service (separate terminal)
kubectl port-forward service/ping-pong-svc 8092:1234

# Test applications
curl localhost:8091/logs
curl localhost:8091/status
curl localhost:8092/pingpong
```

## 📊 Expected Responses

### Log Output Application

**GET /logs** - All stored log entries:
```json
{
  "logs": [
    {
      "timestamp": "2024-09-01T12:00:00Z",
      "value": "550e8400-e29b-41d4-a716-446655440000"
    },
    {
      "timestamp": "2024-09-01T12:00:05Z",
      "value": "550e8400-e29b-41d4-a716-446655440001"
    }
  ]
}
```

**GET /status** - Status with last 10 log entries (default):
```json
{
  "status": "ready",
  "logs": [
    {
      "timestamp": "2024-09-01T12:00:15Z",
      "value": "550e8400-e29b-41d4-a716-446655440003"
    }
  ]
}
```

**GET /status?n=5** - Status with last 5 log entries:
```json
{
  "status": "ready",
    "logs": [
    {
      "timestamp": "2025-09-01T05:35:00.992022875Z",
      "value": "7784aa04-6cc7-44e7-9ca4-af4172ee9a9f"
    },
    {
      "timestamp": "2025-09-01T05:35:05.99636521Z",
      "value": "7784aa04-6cc7-44e7-9ca4-af4172ee9a9f"
    },
    {
      "timestamp": "2025-09-01T05:35:10.996151948Z",
      "value": "7784aa04-6cc7-44e7-9ca4-af4172ee9a9f"
    },
    {
      "timestamp": "2025-09-01T05:35:15.992089615Z",
      "value": "7784aa04-6cc7-44e7-9ca4-af4172ee9a9f"
    },
    {
      "timestamp": "2025-09-01T05:35:20.99616667Z",
      "value": "7784aa04-6cc7-44e7-9ca4-af4172ee9a9f"
    }
  ]
}
```

### Ping Pong Application

**GET /pingpong** - First request:
```json
{
  "count": 1
}
```

**Subsequent Ping Pong Requests:**
```json
{
  "count": 2
}
```
```json
{
  "count": 3
}
```

## 🛠️ Local Development

### Log Output Application
```bash
cd log_output
make build       # Build the application
make run         # Run locally on port 8091
make watch       # Live reload with Air
```

### Ping Pong Application
```bash
cd ping_pong
make build       # Build the application
make run         # Run locally on port 8092
make watch       # Live reload with Air
```

## ⚙️ Configuration

### Environment Variables
- **Log Output:**
  - `PORT`: Server port (default: 8091)
- **Ping Pong:**
  - `PORT`: Server port (default: 8092)

### Kubernetes Resources

#### Services
- **log-output-svc:** ClusterIP service (port 2345 → 8091)
- **ping-pong-svc:** ClusterIP service (port 1234 → 8092)

#### Deployments
- **log-output-dep:** Manages log_output application pods
- **ping-pong-dep:** Manages ping_pong application pods

#### Shared Ingress
- **Routes `/logs` → log-output-svc:2345**
- **Routes `/status` → log-output-svc:2345**
- **Routes `/pingpong` → ping-pong-svc:1234**
- **External access via port 80 (mapped to 8081 locally)**

### Network Architecture
```
Internet → k3d LoadBalancer:8081 → Shared Ingress:80
                                      ├─ /logs → log-output-svc:2345 → log_output:8091
                                      ├─ /status → log-output-svc:2345 → log_output:8091
                                      └─ /pingpong → ping-pong-svc:1234 → ping_pong:8092
```



kubectl test through busybox using dns service with Service name
```bash
kubectl exec -it my-busybox -- wget -qO - http://pingpong-svc:2366/pingpong | jq .
```


kubectl test through busybox using ClusterIP
```bash
kubectl exec -it my-busybox -- wget -qO - http://10.43.122.68:2366/pingpong | jq .
```

kubectl test through busybox Direct pod IP access
```bash
 kubectl exec -it my-busybox -- wget -qO - http://10.42.1.11:8096/pingpong | jq .
 kubectl exec -it my-busybox -- wget -qO - http://10.42.1.11:8095/
```


