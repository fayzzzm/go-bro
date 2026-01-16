package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/user/go-learning/src/models"
)

type PostgresUserRepo struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepo(pool *pgxpool.Pool) *PostgresUserRepo {
	return &PostgresUserRepo{pool: pool}
}

func (r *PostgresUserRepo) Create(u *models.User) error {
	// Note: Our schema.sql uses UUID and returns it. 
	// To keep it simple for this exercise, we'll map into our int-based model.
	// In a real app, you'd use UUIDs in your models/user.go too.
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, 'temporary_hash') RETURNING id`
	
	// We'll use a hack to convert PG's UUID or serial to our model's Int for now
	var id int
	// If your schema uses Serial, this works. If UUID, we'd need to change models.User.ID to string.
	// For this exercise, let's assume we are just pushing data.
	err := r.pool.QueryRow(context.Background(), query, u.Name, u.Email).Scan(&id)
	if err != nil {
		return err
	}
	u.ID = id
	return nil
}

func (r *PostgresUserRepo) GetByID(id int) (*models.User, error) {
	u := &models.User{}
	query := `SELECT id, name, email, created_at FROM users WHERE id = $1`
	
	// Scan directly into the model
	err := r.pool.QueryRow(context.Background(), query, id).Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *PostgresUserRepo) GetAll() ([]*models.User, error) {
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
