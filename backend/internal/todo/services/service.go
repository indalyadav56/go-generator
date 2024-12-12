package services

import (
	"backend/internal/todo/models"
	"backend/internal/todo/repository"
	"backend/pkg/logger"
)

type TodoService interface {
	Create(todo *models.Todo) (*models.Todo, error)
	GetByID(id string) (*models.Todo, error)
	Update(id string, todo *models.Todo) (*models.Todo, error)
	Delete(id string) error
	GetAll() ([]*models.Todo, error)
	Search(query string) ([]*models.Todo, error)
}

type todoService struct {
	todoRepo repository.TodoRepository
	log      logger.Logger
}

func NewTodoService(repo repository.TodoRepository, log logger.Logger) *todoService {
	return &todoService{
		todoRepo: repo,
		log:      log,
	}
}

// Create creates a new todo
func (s *todoService) Create(todo *models.Todo) (*models.Todo, error) {
	return s.todoRepo.Insert(todo)
}

// GetByID retrieves a todo by its ID
func (s *todoService) GetByID(id string) (*models.Todo, error) {
	return s.todoRepo.FindByID(id)
}

// Update updates an existing todo by its ID
func (s *todoService) Update(id string, todo *models.Todo) (*models.Todo, error) {
	existingTodo, err := s.todoRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields (example: only certain fields are updated)
	existingTodo.Field1 = todo.Field1
	existingTodo.Field2 = todo.Field2

	return s.todoRepo.Update(id, existingTodo)
}

// Delete removes a todo by its ID
func (s *todoService) Delete(id string) error {
	return s.todoRepo.Delete(id)
}

// GetAll retrieves all todo records
func (s *todoService) GetAll() ([]*models.Todo, error) {
	return s.todoRepo.FindAll()
}

// Search allows searching for todo entities based on a query
func (s *todoService) Search(query string) ([]*models.Todo, error) {
	return s.todoRepo.Search(query)
}
