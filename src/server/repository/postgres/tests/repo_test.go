package tests

import (
	"context"
	"os"
	"testing"

	"github.com/fayzzzm/go-bro/repository/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestPostgresRepositories(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		t.Skip("DATABASE_URL not set, skipping integration test")
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// 1. Test AuthRepo
	t.Run("AuthRepo", func(t *testing.T) {
		repo := postgres.NewAuthRepo(pool)

		email := "test-repo@example.com"
		name := "Repo Tester"
		pwd := "hashed-pwd"

		// Signup
		user, err := repo.Signup(ctx, name, email, pwd)
		if err != nil {
			t.Fatalf("Signup failed: %v", err)
		}
		if user.Email != email {
			t.Errorf("Expected email %s, got %s", email, user.Email)
		}

		// GetUserByEmail
		userWithPwd, err := repo.GetUserByEmail(ctx, email)
		if err != nil {
			t.Fatalf("GetUserByEmail failed: %v", err)
		}
		if userWithPwd.PasswordHash != pwd {
			t.Errorf("Expected pwd %s, got %s", pwd, userWithPwd.PasswordHash)
		}
	})

	// 2. Test TodoRepo
	t.Run("TodoRepo", func(t *testing.T) {
		authRepo := postgres.NewAuthRepo(pool)
		todoRepo := postgres.NewTodoRepo(pool)

		// Create a user first
		email := "todo-repo@example.com"
		user, err := authRepo.Signup(ctx, "Todo User", email, "pwd")
		if err != nil {
			// user might already exist from previous run, let's try to get it
			userWithPwd, err2 := authRepo.GetUserByEmail(ctx, email)
			if err2 != nil {
				t.Fatalf("Failed to create/get user for todo test: %v", err2)
			}
			user = &userWithPwd.User
		}

		// Create Todo
		title := "Test Todo Repo"
		desc := "Repo integration test"
		todo, err := todoRepo.Create(ctx, user.ID, title, &desc)
		if err != nil {
			t.Fatalf("Failed to create todo: %v", err)
		}
		if todo.Title != title {
			t.Errorf("Expected title %s, got %s", title, todo.Title)
		}

		// List Todos
		todos, err := todoRepo.GetByUser(ctx, user.ID, 10, 0)
		if err != nil {
			t.Fatalf("Failed to list todos: %v", err)
		}
		if len(todos) == 0 {
			t.Error("Expected at least one todo, got none")
		}

		// Toggle Todo
		updated, err := todoRepo.Toggle(ctx, todo.ID, user.ID)
		if err != nil {
			t.Fatalf("Failed to toggle todo: %v", err)
		}
		if updated.Completed != !todo.Completed {
			t.Error("Toggle did not change completed status")
		}

		// Delete Todo
		err = todoRepo.Delete(ctx, todo.ID, user.ID)
		if err != nil {
			t.Fatalf("Failed to delete todo: %v", err)
		}
	})
}
