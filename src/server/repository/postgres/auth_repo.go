package postgres

import (
	"context"

	"github.com/fayzzzm/go-bro/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepo struct {
	pool *pgxpool.Pool
}

func NewAuthRepo(pool *pgxpool.Pool) *AuthRepo {
	return &AuthRepo{pool: pool}
}

func (r *AuthRepo) Signup(ctx context.Context, name, email, passwordHash string) (*models.User, error) {
	payload := map[string]any{
		"name":          name,
		"email":         email,
		"password_hash": passwordHash,
	}
	return queryOne[models.User](ctx, r.pool, "SELECT * FROM users.create($1)", payload)
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*models.UserWithPassword, error) {
	payload := map[string]any{
		"email": email,
	}
	return queryOne[models.UserWithPassword](ctx, r.pool, "SELECT * FROM users.get_by_email($1)", payload)
}
