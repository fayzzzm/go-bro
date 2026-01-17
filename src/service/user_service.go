package service

import (
	"errors"
	"github.com/fayzzzm/go-bro/src/models"
	"strings"
)

// UserRepository defines the contract for user data access.
// In Hexagonal Architecture, this is an "Output Port".
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetAll() ([]*models.User, error)
}

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(name, email string) (*models.User, error) {
	// Business Logic: Validation
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email address")
	}

	user := &models.User{
		Name:  name,
		Email: email,
	}

	err := s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(id int) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) ListUsers() ([]*models.User, error) {
	return s.repo.GetAll()
}
