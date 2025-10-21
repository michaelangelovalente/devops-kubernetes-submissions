package server

import (
	"net/http"

	"todo_app/web"
	"todo_app/web/views"

	common_server "common/server"

	"github.com/a-h/templ"
)

func (s *AppServer) RegisterRoutes() http.Handler {
	r := common_server.NewRouter()
	r.Handle("/static/*", http.FileServer(http.FS(web.Files)))

	r.Get("/image", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, s.ImageDir+"/background.jpg")
	})

	r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
		task := r.FormValue("task")
		newTodo, _ := s.TodoClient.AddTodo(task)
		templ.Handler(views.Todo(newTodo)).ServeHTTP(w, r)
	})

	r.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		todos, _ := s.TodoClient.GetTodos()

		templ.Handler(web.Base(todos)).ServeHTTP(w, r)
	}))

	return r
}
