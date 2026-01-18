package memory

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/fayzzzm/go-bro/models"
)

// UserRepo is an in-memory implementation of the service.UserRepository interface.
// It simulates the behavior of SQL functions for testing/development purposes.
type UserRepo struct {
	mu     sync.RWMutex
	users  map[int]*models.User
	nextID int
}

func NewUserRepo() *UserRepo {
	return &UserRepo{
		users:  make(map[int]*models.User),
		nextID: 1,
	}
}

// RegisterUser simulates the fn_register_user SQL function behavior.
func (r *UserRepo) RegisterUser(ctx context.Context, name, email string) (*models.User, error) {
	// Simulate SQL function validation
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("name cannot be empty")
	}
	if !strings.Contains(email, "@") {
		return nil, errors.New("invalid email address")
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate email
	for _, u := range r.users {
		if u.Email == email {
			return nil, errors.New("email already exists")
		}
	}

	user := &models.User{
		ID:        r.nextID,
		Name:      strings.TrimSpace(name),
		Email:     strings.ToLower(strings.TrimSpace(email)),
		CreatedAt: time.Now(),
	}
	r.users[user.ID] = user
	r.nextID++

	return user, nil
}

// GetByID simulates the fn_get_user_by_id SQL function behavior.
func (r *UserRepo) GetByID(ctx context.Context, id int) (*models.User, error) {
	if id <= 0 {
		return nil, errors.New("invalid user id")
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

// GetAll simulates the fn_list_users SQL function behavior.
func (r *UserRepo) GetAll(ctx context.Context, limit, offset int) ([]*models.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Apply defaults
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	// Collect all users and sort by ID (simulating ORDER BY)
	allUsers := make([]*models.User, 0, len(r.users))
	for _, u := range r.users {
		allUsers = append(allUsers, u)
	}

	// Apply offset
	if offset >= len(allUsers) {
		return []*models.User{}, nil
	}
	allUsers = allUsers[offset:]

	// Apply limit
	if limit < len(allUsers) {
		allUsers = allUsers[:limit]
	}

	return allUsers, nil
}
