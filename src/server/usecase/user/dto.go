package user

import "time"

// Input DTOs - What the use case receives
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

// Output DTOs - What the use case returns
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
	Users []*UserOutput `json:"users"`
	Total int           `json:"total"`
}
