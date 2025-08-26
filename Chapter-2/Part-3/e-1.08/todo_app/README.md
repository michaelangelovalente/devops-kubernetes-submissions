# Todo App - Exercise 1.08

A Todo application built with Go, HTMX, and TailwindCSS for the DevOps with Kubernetes Course.

## 🔗 Links

- **Docker Hub:** [michaelangelovalente/todo_app:ex1.08](https://hub.docker.com/r/michaelangelovalente/todo_app)
- **Source Code:** [GitHub Repository](https://github.com/michaelangelovalente/devops-kubernetes-submissions/tree/main/Chapter-2/Part-3/e-1.08/todo_app)

## 📋 Description

- **Backend:** Go HTTP server with structured logging and graceful shutdown
- **Frontend:** HTMX for dynamic interactions without JavaScript complexity
- **Styling:** TailwindCSS
- **Infrastructure:** Docker containerization and Kubernetes deployment with Ingress

This application demonstrates how to use Kubernetes Ingress to enable external HTTP access to pods running in a cluster, fulfilling exercise 1.08 requirements.
## 🐳 Docker Commands

### Build Image
```bash
docker build -t <username>/todo_app:ex1.08 .
```

### Push to Docker Hub
```bash
docker push <username>/todo_app:ex1.08
```

### Run Locally with Docker
```bash
docker run -p 3005:3005 <username>/todo_app:ex1.08
```

## ☸️ Kubernetes Deployment

### Create k3d Cluster
```bash
k3d cluster create --port 8090:30090@agent:0 -p 8082:80@loadbalancer --agents 2
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

# Check services and ingress
kubectl get svc,ing

# Check environment variables
kubectl exec <pod-name> -- printenv | grep APP_PORT
```

### Access Application

**Method 1: Ingress (Primary for Exercise 1.08)**
```bash
# Access via Ingress through load balancer
# If using k3d with port mapping 8082:80@loadbalancer
curl localhost:8082/

# Or access directly in browser:
# http://localhost:8082/
```

**Method 2: Port Forward (Alternative)**
```bash
# Port forward to access the application (development/debugging)
kubectl port-forward service/todo-app-svc 8090:1234

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
    <title>TODO App</title>
    <script src="https://unpkg.com/htmx.org/dist/htmx.min.js"></script>
    <script src="https://cdn.tailwindcss.com"></script>
</head>

<body class="bg-gray-50 min-h-screen">
    <div class="container mx-auto max-w-md p-6">
        <header class="text-center mb-8">
            <h1 class="text-3xl font-bold text-gray-800 mb-2">TODO App</h1>
            <p class="text-gray-600">Exercise 1.08</p>
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

**Kubernetes Resources:**
- **Deployment:** Manages application pods and replica sets
- **Service Type:** ClusterIP (internal cluster communication)
- **Service Port:** 1234 (cluster-internal port)
- **Target Port:** 3005 (application port)
- **Ingress:** Provides external HTTP access on port 80

**Network Configuration:**
- **External Access:** Available via browser through Ingress controller
- **Internal Access:** Available to other cluster services via ClusterIP service on port 1234
- **Application Port:** Listens on port 3005 inside the container

