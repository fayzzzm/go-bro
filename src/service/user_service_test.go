package service

import (
	"errors"
	"testing"

	"github.com/fayzzzm/go-bro/src/models"
)

// MockUserRepo is a manual mock implementation of UserRepository
type MockUserRepo struct {
	users map[int]*models.User
	err   error
}

func (m *MockUserRepo) Create(u *models.User) error {
	if m.err != nil {
		return m.err
	}
	u.ID = len(m.users) + 1
	m.users[u.ID] = u
	return nil
}

func (m *MockUserRepo) GetByID(id int) (*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, ok := m.users[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return user, nil
}

func (m *MockUserRepo) GetAll() ([]*models.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	var res []*models.User
	for _, u := range m.users {
		res = append(res, u)
	}
	return res, nil
}

func TestUserService_RegisterUser(t *testing.T) {
	t.Run("successful registration", func(t *testing.T) {
		mockRepo := &MockUserRepo{users: make(map[int]*models.User)}
		svc := NewUserService(mockRepo)

		user, err := svc.RegisterUser("John Doe", "john@example.com")

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

	t.Run("invalid email", func(t *testing.T) {
		mockRepo := &MockUserRepo{users: make(map[int]*models.User)}
		svc := NewUserService(mockRepo)

		_, err := svc.RegisterUser("John Doe", "invalid-email")

		if err == nil {
			t.Error("expected error for invalid email, got nil")
		}
		if err.Error() != "invalid email address" {
			t.Errorf("expected 'invalid email address' error, got %v", err)
		}
	})

	t.Run("empty name", func(t *testing.T) {
		mockRepo := &MockUserRepo{users: make(map[int]*models.User)}
		svc := NewUserService(mockRepo)

		_, err := svc.RegisterUser("", "john@example.com")

		if err == nil {
			t.Error("expected error for empty name, got nil")
		}
	})
}
