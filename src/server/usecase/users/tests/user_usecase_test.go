package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/fayzzzm/go-bro/models"
	"github.com/fayzzzm/go-bro/usecase/users"
)

// MockUserService implements service.UserServicer for testing use cases
type MockUserService struct {
	users map[int]*models.User
	err   error
}

func (m *MockUserService) RegisterUser(ctx context.Context, name, email string) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}

	// Simulate service behavior
	if name == "" {
		return nil, errors.New("name cannot be empty")
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

func (m *MockUserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (m *MockUserService) ListUsers(ctx context.Context, limit, offset int) ([]models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	var res []models.User
	for _, u := range m.users {
		res = append(res, *u)
	}
	return res, nil
}

func TestUserUseCase_RegisterUser(t *testing.T) {
	ctx := context.Background()

	t.Run("successful registration returns proper output DTO", func(t *testing.T) {
		mockSvc := &MockUserService{users: make(map[int]*models.User)}
		uc := users.NewUseCase(mockSvc)

		input := users.RegisterUserInput{Name: "John Doe", Email: "john@example.com"}
		output, err := uc.RegisterUser(ctx, input)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !output.Success {
			t.Error("expected Success to be true")
		}
		if output.User == nil {
			t.Error("expected User to not be nil")
		}
		if output.User.Name != "John Doe" {
			t.Errorf("expected name John Doe, got %s", output.User.Name)
		}
		if output.Message != "User registered successfully" {
			t.Errorf("expected success message, got %s", output.Message)
		}
	})

	t.Run("failed registration returns error in output", func(t *testing.T) {
		mockSvc := &MockUserService{users: make(map[int]*models.User)}
		uc := users.NewUseCase(mockSvc)

		input := users.RegisterUserInput{Name: "", Email: "john@example.com"}
		output, err := uc.RegisterUser(ctx, input)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if output.Success {
			t.Error("expected Success to be false")
		}
		if output.User != nil {
			t.Error("expected User to be nil on failure")
		}
	})
}

func TestUserUseCase_GetUser(t *testing.T) {
	ctx := context.Background()

	t.Run("user found returns output with Found=true", func(t *testing.T) {
		mockSvc := &MockUserService{users: map[int]*models.User{
			1: {ID: 1, Name: "Test User", Email: "test@test.com", CreatedAt: time.Now()},
		}}
		uc := users.NewUseCase(mockSvc)

		input := users.GetUserInput{ID: 1}
		output, err := uc.GetUser(ctx, input)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if !output.Found {
			t.Error("expected Found to be true")
		}
		if output.User.ID != 1 {
			t.Errorf("expected ID 1, got %d", output.User.ID)
		}
	})

	t.Run("user not found returns output with Found=false", func(t *testing.T) {
		mockSvc := &MockUserService{users: make(map[int]*models.User)}
		uc := users.NewUseCase(mockSvc)

		input := users.GetUserInput{ID: 999}
		output, err := uc.GetUser(ctx, input)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if output.Found {
			t.Error("expected Found to be false")
		}
	})
}

func TestUserUseCase_ListUsers(t *testing.T) {
	ctx := context.Background()

	t.Run("list users with default pagination", func(t *testing.T) {
		mockSvc := &MockUserService{users: map[int]*models.User{
			1: {ID: 1, Name: "User 1", Email: "user1@test.com", CreatedAt: time.Now()},
			2: {ID: 2, Name: "User 2", Email: "user2@test.com", CreatedAt: time.Now()},
		}}
		uc := users.NewUseCase(mockSvc)

		input := users.ListUsersInput{} // Empty uses defaults
		output, err := uc.ListUsers(ctx, input)

		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if output.Total != 2 {
			t.Errorf("expected Total 2, got %d", output.Total)
		}
		if len(output.Users) != 2 {
			t.Errorf("expected 2 users, got %d", len(output.Users))
		}
	})

	t.Run("applies default limit when zero", func(t *testing.T) {
		mockSvc := &MockUserService{users: make(map[int]*models.User)}
		uc := users.NewUseCase(mockSvc)

		input := users.ListUsersInput{Limit: 0, Offset: -5}
		_, err := uc.ListUsers(ctx, input)

		// Should not error - defaults applied
		if err != nil {
			t.Errorf("expected no error with default limits, got %v", err)
		}
	})
}
