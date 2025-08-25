package routes

import (
	"net/http"

	"log_output/internal/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func RegisterRoutes(app *app.Application) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// r.Get("/", HelloWorldHandler)
	r.Get("/logs", app.LogMemoryHandler.GetAllLogs)
	r.Get("/status", app.LogMemoryHandler.GetLastLogsAndStatus)

	return r
}

// Basic health and ready handlers
// func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
// 	// simple health check
// 	utils.WriteJson(w, http.StatusOK, utils.Envelope{"status": "ready"})
// }

// func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
// 	resp := make(map[string]string)
// 	resp["message"] = "Hello World"
//
// 	jsonResp, err := json.Marshal(resp)
// 	if err != nil {
// 		log.Fatalf("error handling JSON marshal. Err: %v", err)
// 	}
//
// 	_, _ = w.Write(jsonResp)
// }
