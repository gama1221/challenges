package todo

import "todo-api/internal/model"

type TodoRepository interface {
	Save(todo model.Todo) error
	FindByID(id string) (model.Todo, error)
	DeleteByID(id string) error
	ExistsByID(id string) (bool, error)
	ListTodos(lastID string, limit int,status string, sortOrder string) ([]model.Todo, string, error)
}
