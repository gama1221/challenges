package handler

import (
	"errors"
	"net/http"
	"strconv"
	"todo-api/internal/constants"
	"todo-api/internal/handler/request"
	"todo-api/internal/handler/response"
	"todo-api/internal/model"
	"todo-api/internal/service"

	"github.com/gin-gonic/gin"
)

type TodoHandler interface {
	RegisterTodoRoutes(router *gin.RouterGroup)
	createTodo(c *gin.Context)
	getTodoByID(c *gin.Context)
	updateTodo(c *gin.Context)
	deleteTodo(c *gin.Context)
	listTodos(c *gin.Context)
}

type todoHandler struct {
	todoService service.TodoService
}

func NewTodoHandler(t service.TodoService) TodoHandler {
	return &todoHandler{
		todoService: t,
	}
}

func (t todoHandler) RegisterTodoRoutes(router *gin.RouterGroup) {
	router.POST("/todo", t.createTodo)
	router.GET("/todo/:id", t.getTodoByID)
	router.PUT("/todo/:id", t.updateTodo)
	router.DELETE("/todo/:id", t.deleteTodo)
	router.GET("/todos", t.listTodos)

}

func (t todoHandler) createTodo(c *gin.Context) {

	var todoRequest request.TodoRequestBody

	if err := c.ShouldBindJSON(&todoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todoModel := getTodo(todoRequest)

	createdTodo, err := t.todoService.CreateTodo(todoModel)
	if err != nil {
		if errors.Is(err, service.ErrUserDoesNotExist) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdTodo)
}

func (t todoHandler) getTodoByID(c *gin.Context) {
	id := c.Param(constants.Id)

	todo, err := t.todoService.GetTodoByID(id)
	if err != nil {
		if errors.Is(err, service.ErrTodoDoesNotExist) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (t todoHandler) updateTodo(c *gin.Context) {
	id := c.Param(constants.Id)
	var todoRequest request.UpdateTodoRequestBody

	// Extract request body
	if err := c.ShouldBindJSON(&todoRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	todoModel := getUpdateTodo(todoRequest)

	todo, err := t.todoService.UpdateTodoByID(id, todoModel)
	if err != nil {
		if errors.Is(err, service.ErrTodoDoesNotExist) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (t todoHandler) deleteTodo(c *gin.Context) {
	id := c.Param(constants.Id)

	if err := t.todoService.DeleteTodoByID(id); err != nil {
		if errors.Is(err, service.ErrTodoDoesNotExist) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}

func (t todoHandler) listTodos(c *gin.Context) {
	lastID := c.Query(constants.LastId)
	pageSizeStr := c.Query(constants.PageSize)
	status := c.Query(constants.Status)
	sortOrder := c.Query(constants.SortOrder)

	if status == "" {
		status = string(model.StatusPending)
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	todos, nextLastID, err := t.todoService.ListTodos(lastID, pageSize, status, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := response.TodoResponse{
		Todos:      todos,
		NextLastID: nextLastID,
	}

	c.JSON(http.StatusOK, response)
}

func getTodo(todo request.TodoRequestBody) model.Todo {
	return model.Todo{
		UserID:      todo.UserID,
		Title:       todo.Title,
		Description: todo.Description,
	}
}
func getUpdateTodo(todo request.UpdateTodoRequestBody) model.Todo {
	return model.Todo{
		Title:       todo.Title,
		Description: todo.Description,
		Status:      model.TodoStatus(todo.Status),
	}
}
