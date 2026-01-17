package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/fayzzzm/go-bro/src/controller"
	"github.com/fayzzzm/go-bro/src/repository/postgres"
	"github.com/fayzzzm/go-bro/src/routes"
	"github.com/fayzzzm/go-bro/src/service"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			// 1. Database Pool
			NewDatabasePool,
			// 2. Repositories (Adapters) - Isolated in their own sub-package
			fx.Annotate(
				postgres.NewUserRepo,
				fx.As(new(service.UserRepository)),
			),
			// 3. Services (Core)
			fx.Annotate(
				service.NewUserService,
				fx.As(new(controller.UserServicer)),
			),
			// 4. Controllers (Adapters)
			controller.NewUserController,
			// 5. Framework (Gin)
			NewGinEngine,
		),
		fx.Invoke(
			// 6. Setup Routes and start server logic
			RegisterRoutes,
		),
	).Run()
}

// NewDatabasePool creates a connection pool and handles its shutdown via fx.Lifecycle
func NewDatabasePool(lc fx.Lifecycle) (*pgxpool.Pool, error) {
	ctx := context.Background()
	connStr := "postgres://gouser:gopassword@postgres:5432/godb?sslmode=disable"
	
	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing database pool...")
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

// NewGinEngine initializes the Gin framework
func NewGinEngine() *gin.Engine {
	r := gin.Default()
	return r
}

// RegisterRoutes ties everything together and starts the HTTP server
func RegisterRoutes(lc fx.Lifecycle, r *gin.Engine, userCtrl *controller.UserController) {
	routes.SetupRoutes(r, userCtrl)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("🚀 Clean Architecture API starting on :%s (via fx)...", port)
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatalf("Failed to start server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down HTTP server...")
			return server.Shutdown(ctx)
		},
	})
}
