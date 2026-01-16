package service

import (
	"errors"
	"github.com/user/go-learning/src/models"
	"github.com/user/go-learning/src/repository"
	"strings"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
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
