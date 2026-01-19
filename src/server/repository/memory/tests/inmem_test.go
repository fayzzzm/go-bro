package tests

import (
	"context"
	"testing"

	"github.com/fayzzzm/go-bro/repository/memory"
)

func TestInMemUserRepo(t *testing.T) {
	repo := memory.NewUserRepo()
	ctx := context.Background()

	// Test Register
	user, err := repo.RegisterUser(ctx, "InMem User", "inmem@test.com")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}
	if user.ID != 1 {
		t.Errorf("Expected ID 1, got %d", user.ID)
	}

	// Test Get
	found, err := repo.GetByID(ctx, 1)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if found.Email != "inmem@test.com" {
		t.Errorf("Expected email inmem@test.com, got %s", found.Email)
	}

	// Test Duplicate
	_, err = repo.RegisterUser(ctx, "Other", "inmem@test.com")
	if err == nil {
		t.Error("Expected error for duplicate email, got nil")
	}
}
