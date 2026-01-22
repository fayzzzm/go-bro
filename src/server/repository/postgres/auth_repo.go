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

func (r *AuthRepo) Signup(ctx context.Context, name, email, passwordHash string) (*models.User, error) {
	return queryOne[models.User](ctx, r.pool,
		"SELECT * FROM users.create($1, $2, $3)",
		name, email, passwordHash,
	)
}

func (r *AuthRepo) GetUserByEmail(ctx context.Context, email string) (*models.UserWithPassword, error) {
	user, err := queryOne[models.UserWithPassword](ctx, r.pool,
		"SELECT * FROM users.get_by_email($1)",
		email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}
