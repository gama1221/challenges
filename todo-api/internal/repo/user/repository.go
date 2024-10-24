package user

import (
	"todo-api/internal/model"
)

type UserRepository interface {
	FindByID(id string) (bool, error)
	FindByUsername(username string) (bool, error)
	Save(user model.User) error
}
