package impl

import (
	"errors"
	"fmt"
	"todo-api/internal/model"
	"todo-api/internal/repo/user"

	"github.com/gocql/gocql"
)

var ErrUserNotFound = errors.New("user not found")
var ErrUserDoesNotExist = errors.New("user does not exist")

type userRepositoryImpl struct {
	session *gocql.Session
}

func NewUserRepository(session *gocql.Session) user.UserRepository {
	return &userRepositoryImpl{session: session}
}

func (r *userRepositoryImpl) FindByUsername(username string) (bool, error) {
	query := `SELECT user_id FROM todoapp.users WHERE username = ? LIMIT 1 ALLOW FILTERING`
	var userID string

	if err := r.session.Query(query, username).Consistency(gocql.One).Scan(&userID); err != nil {
		if err == gocql.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *userRepositoryImpl) FindByID(id string) (bool, error) {
	query := "SELECT user_id FROM todoapp.users WHERE user_id = ? LIMIT 1"
	var userID string

	if err := r.session.Query(query, id).Consistency(gocql.One).Scan(&userID); err != nil {
		if err != nil {
			fmt.Println(err)
			return false, ErrUserDoesNotExist
		}
		return false, err
	}

	return true, nil
}

func (r *userRepositoryImpl) Save(user model.User) error {
	query := `INSERT INTO todoapp.users (user_id, username, email, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`
	if err := r.session.Query(query, user.UserID, user.Username, user.Email, user.CreatedAt, user.UpdatedAt).Consistency(gocql.One).Exec(); err != nil {
		return err
	}
	return nil
}
