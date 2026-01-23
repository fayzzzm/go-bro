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
	payload := UserRequest{
		Name:  &name,
		Email: &email,
	}
	return queryOne[models.User](ctx, r.pool, "SELECT * FROM users.create($1)", payload)
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	payload := UserRequest{
		ID: &id,
	}
	return queryOne[models.User](ctx, r.pool, "SELECT * FROM users.get($1)", payload)
}

func (r *UserRepo) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	payload := UserRequest{
		LimitVal:  &limit,
		OffsetVal: &offset,
	}
	return queryRows[models.User](ctx, r.pool, "SELECT * FROM users.list($1)", payload)
}
