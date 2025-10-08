package server

import (
	"net/http"

	"log_output/internal/app"

	common_server "common/server"
)

func RegisterRoutes(app *app.Application) http.Handler {
	r := common_server.NewRouter()
	r.Get("/logs", app.LogMemoryHandler.GetAllLogs)
	r.Get("/status", app.LogMemoryHandler.GetLastLogsAndStatus)

	return r
}
