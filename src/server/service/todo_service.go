package service

import (
	"context"

	"github.com/fayzzzm/go-bro/models"
)

type TodoRepository interface {
	Create(ctx context.Context, userID int, title string, description *string) (*models.Todo, error)
	GetByUser(ctx context.Context, userID int, limit, offset int) ([]models.Todo, error)
	GetByID(ctx context.Context, todoID, userID int) (*models.Todo, error)
	Update(ctx context.Context, userID int, todo *models.Todo) (*models.Todo, error)
	Delete(ctx context.Context, todoID, userID int) error
	Toggle(ctx context.Context, todoID, userID int) (*models.Todo, error)
}

type TodoServicer interface {
	Create(ctx context.Context, userID int, title string, description *string) (*models.Todo, error)
	GetByUser(ctx context.Context, userID int, limit, offset int) ([]models.Todo, error)
	GetByID(ctx context.Context, todoID, userID int) (*models.Todo, error)
	Update(ctx context.Context, userID int, todo *models.Todo) (*models.Todo, error)
	Delete(ctx context.Context, todoID, userID int) error
	Toggle(ctx context.Context, todoID, userID int) (*models.Todo, error)
}

type TodoService struct {
	repo TodoRepository
}

func NewTodoService(repo TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) Create(ctx context.Context, userID int, title string, description *string) (*models.Todo, error) {
	return s.repo.Create(ctx, userID, title, description)
}

func (s *TodoService) GetByUser(ctx context.Context, userID int, limit, offset int) ([]models.Todo, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return s.repo.GetByUser(ctx, userID, limit, offset)
}

func (s *TodoService) GetByID(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	return s.repo.GetByID(ctx, todoID, userID)
}

func (s *TodoService) Update(ctx context.Context, userID int, todo *models.Todo) (*models.Todo, error) {
	return s.repo.Update(ctx, userID, todo)
}

func (s *TodoService) Delete(ctx context.Context, todoID, userID int) error {
	return s.repo.Delete(ctx, todoID, userID)
}

func (s *TodoService) Toggle(ctx context.Context, todoID, userID int) (*models.Todo, error) {
	return s.repo.Toggle(ctx, todoID, userID)
}
