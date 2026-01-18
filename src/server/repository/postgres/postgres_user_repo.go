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

// RegisterUser calls the fn_register_user SQL function.
// The function handles validation (empty name, email format) and returns the created user.
// Errors from the SQL function are returned as-is (e.g., RAISE EXCEPTION becomes a Go error).
func (r *UserRepo) RegisterUser(ctx context.Context, name, email string) (*models.User, error) {
	u := &models.User{}

	// Call the SQL function that handles all validation and insertion
	query := `SELECT * FROM fn_register_user($1, $2)`

	err := r.pool.QueryRow(ctx, query, name, email).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetByID calls the fn_get_user_by_id SQL function.
func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	u := &models.User{}

	query := `SELECT * FROM fn_get_user_by_id($1)`

	err := r.pool.QueryRow(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// GetAll calls the fn_list_users SQL function with pagination.
func (r *UserRepo) GetAll(ctx context.Context, limit, offset int) ([]*models.User, error) {
	query := `SELECT * FROM fn_list_users($1, $2)`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		u := &models.User{}
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, nil
}
