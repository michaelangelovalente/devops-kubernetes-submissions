# Todo App - Exercise 1.06

A Todo application built with Go, HTMX, and TailwindCSS for the DevOps with Kubernetes Course.

## 🔗 Links

- **Docker Hub:** [michaelangelovalente/todo_app:ex1.06](https://hub.docker.com/layers/michaelangelovalente/todo_app/ex1.06/images/sha256-29e7c9f66809b52e63c4a9aa31c10aa6982551a661c2bda4c948092073a0ab28)
- **Source Code:** [GitHub Repository](https://github.com/michaelangelovalente/devops-kubernetes-submissions/tree/main/Chapter-2/Part-3/e-1.06/todo_app)

## 📋 Description

- **Backend:** Go HTTP server with structured logging and graceful shutdown
- **Frontend:** HTMX for dynamic interactions without JavaScript complexity
- **Styling:** TailwindCSS
- **Infrastructure:** Docker containerization and Kubernetes deployment with NodePort Service

This application demonstrates how to use a NodePort Service to enable external access to pods running in a Kubernetes cluster, fulfilling exercise 1.06 requirements.
## 🐳 Docker Commands

### Build Image
```bash
docker build -t <username>/todo_app:ex1.06 .
```

### Push to Docker Hub
```bash
docker push <username>/todo_app:ex1.06
```

### Run Locally with Docker
```bash
docker run -p 8080:8080 <username>/todo_app:ex1.06
```

## ☸️ Kubernetes Deployment

### Create k3d Cluster
```bash
k3d cluster create --port 8089:30080@agent:0 -p 8081:2345@loadbalancer --agents 3
```

### Configure kubectl Context
```bash
kubectl config get-contexts
kubectl config use-context k3d-k3s-default
```

### Deploy Application
```bash
# Apply both deployment and service manifests
kubectl apply -f manifests/
```

### Verify Deployment
```bash
# Check deployment status
kubectl get deployments

# Check pods
kubectl get pods

# Check services (including NodePort)
kubectl get services

# Check environment variables
kubectl exec <pod-name> -- printenv | grep APP_PORT
```

### Access Application

**Method 1: NodePort Service (Primary for Exercise 1.06)**
```bash
# Access via NodePort (assuming k3d cluster with port mapping)
# First, check the NodePort assigned:
kubectl get service todo-app-svc

# If using k3d with port mapping 8089:30080@agent:0
# and service uses NodePort 30080:
curl localhost:8089/

# Or access directly via cluster node IP:
# Get node IP:
kubectl get nodes -o wide
# Access: http://<NODE_IP>:30080/
```

**Method 2: Port Forward (Alternative)**
```bash
# Port forward to access the application (development/debugging)
kubectl port-forward service/todo-app-svc 8090:8080

# Test the application
curl localhost:8090/
```

#### Expected Output

```html
<!doctype html>
<html lang="en" class="h-screen">

<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width,initial-scale=1">
    <title>TODO App - Exercise 1.05</title>
    <script src="https://unpkg.com/htmx.org/dist/htmx.min.js"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</head>

<body class="bg-gray-50 min-h-screen">
    <div class="container mx-auto max-w-md p-6">
        <header class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-800 mb-2">TODO App</h1>
            <p class="text-gray-600">Exercise 1.05</p>
        </header>
        <main class="bg-white rounded-lg shadow-md p-6">
            <form class="mb-6" hx-post="/todos" hx-target="#todo-list" hx-swap="beforeend">
                <div class="flex gap-2"><input type="text" name="task" placeholder="Add a new todo..."
                        class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                        required> <button type="submit"
                        class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500">Add</button>
                </div>
            </form>
            <div id="todo-list" class="space-y-2"><!-- Todo items will be inserted here --></div>
        </main>
    </div>
</body>

</html>
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
- **Live Reload:** Use `make watch` for automatic rebuilding during development
- **HTMX Integration:** Dynamic todo interactions without page refreshes
- **TailwindCSS:** Responsive design with modern styling
- **Graceful Shutdown:** Clean server termination handling

---

## ⚙️ Configuration

**Environment Variables:**
- `APP_PORT`: Server port (default: 8080, Kubernetes deployment uses 3005)

**Service Configuration:**
- **Service Type:** NodePort (enables external access without LoadBalancer)
- **Target Port:** 8080 (application port)
- **Service Port:** 8080 (cluster-internal port)
- **NodePort:** 30080 (external access port on cluster nodes)

