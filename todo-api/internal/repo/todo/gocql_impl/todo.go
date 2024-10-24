package gocqlimpl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"todo-api/internal/constants"
	"todo-api/internal/model"
	"todo-api/internal/repo/todo"

	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
)

type todoRepositoryImpl struct {
	session       *gocql.Session
	opensearchURL string
}

func NewTodoRepository(session *gocql.Session, opensearchURL string) todo.TodoRepository {
	return &todoRepositoryImpl{session: session, opensearchURL: opensearchURL}
}

// type todoRepositoryImpl struct {
// 	session       *gocql.Session
// 	opensearchURL string
// }

// func NewTodoRepository(session *gocql.Session, opensearchURL string) todo.TodoRepository {
// 	return &todoRepositoryImpl{
// 		session:       session,
// 		opensearchURL: opensearchURL,
// 	}
// }

func (r *todoRepositoryImpl) logToOpenSearch(event, message string) {
	logEntry := map[string]interface{}{
		"event":   event,
		"message": message,
	}

	jsonData, err := json.Marshal(logEntry)
	if err != nil {
		logrus.WithError(err).Error("Failed to marshal log entry for OpenSearch")
		return
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/todo-logs/_doc", r.opensearchURL), bytes.NewBuffer(jsonData))
	if err != nil {
		logrus.WithError(err).Error("Failed to create request for OpenSearch")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth("admin", "Getinet@123!") // Use your credentials

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.WithError(err).Error("Failed to send log entry to OpenSearch")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		logrus.Errorf("OpenSearch responded with status: %s", resp.Status)
	}
}

func (r *todoRepositoryImpl) Save(todo model.Todo) error {
	query := `INSERT INTO todoapp.todos (id, user_id, title, description, status, created, updated) VALUES (?, ?, ?, ?, ?, ?, ?)`
	logrus.WithFields(logrus.Fields{
		"id":          todo.ID,
		"user_id":     todo.UserID,
		"title":       todo.Title,
		"description": todo.Description,
		"status":      todo.Status,
	}).Info("Saving todo")

	r.logToOpenSearch("Saving todo", fmt.Sprintf("ID: %s, UserID: %s", todo.ID, todo.UserID))

	if err := r.session.Query(query, todo.ID, todo.UserID, todo.Title, todo.Description, todo.Status, todo.CreatedAt, todo.UpdatedAt).Consistency(gocql.One).Exec(); err != nil {
		logrus.WithError(err).Error("Error saving todo")
		r.logToOpenSearch("Error saving todo", err.Error())
		return err
	}
	logrus.Info("Todo saved successfully")
	r.logToOpenSearch("Todo saved", fmt.Sprintf("Todo with ID: %s saved successfully", todo.ID))
	return nil
}

func (r *todoRepositoryImpl) FindByID(id string) (model.Todo, error) {
	query := `SELECT id, user_id, title, description, status, created, updated FROM todoapp.todos WHERE id = ?`
	var todo model.Todo

	logrus.WithField("id", id).Info("Finding todo by ID")
	r.logToOpenSearch("Finding todo by ID", id)

	if err := r.session.Query(query, id).Consistency(gocql.One).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
		if err == gocql.ErrNotFound {
			logrus.Warn("Todo not found")
			return model.Todo{}, nil
		}
		logrus.WithError(err).Error("Error finding todo")
		r.logToOpenSearch("Error finding todo", err.Error())
		return model.Todo{}, err
	}

	logrus.Info("Todo found successfully")
	r.logToOpenSearch("Todo found", fmt.Sprintf("Todo with ID: %s found", todo.ID))
	return todo, nil
}

func (r *todoRepositoryImpl) DeleteByID(id string) error {
	query := `DELETE FROM todoapp.todos WHERE id = ?`
	logrus.WithField("id", id).Info("Deleting todo by ID")
	r.logToOpenSearch("Deleting todo by ID", id)

	if err := r.session.Query(query, id).Consistency(gocql.One).Exec(); err != nil {
		logrus.WithError(err).Error("Error deleting todo")
		r.logToOpenSearch("Error deleting todo", err.Error())
		return err
	}
	logrus.Info("Todo deleted successfully")
	r.logToOpenSearch("Todo deleted", fmt.Sprintf("Todo with ID: %s deleted successfully", id))
	return nil
}

func (r *todoRepositoryImpl) ExistsByID(id string) (bool, error) {
	var todoID string
	query := `SELECT id FROM todoapp.todos WHERE id = ? LIMIT 1`
	logrus.WithField("id", id).Info("Checking if todo exists by ID")
	r.logToOpenSearch("Checking if todo exists by ID", id)

	if err := r.session.Query(query, id).Consistency(gocql.One).Scan(&todoID); err != nil {
		if err == gocql.ErrNotFound {
			logrus.Info("Todo does not exist")
			r.logToOpenSearch("Todo does not exist", id)
			return false, nil
		}
		logrus.WithError(err).Error("Error checking todo existence")
		r.logToOpenSearch("Error checking todo existence", err.Error())
		return false, err
	}
	logrus.Info("Todo exists")
	r.logToOpenSearch("Todo exists", fmt.Sprintf("Todo with ID: %s exists", id))
	return true, nil
}

func (r *todoRepositoryImpl) ListTodos(lastID string, limit int, status string, sortOrder string) ([]model.Todo, string, error) {
	var todos []model.Todo
	var nextLastID string

	var query string
	if lastID == "" {
		query = `SELECT id, user_id, title, description, status, created, updated FROM todoapp.todos WHERE status = ? LIMIT ? ALLOW FILTERING`
	} else {
		query = `SELECT id, user_id, title, description, status, created, updated FROM todoapp.todos WHERE id > ? AND status = ? LIMIT ? ALLOW FILTERING`
	}

	logrus.WithFields(logrus.Fields{
		"lastID":    lastID,
		"limit":     limit,
		"status":    status,
		"sortOrder": sortOrder,
	}).Info("Listing todos")

	r.logToOpenSearch("Listing todos", fmt.Sprintf("LastID: %s, Limit: %d, Status: %s", lastID, limit, status))

	q := r.session.Query(query).Consistency(gocql.One)

	if lastID == "" {
		q.Bind(status, limit)
	} else {
		q.Bind(lastID, status, limit)
	}

	iter := q.Iter()
	defer iter.Close()
	var todo model.Todo

	for iter.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt) {
		todos = append(todos, todo)
		nextLastID = todo.ID
	}

	if err := iter.Close(); err != nil {
		logrus.WithError(err).Error("Error closing iterator")
		r.logToOpenSearch("Error closing iterator", err.Error())
		return nil, "", err
	}

	if len(todos) == 0 {
		logrus.Info("No todos found")
		r.logToOpenSearch("No todos found", "No todos match the criteria.")
		return todos, nextLastID, nil // Return empty list if no todos found
	}

	logrus.WithField("count", len(todos)).Info("Todos fetched successfully")
	r.logToOpenSearch("Todos fetched", fmt.Sprintf("Count: %d", len(todos)))

	sort.Slice(todos, func(i, j int) bool {
		if sortOrder == constants.DescS || sortOrder == constants.DescC {
			return todos[i].CreatedAt.After(todos[j].CreatedAt) // Descending order
		}
		return todos[i].CreatedAt.Before(todos[j].CreatedAt) // Ascending order
	})

	return todos, nextLastID, nil
}
