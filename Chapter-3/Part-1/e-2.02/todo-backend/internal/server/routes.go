package server

import (
	"net/http"

	common_server "common/server"
)

func (s *AppServer) RegisterRoutes() http.Handler {
	r := common_server.NewRouter()
	r.Post("/todo", s.todoHandler.AddTodo)
	r.Get("/todos", s.todoHandler.GetTodos)
	return r
}
