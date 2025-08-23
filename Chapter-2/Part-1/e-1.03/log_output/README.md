# Log Output - Exercise 1.1

A Go-based HTTP server with background logging for the DevOps Kubernetes Course.

## 🔗 Links

- **Docker Hub:** [michaelangelovalente/log-output:ex1.03]()
- **Source Code:** [GitHub Repository]()

## 📋 Description

The `log_output` application is a Go-based HTTP server that generates timestamped log entries every 5 seconds.
Each entry contains a random UUID string stored in memory and displayed with timestamps.
This version uses `./manifests/deployment.yaml` for a declarative release approach

Built with modern **Go 1.24+** practices following clean architecture patterns, this application demonstrates:
- Background goroutines with graceful shutdown
- In-memory storage with interfaces
- CORS-enabled HTTP server
- Structured logging and UUID generation


## 🐳 Docker Commands

### Build Image
```bash
docker build -t <username>/log-output:ex1.03 .
```

### Push to Docker Hub
```bash
docker push <username>/log-output:ex1.03
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

### Declarative Deployment
```bash
kubectl apply -f manifests/
```

### Monitor Deployment
```bash
kubectl get deployments
kubectl get pods
kubectl logs -f <pod-name>
```

---

**Environment Variables:**
- `PORT`: Server port (default: 3000)

**Features:**
- Background UUID logging every 5 seconds
- CORS-enabled HTTP endpoints
- Graceful shutdown with context cancellation
- Declarative deployment through manifests/deployment.yaml

