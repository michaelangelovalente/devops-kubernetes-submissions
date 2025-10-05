package server

import (
	"net/http"

	"todo_app/internal/todo"
	"todo_app/web"
	"todo_app/web/views"

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

	r.Handle("/static/*", http.FileServer(http.FS(web.Files)))

	r.Get("/image", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s.ImageDir+"/background.jpg")
	})

	r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
		task := r.FormValue("task")
		newTodo := todo.AddTodo(task)
		templ.Handler(views.Todo(newTodo)).ServeHTTP(w, r)
	})

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(web.Base(todo.GetTodos())).ServeHTTP(w, r)
	}))

	return r
}
