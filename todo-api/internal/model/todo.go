package model

import (
	"time"
)

type Todo struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Status      TodoStatus `json:"status"`
	CreatedAt   time.Time  `json:"created"`
	UpdatedAt   time.Time  `json:"updated"`
}

type TodoStatus string

const (
	StatusPending    TodoStatus = "Pending"
	StatusInProgress TodoStatus = "In Progress"
	StatusCompleted  TodoStatus = "Completed"
	StatusArchived   TodoStatus = "Archived"
	StatusCancelled  TodoStatus = "Cancelled"
)
