package repository

import (
	"errors"
	"github.com/user/go-learning/src/models"
	"sync"
	"time"
)

// UserRepository defines the contract for user data access
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id int) (*models.User, error)
	GetAll() ([]*models.User, error)
}

// InMemUserRepo is an in-memory implementation of UserRepository
type InMemUserRepo struct {
	mu    sync.RWMutex
	users map[int]*models.User
	nextID int
}

func NewInMemUserRepo() *InMemUserRepo {
	return &InMemUserRepo{
		users:  make(map[int]*models.User),
		nextID: 1,
	}
}

func (r *InMemUserRepo) Create(u *models.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	u.ID = r.nextID
	u.CreatedAt = time.Now()
	r.users[u.ID] = u
	r.nextID++
	return nil
}

func (r *InMemUserRepo) GetByID(id int) (*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *InMemUserRepo) GetAll() ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*models.User, 0, len(r.users))
	for _, u := range r.users {
		users = append(users, u)
	}
	return users, nil
}
