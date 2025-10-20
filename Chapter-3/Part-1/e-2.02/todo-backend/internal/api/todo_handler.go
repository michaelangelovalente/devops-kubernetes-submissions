package api

import (
	"common/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"todo-backend/internal/store"
)

type TodoHandler struct {
	logger *log.Logger
}

func NewTodoHandler(logger *log.Logger) *TodoHandler {
	return &TodoHandler{
		logger: logger,
	}
}

func (th *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	todos := store.GetTodos()
	utils.WriteJSON(w, http.StatusOK,
		utils.Envelope{
			"data": todos,
		},
	)
}

func (th *TodoHandler) AddTodo(w http.ResponseWriter, r *http.Request) {
	var todoEntry struct {
		Data string `json:"data"`
	}
	err := json.NewDecoder(r.Body).Decode(&todoEntry)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest,
			utils.Envelope{
				"error": fmt.Sprintf("could not decode request %v", err),
			},
		)
		return
	}
	_ = store.AddTodo(todoEntry.Data)

	todos := store.GetTodos()
	utils.WriteJSON(w, http.StatusOK,
		utils.Envelope{
			"data": todos,
		},
	)

}

// func (th *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
// 	todos := store.GetTodos()
// 	utils.WriteJSON(w, http.StatusOK,
// 		utils.Envelope{
// 			"data": todos,
// 		},
// 	)
// }
