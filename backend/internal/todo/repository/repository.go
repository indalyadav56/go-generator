package repository

import (
	"backend/internal/todo/models"
	"backend/pkg/logger"
	"database/sql"
	"fmt"
)

type TodoRepository interface {
	Insert(todo *models.Todo) (*models.Todo, error)
	Update(todo *models.Todo) (*models.Todo, error)
	FindByID(id string) (*models.Todo, error)
	List(page, pageSize int) ([]models.Todo, error)
	Delete(id string) error
}

type todoRepository struct {
	db  *sql.DB
	log logger.Logger
}

func NewTodoRepository(db *sql.DB, log logger.Logger) *todoRepository {
	return &todoRepository{
		db:  db,
		log: log,
	}
}

// Insert inserts a new record into the database
func (r *todoRepository) Insert(todo *models.Todo) (*models.Todo, error) {
	query := "INSERT INTO todos (field1, field2) VALUES (?, ?) RETURNING id"
	err := r.db.QueryRow(query, todo.Field1, todo.Field2).Scan(&todo.ID)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// FindByID retrieves a record by its ID from the database
func (r *todoRepository) FindByID(id string) (*models.Todo, error) {
	query := "SELECT id, field1, field2 FROM todos WHERE id = ?"
	row := r.db.QueryRow(query, id)

	var todo models.Todo
	if err := row.Scan(&todo.ID, &todo.Field1, &todo.Field2); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &todo, nil
}

// List retrieves a paginated list of records from the database
func (r *todoRepository) List(page, pageSize int) ([]models.Todo, error) {
	offset := (page - 1) * pageSize
	query := "SELECT id, field1, field2 FROM todos LIMIT ? OFFSET ?"
	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.ID, &todo.Field1, &todo.Field2); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// Update updates an existing record in the database
func (r *todoRepository) Update(todo *models.Todo) (*models.Todo, error) {
	query := "UPDATE todos SET field1 = ?, field2 = ? WHERE id = ?"
	_, err := r.db.Exec(query, todo.Field1, todo.Field2, todo.ID)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

// Delete removes a record from the database by ID
func (r *todoRepository) Delete(id string) error {
	query := "DELETE FROM todos WHERE id = ?"
	_, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
