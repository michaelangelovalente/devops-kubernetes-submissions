# Todo App - Exercise 1.4

A Go-based HTTP server application for the DevOps Kubernetes Course.

## 🔗 Links

- **Docker Hub:** [michaelangelovalente/todo_app:ex1.04](https://hub.docker.com/layers/michaelangelovalente/todo_app/ex1.04/images/sha256-e4bc8296f2beab45586673010335e9e2f1773001bfe5bfa42cb6ed8c379f4861)
- **Source Code:** [GitHub Repository](https://github.com/michaelangelovalente/devops-kubernetes-submissions/tree/main/Chapter-2/Part-1/e-1.04/todo_app)

## 📋 Description

The `todo_app` is a Go-based HTTP server that meets exercise 1.4 requirements by logging `Server started in port NNNN` on startup.
The port number is configurable via environment variable, defaulting to **8080** if not specified.
This version of `todo_app` is identical to the version 1.02, the only difference is that it allows Declarative Deployment through `manifests/deployment.yaml` (exercise 1.04 requirement)

## 🐳 Docker Commands

### Build Image
```bash
docker build -t <username>/todo_app:ex1.04 .
```

**Alternative using Docker Compose:**
```bash
docker compose up
docker image tag todo_app:latest <username>/todo_app:ex1.04
```

### Push to Docker Hub
```bash
docker push <username>/todo_app:ex1.04
```

## ☸️ Kubernetes Deployment

### Deploy Application
```bash
kubectl apply -f manifests/
```

### Monitor Deployment
```bash
kubectl get deployments
kubectl get replicasets
kubectl get pods
kubectl logs -f <pod-name>
```

---

**Environment Variables:**
- `PORT`: Server port (default: 8080)
