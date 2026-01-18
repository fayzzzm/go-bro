package user

import "context"

// UserUseCase defines the contract for user-related use cases.
// This is the Input Port that controllers will depend on.
type UserUseCase interface {
	RegisterUser(ctx context.Context, input RegisterUserInput) (*RegisterUserOutput, error)
	GetUser(ctx context.Context, input GetUserInput) (*GetUserOutput, error)
	ListUsers(ctx context.Context, input ListUsersInput) (*ListUsersOutput, error)
}
