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

	r.Handle("/static/*", http.FileServer(http.FS(web.Files)))

	r.Get("/image", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s.ImageDir+"/background.jpg")
	})

	r.Get("/", templ.Handler(web.Base()).ServeHTTP)
	return r
}
