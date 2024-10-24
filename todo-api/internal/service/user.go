package service

import (
	"errors"
	"time"
	"todo-api/internal/model"
	"todo-api/internal/repo/user"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserDoesNotExist = errors.New("user does not exist")

type UserService interface {
	CreateUser(user model.User) (model.User, error)
}

type userService struct {
	userRepo user.UserRepository
}

func NewUserService(userRepo user.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CheckUserExists(userID string) (bool, error) {
	exists, err := s.userRepo.FindByID(userID)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (s *userService) CreateUser(user model.User) (model.User, error) {
	exists, _ := s.userRepo.FindByUsername(user.Username)
	if exists {
		return model.User{}, ErrUserAlreadyExists
	}

	user.UserID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return user, s.userRepo.Save(user)
}
