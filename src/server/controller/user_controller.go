package controller

import (
	"net/http"
	"strconv"

	"github.com/fayzzzm/go-bro/middleware"
	"github.com/fayzzzm/go-bro/pkg/reply"
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
	id, err := strconv.Atoi(c.Param("id"))
	if reply.Error(c, http.StatusBadRequest, "invalid id format", err) {
		return
	}

	input := usecase.GetUserInput{ID: id}
	output, err := ctrl.uc.GetUser(c.Request.Context(), input)

	if reply.Error(c, http.StatusNotFound, "user not found", err) {
		return
	}

	reply.OK(c, output.User)
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	input := middleware.GetBody[usecase.RegisterUserInput](c)

	output, err := ctrl.uc.RegisterUser(c.Request.Context(), input)
	if reply.Error(c, http.StatusUnprocessableEntity, "could not create user", err) {
		return
	}

	reply.Created(c, output)
}

func (ctrl *UserController) ListUsers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	input := usecase.ListUsersInput{
		Limit:  limit,
		Offset: offset,
	}

	output, err := ctrl.uc.ListUsers(c.Request.Context(), input)
	if reply.InternalError(c, err) {
		return
	}

	reply.OK(c, output)
}
