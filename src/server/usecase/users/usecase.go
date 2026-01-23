package users

import (
	"context"
	"time"

	"github.com/fayzzzm/go-bro/service"
)

// --- Interfaces ---

// UseCase defines the contract for user-related use cases.
type UseCase interface {
	RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error)
	GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error)
	ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error)
}

// --- DTOs ---

// Input DTOs
type RegisterUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetUserInput struct {
	ID int `json:"id"`
}

type ListUsersInput struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// Output DTOs
type UserOutput struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterUserOutput struct {
	User    *UserOutput `json:"user"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
}

type GetUserOutput struct {
	User  *UserOutput `json:"user,omitempty"`
	Found bool        `json:"found"`
}

type ListUsersOutput struct {
	Users []UserOutput `json:"users"`
	Total int          `json:"total"`
}

// --- Implementation ---

type UseCaseImpl struct {
	userService service.UserServicer
}

func NewUseCase(svc service.UserServicer) *UseCaseImpl {
	return &UseCaseImpl{userService: svc}
}

func (uc *UseCaseImpl) RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	user, err := uc.userService.RegisterUser(ctx, input.Name, input.Email)
	if err != nil {
		return &RegisterUserOutput{
			Success: false,
			Message: err.Error(),
		}, err
	}

	return &RegisterUserOutput{
		User: &UserOutput{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		Success: true,
		Message: "User registered successfully",
	}, nil
}

func (uc *UseCaseImpl) GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
	user, err := uc.userService.GetUser(ctx, input.ID)
	if err != nil {
		return &GetUserOutput{
			Found: false,
		}, err
	}

	return &GetUserOutput{
		User: &UserOutput{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
		Found: true,
	}, nil
}

func (uc *UseCaseImpl) ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 100
	}
	offset := input.Offset
	if offset < 0 {
		offset = 0
	}

	users, err := uc.userService.ListUsers(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	output := &ListUsersOutput{
		Users: make([]UserOutput, len(users)),
		Total: len(users),
	}

	for i, u := range users {
		output.Users[i] = UserOutput{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		}
	}

	return output, nil
}
