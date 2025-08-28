package routes

import (
	"net/http"

	"ping_pong/internal/app"

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
	r.Get("/test", app.PingpongHandler.Get)

	return r
}

// func (a *Application) setupRoutes() {
//     r := a.router
//
//     // Global middleware
//     r.Use(middleware.RequestID)
//     r.Use(middleware.RealIP)
//     r.Use(middleware.Logger)
//     r.Use(middleware.Recoverer)
//
//     // Add context timeout middleware for all requests
//     r.Use(func(next http.Handler) http.Handler {
//         return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//             // Create a timeout context for each request
//             ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
//             defer cancel()
//
//             next.ServeHTTP(w, r.WithContext(ctx))
//         })
//     })
//
//     // Health check endpoints (no timeout needed)
//     r.Get("/health/live", a.handlers.Log.HealthLive)
//     r.Get("/health/ready", a.handlers.Log.HealthReady)
//
//     // API routes
//     r.Route("/api/v1", func(r chi.Router) {
//         r.Get("/logs", a.handlers.Log.GetCurrentLog)
//         r.Get("/pingpong", a.handlers.Ping.HandlePing)
//     })
// }
