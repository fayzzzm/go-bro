package service

import (
	"context"
	"errors"

	"github.com/fayzzzm/go-bro/models"
	"github.com/fayzzzm/go-bro/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepository interface {
	Signup(ctx context.Context, name, email, hashedPassword string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserWithPassword, error)
}

type AuthServicer interface {
	Signup(ctx context.Context, name, email, password string) (*models.User, string, error)
	Login(ctx context.Context, email, password string) (*models.User, string, error)
}

type AuthService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Signup(ctx context.Context, name, email, password string) (*models.User, string, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	// Create user
	user, err := s.repo.Signup(ctx, name, email, string(hashedPassword))
	if err != nil {
		return nil, "", err
	}

	// Generate token
	token, err := auth.GenerateToken(user.ID, user.Email)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	// Get user with password
	userWithPassword, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(userWithPassword.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	// Generate token
	token, err := auth.GenerateToken(userWithPassword.ID, userWithPassword.Email)
	if err != nil {
		return nil, "", err
	}

	return &models.User{
		ID:        userWithPassword.ID,
		Name:      userWithPassword.Name,
		Email:     userWithPassword.Email,
		CreatedAt: userWithPassword.CreatedAt,
	}, token, nil
}
