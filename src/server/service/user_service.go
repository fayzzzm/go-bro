package service

import (
	"context"

	"github.com/fayzzzm/go-bro/models"
)

// UserRepository defines the contract for user data access.
// In Hexagonal Architecture, this is an "Output Port".
// The repository now calls SQL functions for all operations.
type UserRepository interface {
	// RegisterUser calls the SQL function fn_register_user
	RegisterUser(ctx context.Context, name, email string) (*models.User, error)
	// GetByID calls the SQL function fn_get_user_by_id
	GetByID(ctx context.Context, id int) (*models.User, error)
	// GetAll calls the SQL function fn_list_users
	GetAll(ctx context.Context, limit, offset int) ([]*models.User, error)
}

// UserServicer is the interface that use cases depend on.
// This is the Input Port for the service layer.
type UserServicer interface {
	RegisterUser(ctx context.Context, name, email string) (*models.User, error)
	GetUser(ctx context.Context, id int) (*models.User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error)
}

// UserService is a thin wrapper that delegates to the repository.
// Business logic validation is now handled by SQL functions.
type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(ctx context.Context, name, email string) (*models.User, error) {
	// All validation and business logic is now in the SQL function fn_register_user
	return s.repo.RegisterUser(ctx, name, email)
}

func (s *UserService) GetUser(ctx context.Context, id int) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	return s.repo.GetAll(ctx, limit, offset)
}
