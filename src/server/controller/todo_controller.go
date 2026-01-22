package controller

import (
	"strconv"

	"github.com/fayzzzm/go-bro/middleware"
	"github.com/fayzzzm/go-bro/models"
	"github.com/fayzzzm/go-bro/pkg/reply"
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

func (ctrl *TodoController) Create(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	req := middleware.GetBody[CreateTodoRequest](c)

	todo, err := ctrl.todoService.Create(c.Request.Context(), userID, req.Title, req.Description)
	if reply.InternalError(c, err) {
		return
	}

	reply.Created(c, gin.H{"todo": todo})
}

func (ctrl *TodoController) List(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	todos, err := ctrl.todoService.GetByUser(c.Request.Context(), userID, limit, offset)
	if reply.InternalError(c, err) {
		return
	}

	if todos == nil {
		todos = []models.Todo{}
	}

	reply.OK(c, gin.H{"todos": todos, "count": len(todos)})
}

func (ctrl *TodoController) GetByID(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	todoID, _ := strconv.Atoi(c.Param("id"))

	todo, err := ctrl.todoService.GetByID(c.Request.Context(), todoID, userID)
	if reply.NotFound(c, err) {
		return
	}

	reply.OK(c, gin.H{"todo": todo})
}

func (ctrl *TodoController) Update(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	todoID, _ := strconv.Atoi(c.Param("id"))
	req := middleware.GetBody[UpdateTodoRequest](c)

	todo := &models.Todo{
		ID:          todoID,
		Title:       req.Title,
		Description: req.Description,
		Completed:   req.Completed,
	}

	updatedTodo, err := ctrl.todoService.Update(c.Request.Context(), userID, todo)
	if reply.NotFound(c, err) {
		return
	}

	reply.OK(c, gin.H{"todo": updatedTodo})
}

func (ctrl *TodoController) Delete(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	todoID, _ := strconv.Atoi(c.Param("id"))

	if err := ctrl.todoService.Delete(c.Request.Context(), todoID, userID); err != nil {
		reply.NotFound(c, err)
		return
	}

	reply.OK(c, gin.H{"message": "todo deleted"})
}

func (ctrl *TodoController) Toggle(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	todoID, _ := strconv.Atoi(c.Param("id"))

	todo, err := ctrl.todoService.Toggle(c.Request.Context(), todoID, userID)
	if reply.NotFound(c, err) {
		return
	}

	reply.OK(c, gin.H{"todo": todo})
}
