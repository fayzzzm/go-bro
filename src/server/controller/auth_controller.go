package controller

import (
	"context"
	"net/http"
	"os"

	"github.com/fayzzzm/go-bro/middleware"
	"github.com/fayzzzm/go-bro/pkg/reply"
	"github.com/gin-gonic/gin"
)

const (
	AuthCookieName = "auth_token"
	CookieMaxAge   = 24 * 60 * 60 // 24 hours in seconds
)

type authUseCase interface {
	Signup(ctx context.Context, name, email, password string) (interface{}, error)
	Login(ctx context.Context, email, password string) (interface{}, string, error)
}

type AuthController struct {
	usecase authUseCase
}

func NewAuthController(usecase authUseCase) *AuthController {
	return &AuthController{usecase: usecase}
}

type SignupRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	User    interface{} `json:"user"`
	Token   string      `json:"token,omitempty"`
	Message string      `json:"message,omitempty"`
}

func isProduction() bool {
	return os.Getenv("GIN_MODE") == "release"
}

func setAuthCookie(ctx *gin.Context, token string) {
	secure := isProduction()
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(AuthCookieName, token, CookieMaxAge, "/", "", secure, true)
}

func clearAuthCookie(ctx *gin.Context) {
	ctx.SetCookie(AuthCookieName, "", -1, "/", "", false, true)
}

func (c *AuthController) Signup(ctx *gin.Context) {
	req := middleware.GetBody[SignupRequest](ctx)

	user, err := c.usecase.Signup(ctx.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "email already exists" {
			code = http.StatusConflict
		}
		reply.Error(ctx, code, err.Error(), err)
		return
	}

	reply.Created(ctx, AuthResponse{
		User:    user,
		Message: "Account created successfully",
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	req := middleware.GetBody[LoginRequest](ctx)

	user, token, err := c.usecase.Login(ctx.Request.Context(), req.Email, req.Password)
	if reply.Error(ctx, http.StatusUnauthorized, "invalid credentials", err) {
		return
	}

	setAuthCookie(ctx, token)

	reply.OK(ctx, AuthResponse{
		User:    user,
		Token:   token,
		Message: "Login successful",
	})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	clearAuthCookie(ctx)
	reply.OK(ctx, gin.H{"message": "Logged out successfully"})
}

func (c *AuthController) Me(ctx *gin.Context) {
	userID, _ := middleware.GetUserID(ctx)
	email, _ := middleware.GetUserEmail(ctx)

	reply.OK(ctx, gin.H{
		"user_id": userID,
		"email":   email,
	})
}
