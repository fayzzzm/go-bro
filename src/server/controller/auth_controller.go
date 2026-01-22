package controller

import (
	"net/http"
	"os"

	"github.com/fayzzzm/go-bro/middleware"
	"github.com/fayzzzm/go-bro/pkg/reply"
	"github.com/fayzzzm/go-bro/service"
	"github.com/gin-gonic/gin"
)

const (
	AuthCookieName = "auth_token"
	CookieMaxAge   = 24 * 60 * 60 // 24 hours in seconds
)

type AuthController struct {
	authService service.AuthServicer
}

func NewAuthController(authService service.AuthServicer) *AuthController {
	return &AuthController{authService: authService}
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

func setAuthCookie(c *gin.Context, token string) {
	secure := isProduction()
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(AuthCookieName, token, CookieMaxAge, "/", "", secure, true)
}

func clearAuthCookie(c *gin.Context) {
	c.SetCookie(AuthCookieName, "", -1, "/", "", false, true)
}

func (ctrl *AuthController) Signup(c *gin.Context) {
	req := middleware.GetBody[SignupRequest](c)

	user, err := ctrl.authService.Signup(c.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		code := http.StatusBadRequest
		if err.Error() == "email already exists" {
			code = http.StatusConflict
		}
		reply.Error(c, code, err.Error(), err)
		return
	}

	reply.Created(c, AuthResponse{
		User:    user,
		Message: "Account created successfully",
	})
}

func (ctrl *AuthController) Login(c *gin.Context) {
	req := middleware.GetBody[LoginRequest](c)

	user, token, err := ctrl.authService.Login(c.Request.Context(), req.Email, req.Password)
	if reply.Error(c, http.StatusUnauthorized, "invalid credentials", err) {
		return
	}

	setAuthCookie(c, token)

	reply.OK(c, AuthResponse{
		User:    user,
		Token:   token,
		Message: "Login successful",
	})
}

func (ctrl *AuthController) Logout(c *gin.Context) {
	clearAuthCookie(c)
	reply.OK(c, gin.H{"message": "Logged out successfully"})
}

func (ctrl *AuthController) Me(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)
	email, _ := middleware.GetUserEmail(c)

	reply.OK(c, gin.H{
		"user_id": userID,
		"email":   email,
	})
}
