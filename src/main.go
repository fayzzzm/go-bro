package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/go-learning/src/controller"
	"github.com/user/go-learning/src/repository"
	"github.com/user/go-learning/src/routes"
	"github.com/user/go-learning/src/service"
)

func main() {
	// 1. Setup Database Connection Pool
	ctx := context.Background()
	connStr := "postgres://gouser:gopassword@postgres:5432/godb?sslmode=disable"
	
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	// 2. Initialize Clean Architecture Layers
	
	// Data Layer (Repository)
	userRepo := repository.NewPostgresUserRepo(pool)

	// Business Layer (Service)
	userSvc := service.NewUserService(userRepo)

	// Adapter Layer (Controller)
	userCtrl := controller.NewUserController(userSvc)

	// 3. Framework Layer (Gin)
	r := gin.Default()

	// 4. Setup Routes
	routes.SetupRoutes(r, userCtrl)

	// 5. Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Clean Architecture API starting on :%s...", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
