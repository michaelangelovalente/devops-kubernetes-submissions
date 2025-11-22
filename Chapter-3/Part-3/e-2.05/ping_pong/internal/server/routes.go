package server

import (
	"net/http"

	"ping_pong/internal/app"

	common_server "common/server"
)

func RegisterRoutes(app *app.Application) http.Handler {
	r := common_server.NewRouter()

	r.Get("/pingpong", app.PingpongHandler.Update)

	return r
}
