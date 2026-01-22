package tests

import (
	"context"
	"testing"

	"github.com/fayzzzm/go-bro/models"
	"github.com/fayzzzm/go-bro/service"
	"golang.org/x/crypto/bcrypt"
)

// MockAuthRepository implements service.AuthRepository
type MockAuthRepository struct {
	SignupFunc         func(ctx context.Context, name, email, hashedPassword string) (*models.User, error)
	GetUserByEmailFunc func(ctx context.Context, email string) (*models.UserWithPassword, error)
}

func (m *MockAuthRepository) Signup(ctx context.Context, name, email, hashedPassword string) (*models.User, error) {
	return m.SignupFunc(ctx, name, email, hashedPassword)
}

func (m *MockAuthRepository) GetUserByEmail(ctx context.Context, email string) (*models.UserWithPassword, error) {
	return m.GetUserByEmailFunc(ctx, email)
}

func TestAuthService_Signup(t *testing.T) {
	mockRepo := &MockAuthRepository{
		SignupFunc: func(ctx context.Context, name, email, hashedPassword string) (*models.User, error) {
			return &models.User{ID: 1, Name: name, Email: email}, nil
		},
	}
	authService := service.NewAuthService(mockRepo)

	user, err := authService.Signup(context.Background(), "Test User", "test@example.com", "password123")

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if user.Name != "Test User" {
		t.Errorf("Expected name Test User, got %s", user.Name)
	}
}

func TestAuthService_Login(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// For simplicity, let's just test that it calls the repo.
		mockRepo := &MockAuthRepository{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.UserWithPassword, error) {
				// Re-generating hash here to ensure it's valid
				h, _ := bcrypt.GenerateFromPassword([]byte("password123"), 10)
				return &models.UserWithPassword{
					ID:           1,
					Email:        email,
					PasswordHash: string(h),
				}, nil
			},
		}
		authService := service.NewAuthService(mockRepo)

		user, token, err := authService.Login(context.Background(), "test@example.com", "password123")

		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if user.Email != "test@example.com" {
			t.Errorf("Expected email test@example.com, got %s", user.Email)
		}
		if token == "" {
			t.Error("Expected token, got empty string")
		}
	})

	t.Run("InvalidPassword", func(t *testing.T) {
		mockRepo := &MockAuthRepository{
			GetUserByEmailFunc: func(ctx context.Context, email string) (*models.UserWithPassword, error) {
				h, _ := bcrypt.GenerateFromPassword([]byte("correct-password"), 10)
				return &models.UserWithPassword{
					ID:           1,
					Email:        email,
					PasswordHash: string(h),
				}, nil
			},
		}
		authService := service.NewAuthService(mockRepo)

		_, _, err := authService.Login(context.Background(), "test@example.com", "wrong-password")

		if err == nil || err.Error() != "invalid credentials" {
			t.Errorf("Expected invalid credentials error, got %v", err)
		}
	})
}

// MockTodoRepository implements service.TodoRepository
type MockTodoRepository struct {
	CreateFunc    func(ctx context.Context, userID int, title string, description *string) (*models.Todo, error)
	GetByUserFunc func(ctx context.Context, userID int, limit, offset int) ([]models.Todo, error)
}

func (m *MockTodoRepository) Create(ctx context.Context, userID int, title string, description *string) (*models.Todo, error) {
	return m.CreateFunc(ctx, userID, title, description)
}
func (m *MockTodoRepository) GetByUser(ctx context.Context, userID int, limit, offset int) ([]models.Todo, error) {
	return m.GetByUserFunc(ctx, userID, limit, offset)
}
func (m *MockTodoRepository) GetByID(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	return nil, nil
}
func (m *MockTodoRepository) Update(ctx context.Context, todoID, userID int, title, desc *string, comp *bool) (*models.Todo, error) {
	return nil, nil
}
func (m *MockTodoRepository) Delete(ctx context.Context, todoID, userID int) error { return nil }
func (m *MockTodoRepository) Toggle(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	return nil, nil
}

func TestTodoService_Create(t *testing.T) {
	mockRepo := &MockTodoRepository{
		CreateFunc: func(ctx context.Context, userID int, title string, description *string) (*models.Todo, error) {
			return &models.Todo{ID: 1, UserID: userID, Title: title}, nil
		},
	}
	todoService := service.NewTodoService(mockRepo)

	todo, err := todoService.Create(context.Background(), 1, "Buy Milk", nil)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if todo.Title != "Buy Milk" {
		t.Errorf("Expected title Buy Milk, got %s", todo.Title)
	}
}

func TestTodoService_GetByUser_Defaults(t *testing.T) {
	var capturedLimit int
	mockRepo := &MockTodoRepository{
		GetByUserFunc: func(ctx context.Context, userID int, limit, offset int) ([]models.Todo, error) {
			capturedLimit = limit
			return []models.Todo{}, nil
		},
	}
	todoService := service.NewTodoService(mockRepo)

	todoService.GetByUser(context.Background(), 1, 0, -1)

	if capturedLimit != 100 {
		t.Errorf("Expected limit to default to 100, got %d", capturedLimit)
	}
}
