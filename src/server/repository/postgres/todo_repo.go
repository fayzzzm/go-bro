package postgres

import (
	"context"

	"github.com/fayzzzm/go-bro/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TodoRepo struct {
	pool *pgxpool.Pool
}

func NewTodoRepo(pool *pgxpool.Pool) *TodoRepo {
	return &TodoRepo{pool: pool}
}

func (r *TodoRepo) Create(ctx context.Context, userID int, title string, description *string) (*models.Todo, error) {
	return queryOne[models.Todo](ctx, r.pool,
		"SELECT * FROM todos.create($1, $2, $3)",
		userID, title, description,
	)
}

func (r *TodoRepo) GetByUser(ctx context.Context, userID int, limit, offset int) ([]models.Todo, error) {
	return queryRows[models.Todo](ctx, r.pool,
		"SELECT * FROM todos.list($1, $2, $3)",
		userID, limit, offset,
	)
}

func (r *TodoRepo) GetByID(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	return queryOne[models.Todo](ctx, r.pool,
		"SELECT * FROM todos.get($1, $2)",
		todoID, userID,
	)
}

func (r *TodoRepo) Update(ctx context.Context, todoID, userID int, title *string, description *string, completed *bool) (*models.Todo, error) {
	return queryOne[models.Todo](ctx, r.pool,
		"SELECT * FROM todos.update($1, $2, $3, $4, $5)",
		todoID, userID, title, description, completed,
	)
}

func (r *TodoRepo) Delete(ctx context.Context, todoID, userID int) error {
	_, err := r.pool.Exec(ctx, "SELECT todos.delete($1, $2)", todoID, userID)
	return err
}

func (r *TodoRepo) Toggle(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	return queryOne[models.Todo](ctx, r.pool,
		"SELECT * FROM todos.toggle($1, $2)",
		todoID, userID,
	)
}
