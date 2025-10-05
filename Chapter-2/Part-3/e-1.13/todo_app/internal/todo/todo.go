package todo

import (
	"fmt"
	"sync"
)

type Todo struct {
	ID        int
	Task      string
	Completed bool
}

var (
	inMemoryTodos = []Todo{
		{ID: 1, Task: "Implement in-memory todo list", Completed: true},
		{ID: 2, Task: "Add HTMX for dynamic UI", Completed: false},
		{ID: 3, Task: "Set up Docker for containerization", Completed: false},
	}
	nextID = len(inMemoryTodos) + 1
	mu     sync.RWMutex
)

// initial list todo
var Todos = inMemoryTodos

func GetTodos() []Todo {
	mu.RLock()
	defer mu.RUnlock()
	fmt.Println("GET: ", inMemoryTodos)
	return inMemoryTodos
}

func AddTodo(task string) Todo {
	mu.Lock()
	defer mu.Unlock()
	newTodo := Todo{
		ID:   nextID,
		Task: task,
	}
	inMemoryTodos = append(inMemoryTodos, newTodo)
	fmt.Println("ADD: ", inMemoryTodos)
	nextID++
	return newTodo
}
