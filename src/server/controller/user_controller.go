package controller

import (
	"net/http"
	"strconv"

	usecase "github.com/fayzzzm/go-bro/usecase/user"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	uc usecase.UserUseCase
}

func NewUserController(uc usecase.UserUseCase) *UserController {
	return &UserController{uc: uc}
}

func (ctrl *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	input := usecase.GetUserInput{ID: id}
	output, err := ctrl.uc.GetUser(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !output.Found {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, output.User)
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var input usecase.RegisterUserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	output, err := ctrl.uc.RegisterUser(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error":   err.Error(),
			"success": false,
		})
		return
	}

	c.JSON(http.StatusCreated, output)
}

func (ctrl *UserController) ListUsers(c *gin.Context) {
	// Parse optional query parameters
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	input := usecase.ListUsersInput{
		Limit:  limit,
		Offset: offset,
	}

	output, err := ctrl.uc.ListUsers(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, output)
}
