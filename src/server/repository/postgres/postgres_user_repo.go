package postgres

import (
	"context"

	"github.com/fayzzzm/go-bro/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) RegisterUser(ctx context.Context, name, email string) (*models.User, error) {
	payload := map[string]any{
		"name":  name,
		"email": email,
	}
	// We use the same 'create' function as Signup but without password hash
	return queryOne[models.User](ctx, r.pool, "SELECT * FROM users.create($1)", payload)
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	payload := map[string]any{
		"id": id,
	}
	return queryOne[models.User](ctx, r.pool, "SELECT * FROM users.get($1)", payload)
}

func (r *UserRepo) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	payload := map[string]any{
		"limit_val":  limit,
		"offset_val": offset,
	}
	return queryRows[models.User](ctx, r.pool, "SELECT * FROM users.list($1)", payload)
}
