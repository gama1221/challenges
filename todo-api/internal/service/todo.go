package service

import (
	"errors"
	"time"
	"todo-api/internal/constants"
	"todo-api/internal/model"
	"todo-api/internal/repo/todo"
	"todo-api/internal/repo/user"

	"github.com/google/uuid"
)

var ErrTodoDoesNotExist = errors.New("todo does not exist")
var ErrInvalidStatusTransition = errors.New("invalid status transition")

type TodoService interface {
	CreateTodo(todo model.Todo) (model.Todo, error)
	GetTodoByID(id string) (model.Todo, error)
	UpdateTodoByID(id string, todoRequest model.Todo) (model.Todo, error)
	DeleteTodoByID(id string) error
	ListTodos(lastID string, pageSize int, status string, sortOrder string) ([]model.Todo, string, error)
}

type todoService struct {
	todoRepo todo.TodoRepository
	userRepo user.UserRepository
}

func NewTodoService(todoRepo todo.TodoRepository, userRepo user.UserRepository) TodoService {
	return &todoService{
		todoRepo: todoRepo,
		userRepo: userRepo,
	}
}

func (t *todoService) CreateTodo(todo model.Todo) (model.Todo, error) {
	userExists, err := t.userRepo.FindByID(todo.UserID)
	if !userExists {
		return model.Todo{}, ErrUserDoesNotExist
	}

	todo.Status = model.StatusPending

	if err != nil {
		return model.Todo{}, err
	}

	todo.ID = uuid.New().String()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	if err := t.todoRepo.Save(todo); err != nil {
		return model.Todo{}, err
	}

	return todo, nil
}

func (t *todoService) GetTodoByID(id string) (model.Todo, error) {
	todo, err := t.todoRepo.FindByID(id)
	if err != nil {
		return model.Todo{}, ErrTodoDoesNotExist
	}

	return todo, nil
}

func (t *todoService) UpdateTodoByID(id string, todoRequest model.Todo) (model.Todo, error) {
	existingTodo, err := t.todoRepo.FindByID(id)
	if err != nil {
		return model.Todo{}, ErrTodoDoesNotExist
	}

	if err := validateStatusTransition(string(existingTodo.Status), string(todoRequest.Status)); err != nil {
		return model.Todo{}, err
	}

	existingTodo.Title = todoRequest.Title
	existingTodo.Description = todoRequest.Description
	existingTodo.Status = todoRequest.Status
	existingTodo.UpdatedAt = time.Now()

	if err := t.todoRepo.Save(existingTodo); err != nil {
		return model.Todo{}, err
	}

	return existingTodo, nil
}

func (t *todoService) DeleteTodoByID(id string) error {
	todoExists, err := t.todoRepo.ExistsByID(id)
	if err != nil {
		return err
	}

	if !todoExists {
		return ErrTodoDoesNotExist
	}

	err = t.todoRepo.DeleteByID(id)
	if err != nil {
		return err
	}
	return nil
}

func (t *todoService) ListTodos(lastID string, pageSize int, status string, sortOrder string) ([]model.Todo, string, error) {
	todos, nextLastID, err := t.todoRepo.ListTodos(lastID, pageSize, status, sortOrder)
	if err != nil {
		return nil, "", err
	}

	return todos, nextLastID, nil
}

func validateStatusTransition(currentStatus, newStatus string) error {
	switch currentStatus {
	case constants.Pending:
		if newStatus != constants.InProgres {
			return ErrInvalidStatusTransition
		}
	case constants.InProgres:
		if newStatus != constants.Completed && newStatus != constants.Archived && newStatus != constants.Cancelled {
			return ErrInvalidStatusTransition
		}
	case constants.Completed:
		if newStatus != constants.Archived {
			return ErrInvalidStatusTransition
		}
	case constants.Archived:
		if newStatus != constants.InProgres && newStatus != constants.Pending {
			return ErrInvalidStatusTransition
		}
	case constants.Cancelled:
		return ErrInvalidStatusTransition
	}

	return nil
}
