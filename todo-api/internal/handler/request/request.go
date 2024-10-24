package request

type TodoRequestBody struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpdateTodoRequestBody struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
}
type UserRequestBody struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
