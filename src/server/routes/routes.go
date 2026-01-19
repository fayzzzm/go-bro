package routes

import (
	"github.com/fayzzzm/go-bro/controller"
	"github.com/fayzzzm/go-bro/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	userCtrl *controller.UserController,
	authCtrl *controller.AuthController,
	todoCtrl *controller.TodoController,
) {
	// API v1 group
	api := r.Group("/api/v1")

	// Public routes (no auth required)
	auth := api.Group("/auth")
	{
		auth.POST("/signup", authCtrl.Signup)
		auth.POST("/login", authCtrl.Login)
		auth.POST("/logout", authCtrl.Logout)
	}

	// User routes (public for now, can be protected later)
	users := api.Group("/users")
	{
		users.GET("", userCtrl.ListUsers)
		users.GET("/:id", userCtrl.GetUser)
		users.POST("", userCtrl.CreateUser)
	}

	// Protected routes (require auth)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// Auth
		protected.GET("/me", authCtrl.Me)

		// Todos
		todos := protected.Group("/todos")
		{
			todos.GET("", todoCtrl.List)
			todos.POST("", todoCtrl.Create)
			todos.GET("/:id", todoCtrl.GetByID)
			todos.PUT("/:id", todoCtrl.Update)
			todos.DELETE("/:id", todoCtrl.Delete)
			todos.PATCH("/:id/toggle", todoCtrl.Toggle)
		}
	}
}
