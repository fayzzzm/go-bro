package controller

import (
	"context"
	"strconv"

	"github.com/fayzzzm/go-bro/middleware"
	"github.com/fayzzzm/go-bro/models"
	"github.com/fayzzzm/go-bro/pkg/reply"
	"github.com/gin-gonic/gin"
)

type todoUseCase interface {
	Create(ctx context.Context, userID int, title string, description *string) (*models.Todo, error)
	GetByUser(ctx context.Context, userID, limit, offset int) ([]models.Todo, error)
	GetByID(ctx context.Context, id, userID int) (*models.Todo, error)
	Update(ctx context.Context, userID int, todo *models.Todo) (*models.Todo, error)
	Delete(ctx context.Context, id, userID int) error
	Toggle(ctx context.Context, id, userID int) (*models.Todo, error)
}

type TodoController struct {
	usecase todoUseCase
}

func NewTodoController(usecase todoUseCase) *TodoController {
	return &TodoController{usecase: usecase}
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
	userID, _ := middleware.GetUserID(ctx)
	req := middleware.GetBody[CreateTodoRequest](ctx)

	todo, err := c.usecase.Create(ctx.Request.Context(), userID, req.Title, req.Description)
	if reply.InternalError(ctx, err) {
		return
	}

	reply.Created(ctx, gin.H{"todo": todo})
}

func (c *TodoController) List(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	todos, err := c.usecase.GetByUser(ctx.Request.Context(), userID, limit, offset)
	if reply.InternalError(ctx, err) {
		return
	}

	if todos == nil {
		todos = []models.Todo{}
	}

	reply.OK(ctx, gin.H{"todos": todos, "count": len(todos)})
}

func (c *TodoController) GetByID(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)
	todoID, _ := strconv.Atoi(ctx.Param("id"))

	todo, err := c.usecase.GetByID(ctx.Request.Context(), todoID, userID)
	if reply.NotFound(ctx, err) {
		return
	}

	reply.OK(ctx, gin.H{"todo": todo})
}

func (c *TodoController) Update(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)
	todoID, _ := strconv.Atoi(ctx.Param("id"))
	req := middleware.GetBody[UpdateTodoRequest](ctx)

	todo := &models.Todo{
		ID:          todoID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	}

	updatedTodo, err := c.usecase.Update(ctx.Request.Context(), userID, todo)
	if reply.NotFound(ctx, err) {
		return
	}

	reply.OK(ctx, gin.H{"todo": updatedTodo})
}

func (c *TodoController) Delete(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)
	todoID, _ := strconv.Atoi(ctx.Param("id"))

	if err := c.usecase.Delete(ctx.Request.Context(), todoID, userID); err != nil {
		reply.NotFound(ctx, err)
		return
	}

	reply.OK(ctx, gin.H{"message": "todo deleted"})
}

func (c *TodoController) Toggle(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)
	todoID, _ := strconv.Atoi(ctx.Param("id"))

	todo, err := c.usecase.Toggle(ctx.Request.Context(), todoID, userID)
	if reply.NotFound(ctx, err) {
		return
	}

	reply.OK(ctx, gin.H{"todo": todo})
}
