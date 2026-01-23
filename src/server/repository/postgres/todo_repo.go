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
	payload := TodoRequest{
		UserID:      &userID,
		Title:       &title,
		Description: description,
	}
	return queryOne[models.Todo](ctx, r.pool, "SELECT * FROM todos.create($1)", payload)
}

func (r *TodoRepo) GetByUser(ctx context.Context, userID int, limit, offset int) ([]models.Todo, error) {
	payload := TodoRequest{
		UserID:    &userID,
		LimitVal:  &limit,
		OffsetVal: &offset,
	}
	return queryRows[models.Todo](ctx, r.pool, "SELECT * FROM todos.list($1)", payload)
}

func (r *TodoRepo) GetByID(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	payload := TodoRequest{
		ID:     &todoID,
		UserID: &userID,
	}
	return queryOne[models.Todo](ctx, r.pool, "SELECT * FROM todos.get($1)", payload)
}

func (r *TodoRepo) Update(ctx context.Context, userID int, todo *models.Todo) (*models.Todo, error) {
	payload := TodoRequest{
		ID:          &todo.ID,
		UserID:      &userID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
	}
	return queryOne[models.Todo](ctx, r.pool, "SELECT * FROM todos.update($1)", payload)
}

func (r *TodoRepo) Delete(ctx context.Context, todoID, userID int) error {
	payload := TodoRequest{
		ID:     &todoID,
		UserID: &userID,
	}
	_, err := r.pool.Exec(ctx, "SELECT todos.delete($1)", payload)
	return err
}

func (r *TodoRepo) Toggle(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	payload := TodoRequest{
		ID:     &todoID,
		UserID: &userID,
	}
	return queryOne[models.Todo](ctx, r.pool, "SELECT * FROM todos.toggle($1)", payload)
}
