# 1.2 Project todo_app

## Description

**Docker Hub repository:** [todo_app:ex1.02](https://hub.docker.com/layers/michaelangelovalente/todo_app/ex1.02/images/sha256-eefdc7041060a4df2aa37c456aee658ffaf7d461f83368d2c1407c6a66258709)
**Source Code:** [todo_app source code]()

`todo_app` application is a Go-based HTTP server.
As per exercise requirement (1.2), it logs `Server started in port NNNN`, where NNNN is a environment set port number, if this value isn't set
the application will use port 8080 by default.

## Commands

### Docker Image Management

#### Build Docker Image
```bash
docker build -t <username>/todo_app:ex1.02 .
```

``` bash
# alternative:
docker compose up
docker image tag <username>/todo_app:ex1.02
```

#### Push to Docker Hub
```bash
docker push <username>/todo_app:ex1.02
```

### Kubernetes Deployment


#### Verify / Configure kubectl Context
```bash
kubectl config get-contexts
# if context not first-deploy-cluster
kubectl config use-context log-output-cluster
```

#### Deploy Application
```bash
kubectl create deployment todo-app-ex-1-2 --image=michaelangelovalente/log-output:ex1.02

```

#### Monitor Deployment
```bash
kubectl get deployments
kubectl get replicasets
kubectl get pods
kubectl logs -f <deploymentname-replacesethash-generatedpodsuffix>
```

```
pod name composition notes:
 todo-app-ex-1-2 (deployment name)
 [todo-app-ex-1-2]-85cbd45c46 (deployment name + ReplicaSet hash)
 [[todo-app-ex-1-2]-85cbd45c46]-74tw9  (dployment name + ReplicaSet hash + generated pod suffix)
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
