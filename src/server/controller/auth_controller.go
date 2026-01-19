package controller

import (
	"net/http"
	"os"

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
	Message string      `json:"message,omitempty"`
}

// isProduction checks if we're running in production
func isProduction() bool {
	return os.Getenv("GIN_MODE") == "release"
}

// setAuthCookie sets the JWT token as an HTTP-only cookie
func setAuthCookie(ctx *gin.Context, token string) {
	secure := isProduction() // Only secure in production (requires HTTPS)

	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie(
		AuthCookieName, // name
		token,          // value
		CookieMaxAge,   // maxAge (24 hours)
		"/",            // path
		"",             // domain (empty = current domain)
		secure,         // secure (HTTPS only in production)
		true,           // httpOnly (not accessible via JavaScript)
	)
}

// clearAuthCookie removes the auth cookie
func clearAuthCookie(ctx *gin.Context) {
	ctx.SetCookie(AuthCookieName, "", -1, "/", "", false, true)
}

func (c *AuthController) Signup(ctx *gin.Context) {
	var req SignupRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := c.authService.Signup(ctx.Request.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		errMsg := err.Error()
		if errMsg == "email already exists" {
			ctx.JSON(http.StatusConflict, gin.H{"error": errMsg})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	// Set token as HTTP-only cookie
	setAuthCookie(ctx, token)

	ctx.JSON(http.StatusCreated, AuthResponse{
		User:    user,
		Message: "Account created successfully",
	})
}

func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := c.authService.Login(ctx.Request.Context(), req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Set token as HTTP-only cookie
	setAuthCookie(ctx, token)

	ctx.JSON(http.StatusOK, AuthResponse{
		User:    user,
		Message: "Login successful",
	})
}

func (c *AuthController) Logout(ctx *gin.Context) {
	clearAuthCookie(ctx)
	ctx.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (c *AuthController) Me(ctx *gin.Context) {
	userID, exists := ctx.Get("auth_user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not authenticated"})
		return
	}

	email, _ := ctx.Get("auth_user_email")
	ctx.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   email,
	})
}
