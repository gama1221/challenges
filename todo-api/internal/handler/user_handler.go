package handler

import (
	"net/http"
	"todo-api/internal/handler/request"
	"todo-api/internal/model"
	"todo-api/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("/users", h.createUser)
}

func (h *UserHandler) createUser(c *gin.Context) {
	var userRequest request.UserRequestBody
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := getUser(userRequest)

	u, err := h.userService.CreateUser(user)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}

// getUser converts UserRequestBody to User model
func getUser(userRequest request.UserRequestBody) model.User {
	return model.User{
		Username: userRequest.Username,
		Email:    userRequest.Email,
	}
}
