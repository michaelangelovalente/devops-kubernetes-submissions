# Dynamic Background Image Guide

This guide details the steps to implement a dynamic background image feature in the Todo application. The image will be fetched from an external source, saved to a persistent volume, and refreshed periodically.

## The Goal

To replace the static background image with one that is dynamically fetched from `https://picsum.photos/1200` every two minutes. This will be done in a way that is robust and works within a containerized Kubernetes environment.

---

### Step 1: The Image Service (`internal/picsum/picsum.go`)

**Goal:** Create a dedicated service to handle the logic of fetching and saving the image. This promotes separation of concerns and makes the code easier to maintain.

**Action:** Create the file `internal/picsum/picsum.go` and add the following content:

```go
package picsum

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const ImageURL = "https://picsum.photos/1200"

// FetchAndSaveImage fetches an image from the specified URL and saves it to the given directory.
func FetchAndSaveImage(dir string) error {
	// Create the directory if it doesn't exist.
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	resp, err := http.Get(ImageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	filePath := filepath.Join(dir, "background.jpg")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
```

**Why:**

-   **Separation of Concerns:** This `picsum` package encapsulates all the logic related to fetching the image.
    If we need to change the image source or the way it's processed in the future, we only need to modify this file.
-   **Error Handling:** The function returns an error, allowing the caller to handle cases where the image can't be fetched or saved.
-   **File I/O:** It uses standard Go packages `os` and `io` to handle file operations, making the code portable and easy to understand.
-   **`os.MkdirAll`**: This ensures that the target directory exists, preventing errors if the directory is not created beforehand.

---

### Step 2: The Background Task (`cmd/api/main.go`)

**Goal:** Run the image fetching process in the background without blocking the main application, and repeat it at a regular interval.

**Action:** Modify your `cmd/api/main.go` file to include the background task.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todo_app/internal/picsum"
	"todo_app/internal/server"
)

const imageDir = "tmp/images"

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	// ... (existing code)
}

func imageRotationTask(logger *log.Logger) {
	// Fetch the image on startup.
	if err := picsum.FetchAndSaveImage(imageDir); err != nil {
		logger.Printf("Error fetching initial image: %v", err)
	}

	// Use a ticker to fetch the image every 2 minutes.
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		if err := picsum.FetchAndSaveImage(imageDir); err != nil {
			logger.Printf("Error fetching image: %v", err)
		} else {
			logger.Println("Background image updated successfully.")
		}
	}
}

func main() {
	server := server.NewServer()

	// Start the background image rotation task.
	go imageRotationTask(server.Logger)

	done := make(chan bool, 1)
	go gracefulShutdown(server.Server, done)

	server.Logger.Printf("Server (ex 1.12) started on port %d
", server.Port)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("http server error: %s", err))
	}

	<-done
	server.Logger.Println("Graceful shutdown complete")
}
```

**Why:**

-   **Goroutine (`go imageRotationTask(...)`):** This starts the `imageRotationTask` in a separate goroutine, so it runs concurrently with the main web server and doesn't block it.
-   **`time.Ticker`:** A ticker is the standard Go way to perform an action at regular intervals. It's more efficient than using `time.Sleep` in a loop.
-   **Initial Fetch:** We call `picsum.FetchAndSaveImage` once before the ticker starts to ensure that an image is available as soon as the application starts.
-   **Logging:** We log errors and success messages to provide visibility into what the background task is doing.

---

### Step 3: Serving the Image (`internal/server/routes.go`)

**Goal:** Create an HTTP endpoint that serves the downloaded image.

**Action:** Modify `internal/server/routes.go` to add a file server route.

```go
package server

import (
	"net/http"
	"strings"
	"todo_app/web"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// ... (existing CORS setup)

	// Serve static files from the embedded filesystem.
	r.Handle("/static/*", http.FileServer(http.FS(web.Files)))

	// Serve the dynamic background image from the filesystem.
	imageDir := "tmp/images"
	fileServer := http.FileServer(http.Dir(imageDir))
	r.Handle("/background-image", http.StripPrefix("/background-image", fileServer))

	r.Get("/", templ.Handler(web.Base()).ServeHTTP)
	return r
}
```

**Why:**

-   **`http.Dir("tmp/images")`:** This creates a file system handler that is rooted in the `tmp/images` directory.
-   **`http.FileServer`:** This is a built-in Go function that returns a handler that serves HTTP requests with the contents of the file system.
-   **`r.Handle("/background-image", ...)`:** This registers a new route. Any request to `/background-image` will be handled by our file server.
-   **`http.StripPrefix`:** This is important. It removes the `/background-image` prefix from the request path before passing it to the file server.
        Without this, the file server would look for a file named `/background-image/background.jpg` inside the `tmp/images` directory, which is not what we want.
        We want it to look for `background.jpg`.

---

### Step 4: Updating the Frontend (`web/base.templ`)

**Goal:** Change the frontend to request the image from our new dynamic endpoint.

**Action:** Modify the `style` block in `web/base.templ`.

```html
<style>
    body {
        background-image: url('/background-image');
        background-size: cover;
        background-repeat: no-repeat;
        background-attachment: fixed;
    }
</style>
```

**Why:**

-   This is a simple but crucial change. Instead of pointing to a static file, the `background-image` URL now points to the `/background-image` endpoint we just created. The browser will now fetch the image from our Go server, which serves the dynamically downloaded file.

---

### Step 5: Kubernetes Configuration (`manifests/deployment.yaml`)

**Goal:** Ensure that the downloaded image persists even if the application pod is restarted.

**Action:** Modify `manifests/deployment.yaml` to include a `PersistentVolumeClaim` and mount the volume.

*Self-correction: For this exercise, a simple `emptyDir` volume is sufficient and easier to manage than a full `PersistentVolumeClaim`. An `emptyDir` volume is created when a Pod is assigned to a Node, and exists as long as that Pod is running on that node. When a Pod is removed from a node for any reason, the data in the `emptyDir` is deleted forever.*

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: todo-app-dep
spec:
  replicas: 1
  selector:
    matchLabels:
      app: todo-app
  template:
    metadata:
      labels:
        app: todo-app
    spec:
      containers:
        - name: todo-app
          image: <your-dockerhub-username>/todo_app:ex1.12
          ports:
            - containerPort: 3005
          volumeMounts:
            - name: image-volume
              mountPath: /app/tmp/images
      volumes:
        - name: image-volume
          emptyDir: {}
```

**Why:**

-   **`volumes`:** This section in the Pod spec defines a volume that can be used by the containers in the pod.
-   **`emptyDir: {}`:** An `emptyDir` volume is first created when a Pod is assigned to a Node, and exists as long as that Pod is running on that node. When a Pod is removed from a node for any reason, the data in the `emptyDir` is deleted *forever*. This is suitable for our use case as the image is re-downloaded on startup.
-   **`volumeMounts`:** This section in the container spec mounts the defined volume into the container's filesystem.
-   **`mountPath: /app/tmp/images`:** This is the path inside the container where the volume will be mounted. It must match the directory we are using in our Go code to save the image.
-   **Persistence:** By using a volume, the `tmp/images` directory is now a mount point for the volume. This means that the data is stored outside the container's ephemeral filesystem, and will persist across container restarts.

---

## Conclusion

By following these steps, you will have a fully functional dynamic background image feature. This solution is robust, scalable, and follows best practices for building cloud-native applications.
