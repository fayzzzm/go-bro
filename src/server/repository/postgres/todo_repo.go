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
	payload := map[string]any{
		"user_id":     userID,
		"title":       title,
		"description": description,
	}
	return queryOne[models.Todo](ctx, r.pool, "SELECT * FROM todos.create($1)", payload)
}

func (r *TodoRepo) GetByUser(ctx context.Context, userID int, limit, offset int) ([]models.Todo, error) {
	payload := map[string]any{
		"user_id": userID,
		"limit":   limit,
		"offset":  offset,
	}
	return queryRows[models.Todo](ctx, r.pool, "SELECT * FROM todos.list($1)", payload)
}

func (r *TodoRepo) GetByID(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	payload := map[string]any{
		"id":      todoID,
		"user_id": userID,
	}
	return queryOne[models.Todo](ctx, r.pool, "SELECT * FROM todos.get($1)", payload)
}

func (r *TodoRepo) Update(ctx context.Context, userID int, todo *models.Todo) (*models.Todo, error) {
	// We pass the whole todo object, but we could also wrap it if needed.
	// Since our SQL expects 'user_id' in the payload, we add it.
	payload := map[string]any{
		"id":          todo.ID,
		"user_id":     userID,
		"title":       todo.Title,
		"description": todo.Description,
		"completed":   todo.Completed,
	}
	return queryOne[models.Todo](ctx, r.pool, "SELECT * FROM todos.update($1)", payload)
}

func (r *TodoRepo) Delete(ctx context.Context, todoID, userID int) error {
	payload := map[string]any{
		"id":      todoID,
		"user_id": userID,
	}
	_, err := r.pool.Exec(ctx, "SELECT todos.delete($1)", payload)
	return err
}

func (r *TodoRepo) Toggle(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	payload := map[string]any{
		"id":      todoID,
		"user_id": userID,
	}
	return queryOne[models.Todo](ctx, r.pool, "SELECT * FROM todos.toggle($1)", payload)
}
