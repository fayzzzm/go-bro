package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/fayzzzm/go-bro/src/models"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(u *models.User) error {
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, 'temporary_hash') RETURNING id`
	
	var id int
	err := r.pool.QueryRow(context.Background(), query, u.Name, u.Email).Scan(&id)
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func (r *UserRepo) GetByID(id int) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	
	err := r.pool.QueryRow(context.Background(), query, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetAll() ([]*models.User, error) {
	query := `SELECT id, name, email, created_at FROM users LIMIT 100`
	rows, err := r.pool.Query(context.Background(), query)
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
