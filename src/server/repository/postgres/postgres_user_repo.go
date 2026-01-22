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
	return queryOne[models.User](ctx, r.pool,
		"SELECT * FROM fn_register_user($1, $2)",
		name, email,
	)
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	return queryOne[models.User](ctx, r.pool,
		"SELECT * FROM fn_get_user_by_id($1)",
		id,
	)
}

func (r *UserRepo) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	return queryRows[models.User](ctx, r.pool,
		"SELECT * FROM fn_list_users($1, $2)",
		limit, offset,
	)
}
