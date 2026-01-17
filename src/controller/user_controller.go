package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/fayzzzm/go-bro/src/models"
	"net/http"
	"strconv"
)

// UserServicer defines the expected behavior of the business logic.
// In Hexagonal Architecture, this is an "Input Port".
type UserServicer interface {
	RegisterUser(name, email string) (*models.User, error)
	GetUser(id int) (*models.User, error)
	ListUsers() ([]*models.User, error)
}

type UserController struct {
	svc UserServicer
}

func NewUserController(svc UserServicer) *UserController {
	return &UserController{svc: svc}
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	user, err := ctrl.svc.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.svc.RegisterUser(input.Name, input.Email)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (ctrl *UserController) ListUsers(c *gin.Context) {
	users, err := ctrl.svc.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
