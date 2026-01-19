package postgres

import (
	"context"
	"errors"

	"github.com/fayzzzm/go-bro/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{pool: pool}
}

// UserWithPassword includes password hash for login verification
type UserWithPassword struct {
	models.User
	PasswordHash string
}

func (r *AuthRepo) Signup(ctx context.Context, name, email, passwordHash string) (*models.User, error) {
	var user models.User
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, email, created_at FROM fn_signup_user($1, $2, $3)",
		name, email, passwordHash,
	).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt)

	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*UserWithPassword, error) {
	var user UserWithPassword
	err := r.pool.QueryRow(ctx,
		"SELECT id, name, email, password_hash, created_at FROM fn_get_user_by_email($1)",
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}
