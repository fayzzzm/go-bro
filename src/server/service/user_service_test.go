package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fayzzzm/go-bro/models"
)

// MockUserRepo is a manual mock implementation of UserRepository.
// Since business logic is now in SQL functions, this mock simulates
// what the SQL functions would return (including validation errors).
type MockUserRepo struct {
	users map[int]*models.User
	err   error
}

func (m *MockUserRepo) RegisterUser(ctx context.Context, name, email string) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	// Simulate SQL function validation
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if email == "" || !contains(email, "@") {
		return nil, errors.New("invalid email address")
	}

	user := &models.User{
		ID:        len(m.users) + 1,
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
	}
	m.users[user.ID] = user
	return user, nil
}

func (m *MockUserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserRepo) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	var res []models.User
	for _, u := range m.users {
		res = append(res, *u)
	}
	return res, nil
}

// Helper function since we can't import strings in a test mock
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestUserService_RegisterUser(t *testing.T) {
	ctx := context.Background()

	t.Run("successful registration", func(t *testing.T) {
		mockRepo := &MockUserRepo{users: make(map[int]*models.User)}
		svc := NewUserService(mockRepo)

		user, err := svc.RegisterUser(ctx, "John Doe", "john@example.com")

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user.Name != "John Doe" {
			t.Errorf("expected name John Doe, got %s", user.Name)
		}
		if len(mockRepo.users) != 1 {
			t.Errorf("expected 1 user in repo, got %d", len(mockRepo.users))
		}
	})

	t.Run("invalid email - validation by SQL function", func(t *testing.T) {
		mockRepo := &MockUserRepo{users: make(map[int]*models.User)}
		svc := NewUserService(mockRepo)

		_, err := svc.RegisterUser(ctx, "John Doe", "invalid-email")

		if err == nil {
			t.Error("expected error for invalid email, got nil")
		}
		if err.Error() != "invalid email address" {
			t.Errorf("expected 'invalid email address' error, got %v", err)
		}
	})

	t.Run("empty name - validation by SQL function", func(t *testing.T) {
		mockRepo := &MockUserRepo{users: make(map[int]*models.User)}
		svc := NewUserService(mockRepo)

		_, err := svc.RegisterUser(ctx, "", "john@example.com")

		if err == nil {
			t.Error("expected error for empty name, got nil")
		}
		if err.Error() != "name cannot be empty" {
			t.Errorf("expected 'name cannot be empty' error, got %v", err)
		}
	})
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.Background()

	t.Run("user found", func(t *testing.T) {
		mockRepo := &MockUserRepo{users: map[int]*models.User{
			1: {ID: 1, Name: "Test", Email: "test@test.com"},
		}}
		svc := NewUserService(mockRepo)

		user, err := svc.GetUser(ctx, 1)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if user.ID != 1 {
			t.Errorf("expected ID 1, got %d", user.ID)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := &MockUserRepo{users: make(map[int]*models.User)}
		svc := NewUserService(mockRepo)

		_, err := svc.GetUser(ctx, 999)

		if err == nil {
			t.Error("expected error for non-existent user, got nil")
		}
	})
}

func TestUserService_ListUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("list all users", func(t *testing.T) {
		mockRepo := &MockUserRepo{users: map[int]*models.User{
			1: {ID: 1, Name: "User1", Email: "user1@test.com"},
			2: {ID: 2, Name: "User2", Email: "user2@test.com"},
		}}
		svc := NewUserService(mockRepo)

		users, err := svc.ListUsers(ctx, 100, 0)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if len(users) != 2 {
			t.Errorf("expected 2 users, got %d", len(users))
		}
	})
}
