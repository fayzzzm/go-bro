package controller

import (
	"net/http"
	"strconv"

	"github.com/fayzzzm/go-bro/middleware"
	"github.com/fayzzzm/go-bro/models"
	"github.com/fayzzzm/go-bro/service"
	"github.com/gin-gonic/gin"
)

type TodoController struct {
	todoService service.TodoServicer
}

func NewTodoController(todoService service.TodoServicer) *TodoController {
	return &TodoController{todoService: todoService}
}

type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

func (c *TodoController) Create(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx) // Auth middleware already verified this
	req := middleware.GetBody[CreateTodoRequest](ctx)

	todo, err := c.todoService.Create(ctx.Request.Context(), userID, req.Title, req.Description)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"todo": todo})
}

func (c *TodoController) List(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	todos, err := c.todoService.GetByUser(ctx.Request.Context(), userID, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if todos == nil {
		todos = []models.Todo{}
	}

	ctx.JSON(http.StatusOK, gin.H{"todos": todos, "count": len(todos)})
}

func (c *TodoController) GetByID(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)

	todoID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	todo, err := c.todoService.GetByID(ctx.Request.Context(), todoID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"todo": todo})
}

func (c *TodoController) Update(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)
	todoID, _ := strconv.Atoi(ctx.Param("id")) // Param check could be a middleware too, but keeping it simple
	req := middleware.GetBody[UpdateTodoRequest](ctx)

	todo, err := c.todoService.Update(ctx.Request.Context(), todoID, userID, req.Title, req.Description, req.Completed)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"todo": todo})
}

func (c *TodoController) Delete(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)

	todoID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	if err := c.todoService.Delete(ctx.Request.Context(), todoID, userID); err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "todo deleted"})
}

func (c *TodoController) Toggle(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)

	todoID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid todo id"})
		return
	}

	todo, err := c.todoService.Toggle(ctx.Request.Context(), todoID, userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "todo not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"todo": todo})
}
