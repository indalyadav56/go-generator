package handlers

import (
	"backend/internal/todo/services"
)

type TodoHandler interface {
	// Add your handler methods here
}

type todoHandler struct {
	service services.TodoService
}

func NewTodoHandler(service services.TodoService) TodoHandler {
	return &todoHandler{
		service: service,
	}
}

// Add your handler implementations here
