package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/fayzzzm/go-bro/middleware"
	"github.com/fayzzzm/go-bro/pkg/reply"
	"github.com/fayzzzm/go-bro/usecase/users"
	"github.com/gin-gonic/gin"
)

type userUseCase interface {
	RegisterUser(ctx context.Context, input users.RegisterUserInput) (*users.RegisterUserOutput, error)
	GetUser(ctx context.Context, input users.GetUserInput) (*users.GetUserOutput, error)
	ListUsers(ctx context.Context, input users.ListUsersInput) (*users.ListUsersOutput, error)
}

type UserController struct {
	usecase userUseCase
}

func NewUserController(usecase userUseCase) *UserController {
	return &UserController{usecase: usecase}
}

func (c *UserController) GetUser(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if reply.Error(ctx, http.StatusBadRequest, "invalid id format", err) {
		return
	}

	input := users.GetUserInput{ID: id}
	output, err := c.usecase.GetUser(ctx.Request.Context(), input)

	if reply.Error(ctx, http.StatusNotFound, "user not found", err) {
		return
	}

	reply.OK(ctx, output.User)
}

func (c *UserController) CreateUser(ctx *gin.Context) {
	input := middleware.GetBody[users.RegisterUserInput](ctx)

	output, err := c.usecase.RegisterUser(ctx.Request.Context(), input)
	if reply.Error(ctx, http.StatusUnprocessableEntity, "could not create user", err) {
		return
	}

	reply.Created(ctx, output)
}

func (c *UserController) ListUsers(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	input := users.ListUsersInput{
		Limit:  limit,
		Offset: offset,
	}

	output, err := c.usecase.ListUsers(ctx.Request.Context(), input)
	if reply.InternalError(ctx, err) {
		return
	}

	reply.OK(ctx, output)
}
