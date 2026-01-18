package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/fayzzzm/go-bro/controller"
)

func SetupRoutes(r *gin.Engine, userCtrl *controller.UserController) {
	v1 := r.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", userCtrl.ListUsers)
			users.GET("/:id", userCtrl.GetUser)
			users.POST("", userCtrl.CreateUser)
		}
	}
}
