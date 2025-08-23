# Log Output - Exercise 1.03

A Go-based HTTP server with background logging for the DevOps Kubernetes Course.

## 🔗 Links

- **Docker Hub:** [michaelangelovalente/log-output:ex1.03](https://hub.docker.com/layers/michaelangelovalente/log_output-app/ex1.03/images/sha256-06245083d1e3c81405f029e24e29a4cc80cdfe03ecb023f1aaefce76d6832ba8)
- **Source Code:** [GitHub Repository](https://github.com/michaelangelovalente/devops-kubernetes-submissions/tree/main/Chapter-2/Part-1/e-1.03/log_output)

## 📋 Description

The `log_output` application is a Go-based HTTP server that generates timestamped log entries every 5 seconds.
Each entry contains a random UUID string stored in memory and displayed with timestamps.
This version uses `./manifests/deployment.yaml` for a declarative release approach

Built with **Go 1.24+** following clean architecture patterns, this application demonstrates:
- Background goroutines with graceful shutdown using context cancellation
- In-memory storage with interface-based design
- CORS-enabled HTTP server with configurable timeouts
- Structured logging and UUID generation with 5-second intervals


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
kubectl config use-context k3d-first-deploy-cluster
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
- Background UUID logging every 5 seconds with timestamps
- Declarative Kubernetes deployment through manifests/

