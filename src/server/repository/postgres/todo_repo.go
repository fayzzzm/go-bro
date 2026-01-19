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
	var todo models.Todo
	err := r.pool.QueryRow(ctx,
		"SELECT id, user_id, title, description, completed, created_at, updated_at FROM fn_create_todo($1, $2, $3)",
		userID, title, description,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepo) GetByUser(ctx context.Context, userID int, limit, offset int) ([]models.Todo, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, user_id, title, description, completed, created_at, updated_at FROM fn_get_todos_by_user($1, $2, $3)",
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		err := rows.Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

func (r *TodoRepo) GetByID(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	var todo models.Todo
	err := r.pool.QueryRow(ctx,
		"SELECT id, user_id, title, description, completed, created_at, updated_at FROM fn_get_todo_by_id($1, $2)",
		todoID, userID,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepo) Update(ctx context.Context, todoID, userID int, title *string, description *string, completed *bool) (*models.Todo, error) {
	var todo models.Todo
	err := r.pool.QueryRow(ctx,
		"SELECT id, user_id, title, description, completed, created_at, updated_at FROM fn_update_todo($1, $2, $3, $4, $5)",
		todoID, userID, title, description, completed,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *TodoRepo) Delete(ctx context.Context, todoID, userID int) error {
	_, err := r.pool.Exec(ctx, "SELECT fn_delete_todo($1, $2)", todoID, userID)
	return err
}

func (r *TodoRepo) Toggle(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	var todo models.Todo
	err := r.pool.QueryRow(ctx,
		"SELECT id, user_id, title, description, completed, created_at, updated_at FROM fn_toggle_todo($1, $2)",
		todoID, userID,
	).Scan(&todo.ID, &todo.UserID, &todo.Title, &todo.Description, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return &todo, nil
}
