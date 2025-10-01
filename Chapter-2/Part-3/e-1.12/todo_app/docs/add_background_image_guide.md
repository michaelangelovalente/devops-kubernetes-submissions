# How to Add a Background Image to Your Go Web Application

This guide will walk you through the process of adding a background image to your Go web application.

## 1. Project Structure

For a well-organized project, it's best practice to store static assets like images, CSS, and JavaScript files in a dedicated directory. We'll create a `static` directory inside the `web` directory.

Your project structure should look like this:

```
todo_app/
├── web/
│   ├── static/
│   │   └── example.png
│   ├── base.templ
│   └── views/
...
```

You can create this directory and move your image using the following commands:

```bash
mkdir -p todo_app/web/static
mv todo_app/tmp/example.png todo_app/web/static/example.png
```

## 2. Serving Static Files

Next, you need to configure your Go server to serve files from the `web/static` directory. We'll use a `FileServer` for this.

In `internal/server/routes.go`, modify the `RegisterRoutes` function to include a file server for the `/static/` path.

```go
package server

import (
	"net/http"
	"todo_app/web"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Serve static files
	fs := http.FileServer(http.Dir("web/static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", templ.Handler(web.Base()).ServeHTTP)
	return r
}
```

**Explanation:**

*   `http.Dir("web/static")`: This creates a file system handler for the `web/static` directory.
*   `http.FileServer(...)`: This creates a handler that serves HTTP requests with the contents of the file system.
*   `r.Handle("/static/*", ...)`: This tells the router to handle all requests that start with `/static/`.
*   `http.StripPrefix("/static/", fs)`: This removes the `/static/` prefix from the request path before passing it to the file server. This is important so that the file server can find the correct file. For example, a request to `/static/example.png` will be translated to a request for `example.png` in the `web/static` directory.

## 3. Updating the HTML

Now that the server is configured to serve the image, you can reference it in your HTML. We'll add a `style` tag to the `web/base.templ` file to set the background image.

```html
package web

templ Base() {
	<html>
		<head>
			<title>Todo App</title>
			<style>
				body {
					background-image: url('/static/example.png');
					background-size: cover;
					background-repeat: no-repeat;
					background-attachment: fixed;
				}
			</style>
		</head>
		<body>
			<h1>Todo App</h1>
		</body>
	</html>
}
```

**Explanation:**

*   `background-image: url('/static/example.png');`: This CSS rule sets the background image of the `body` element to the image served at `/static/example.png`.
*   `background-size: cover;`: This scales the image to cover the entire background of the element.
*   `background-repeat: no-repeat;`: This prevents the image from repeating.
*   `background-attachment: fixed;`: This fixes the background image relative to the viewport, so it doesn't scroll with the content.

## 4. Run the Application

Now you can run your application and you should see the background image.

```bash
go run cmd/api/main.go
```

Open your browser and navigate to `http://localhost:3005`. You should see the background image.

## 5. Alternative: Embedding Static Files

Instead of serving files from the disk, you can embed them directly into your Go binary. This is a great option for production builds as it creates a single, self-contained executable.

### 5.1. Update `routes.go`

Modify your `internal/server/routes.go` file to use the `embed` package.

```go
package server

import (
	"embed"
	"net/http"
	"todo_app/web"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//go:embed web/static
var staticFiles embed.FS

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Serve static files from the embedded file system
	fileServer := http.FileServer(http.FS(staticFiles))
	r.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	r.Get("/", templ.Handler(web.Base()).ServeHTTP)
	return r
}
```

**Explanation:**

*   `import "embed"`: This imports the `embed` package.
*   `//go:embed web/static`: This is a compiler directive that tells Go to embed the `web/static` directory into the `staticFiles` variable.
*   `var staticFiles embed.FS`: This declares a variable of type `embed.FS` to hold the embedded files.
*   `http.FileServer(http.FS(staticFiles))`: This creates a file server that serves files from the embedded file system.

### 5.2. Comparison

| Method | Pros | Cons |
| --- | --- | --- |
| `http.Dir()` | - Great for development (no need to recompile for static file changes) | - Requires separate deployment of static files |
| `http.FS()` with `embed` | - Creates a single, self-contained binary for deployment | - Requires recompilation for any changes to static files |