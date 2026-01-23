package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/fayzzzm/go-bro/controller"
	"github.com/fayzzzm/go-bro/repository/postgres"
	"github.com/fayzzzm/go-bro/routes"
	"github.com/fayzzzm/go-bro/service"
	"github.com/fayzzzm/go-bro/usecase/users"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			// 1. Database Pool
			NewDatabasePool,

			// 2. Repositories (Adapters)
			fx.Annotate(
				postgres.NewUserRepo,
				fx.As(new(service.UserRepository)),
			),
			fx.Annotate(
				postgres.NewAuthRepo,
				fx.As(new(service.AuthRepository)),
			),
			fx.Annotate(
				postgres.NewTodoRepo,
				fx.As(new(service.TodoRepository)),
			),

			// 3. Services (Core)
			fx.Annotate(
				service.NewUserService,
				fx.As(new(service.UserServicer)),
			),
			fx.Annotate(
				service.NewAuthService,
				fx.As(new(controller.AuthUseCase)),
			),
			fx.Annotate(
				service.NewTodoService,
				fx.As(new(controller.TodoUseCase)),
			),

			// 4. Use Cases
			fx.Annotate(
				users.NewUseCase,
				fx.As(new(controller.UserUseCase)),
			),

			// 5. Controllers (Adapters)
			controller.NewUserController,
			controller.NewAuthController,
			controller.NewTodoController,

			// 6. Framework (Gin)
			NewGinEngine,
		),
		fx.Invoke(
			// 7. Setup Routes and start server
			RegisterRoutes,
		),
	).Run()
}

// NewDatabasePool creates a connection pool and handles its shutdown via fx.Lifecycle
func NewDatabasePool(lc fx.Lifecycle) (*pgxpool.Pool, error) {
	ctx := context.Background()

	// Get connection string from environment variable
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	// Register custom composite types
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		types := []string{"users.user_request", "todos.todo_request"}
		for _, t := range types {
			dt, err := conn.LoadType(ctx, t)
			if err != nil {
				log.Printf("⚠️ Warning: Failed to load type %s: %v", t, err)
				continue
			}
			conn.TypeMap().RegisterType(dt)
		}
		return nil
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	log.Println("✅ Database connection established")

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing database pool...")
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

// NewGinEngine initializes the Gin framework with CORS
func NewGinEngine() *gin.Engine {
	r := gin.Default()

	// Enable CORS for development
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	return r
}

// RegisterRoutes ties everything together and starts the HTTP server
func RegisterRoutes(
	lc fx.Lifecycle,
	r *gin.Engine,
	userCtrl *controller.UserController,
	authCtrl *controller.AuthController,
	todoCtrl *controller.TodoController,
) {
	routes.SetupRoutes(r, userCtrl, authCtrl, todoCtrl)

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
			log.Printf("🚀 Todo API starting on :%s", port)
			log.Println("📦 Endpoints: /api/v1/auth, /api/v1/todos, /api/v1/users")
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
