package models

import "time"

type Todo struct {
	ID          int       `json:"id" db:"id"`
	UserID      int       `json:"user_id" db:"user_id"`
	Title       *string   `json:"title,omitempty" db:"title"`
	Description *string   `json:"description,omitempty" db:"description"`
	Completed   *bool     `json:"completed,omitempty" db:"completed"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
