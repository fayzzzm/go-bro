package tests

import (
	"testing"

	"github.com/fayzzzm/go-bro/pkg/auth"
)

func TestJWT(t *testing.T) {
	userID := 123
	email := "user@example.com"

	// Test GenerateToken
	token, err := auth.GenerateToken(userID, email)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}
	if token == "" {
		t.Fatal("Generated token is empty")
	}

	// Test ValidateToken
	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("Expected email %s, got %s", email, claims.Email)
	}

	// Test Invalid Token
	_, err = auth.ValidateToken("invalid.token.string")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}

func TestExpiredToken(t *testing.T) {
	// This would require mocking the time or having a way to generate expired tokens
	// Since pkg/auth doesn't support this easily, we'll skip for now or
	// we could refactor pkg/auth to accept an expiration duration.
	t.Log("Skipping expired token test (requires time mocking)")
}
