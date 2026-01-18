package user

import (
	"context"

	"github.com/fayzzzm/go-bro/service"
)

// UserUseCaseImpl implements the UserUseCase interface.
// It orchestrates the flow and transforms between DTOs and domain models.
type UserUseCaseImpl struct {
	userService service.UserServicer
}

func NewUserUseCase(svc service.UserServicer) *UserUseCaseImpl {
	return &UserUseCaseImpl{userService: svc}
}

func (uc *UserUseCaseImpl) RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error) {
	// Call the service (which now just delegates to SQL functions)
	user, err := uc.userService.RegisterUser(ctx, input.Name, input.Email)
	if err != nil {
		return &RegisterUserOutput{
			Success: false,
			Message: err.Error(),
		}, err
	}

	// Transform domain model to output DTO
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

func (uc *UserUseCaseImpl) GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
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

func (uc *UserUseCaseImpl) ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error) {
	// Set defaults
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

	// Transform domain models to output DTOs
	output := &ListUsersOutput{
		Users: make([]*UserOutput, len(users)),
		Total: len(users),
	}

	for i, u := range users {
		output.Users[i] = &UserOutput{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		}
	}

	return output, nil
}
