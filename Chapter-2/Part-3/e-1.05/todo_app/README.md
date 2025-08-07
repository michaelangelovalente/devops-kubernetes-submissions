# Todo App - Exercise 1.5

A Go-based + HTMX application for the DevOps Kubernetes Course exercise 1.05


## 🔗 Links

- **Docker Hub:** [michaelangelovalente/todo_app:ex1.05]()
- **Source Code:** [GitHub Repository]()

## 📋 Description

The `todo_app` is a Go-based HTTP server that meets exercise 1.5 requirements by
starting a backend service with a base end points that redirects to a todo app's landing page
## 🐳 Docker Commands

```bash
docker build -t <username>/todo_app:ex1.05 .
```

### Push to Docker Hub
```bash
docker push <username>/todo_app:ex1.05
```

## ☸️ Kubernetes Deployment

### Configure kubectl Context
```bash
kubectl config get-contexts
# using k3d-k3s-default
```

### Deploy Application
```bash
# creating deployment
kubectl apply -f manifests/deployment.yaml
```

### Check Environment variables APP_PORT
```bash
kubectl exec todo-app-dep-5884c7f66d-ffwr5 -- printenv
```

### Port forward and verify
```bash
kubectl port-forward service/todo-app-svc 8090:8080
```

#### Expected Output
```bash
curl localhost:8090/
```

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
- `APP_PORT`: Server port (default: 8080) / kubernetes manifest deployment.yaml ( port 3005)
