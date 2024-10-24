package response

import "todo-api/internal/model"

type TodoResponse struct {
	Todos      []model.Todo `json:"todos"`
	NextLastID string       `json:"next_last_id,omitempty"`
}
