# Todo App - Exercise 1.2

A Go-based HTTP server application for the DevOps Kubernetes Course.

## 🔗 Links

- **Docker Hub:** [michaelangelovalente/todo_app:ex1.02](https://hub.docker.com/layers/michaelangelovalente/todo_app/ex1.02/images/sha256-a9e0a9948e8ec63389aecf5ecbdcdc2ceee05f86e52115665d7f858b550e50d6)
- **Source Code:** [GitHub Repository](https://github.com/michaelangelovalente/devops-kubernetes-submissions/tree/main/Chapter-2/Part-1/e-1.02/todo_app)

## 📋 Description

The `todo_app` is a Go-based HTTP server that meets exercise 1.2 requirements by logging `Server started in port NNNN` on startup. The port number is configurable via environment variable, defaulting to **8080** if not specified.

## 🐳 Docker Commands

### Build Image
```bash
docker build -t <username>/todo_app:ex1.02 .
```

**Alternative using Docker Compose:**
```bash
docker compose up
docker image tag todo_app:latest <username>/todo_app:ex1.02
```

### Push to Docker Hub
```bash
docker push <username>/todo_app:ex1.02
```

## ☸️ Kubernetes Deployment

### Configure kubectl Context
```bash
kubectl config get-contexts
kubectl config use-context todo-app-cluster  # if not already set
```

### Deploy Application
```bash
kubectl create deployment todo-app-ex-1-2 --image=michaelangelovalente/todo_app:ex1.02
```

### Monitor Deployment
```bash
kubectl get deployments
kubectl get replicasets
kubectl get pods
kubectl logs -f <pod-name>
```

**Pod Naming Convention:**
```
todo-app-ex-1-2                           # Deployment name
todo-app-ex-1-2-85cbd45c46                # + ReplicaSet hash
todo-app-ex-1-2-85cbd45c46-74tw9          # + Pod suffix
```

## 🛠️ Local Development

### Using Makefile Commands
```bash
make build       # Build the application
make run         # Run locally
make docker-run  # Start with Docker Compose
make docker-down # Stop Docker containers
make watch       # Live reload during development
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
- `PORT`: Server port (default: 8080)
