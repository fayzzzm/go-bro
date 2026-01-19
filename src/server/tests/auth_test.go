package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fayzzzm/go-bro/controller"
	"github.com/fayzzzm/go-bro/middleware"
	"github.com/fayzzzm/go-bro/routes"
	"github.com/gin-gonic/gin"
)

// MockAuthService is a mock implementation of AuthServicer for testing
type MockAuthService struct{}

func (m *MockAuthService) Signup(ctx interface{}, name, email, password string) (interface{}, string, error) {
	return map[string]interface{}{
		"id":    1,
		"name":  name,
		"email": email,
	}, "mock-token", nil
}

func (m *MockAuthService) Login(ctx interface{}, email, password string) (interface{}, string, error) {
	return map[string]interface{}{
		"id":    1,
		"name":  "Test User",
		"email": email,
	}, "mock-token", nil
}

// MockTodoService is a mock implementation of TodoServicer
type MockTodoService struct{}

// SetupTestRouter creates a test router with mocked services
func SetupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestHealthEndpoint(t *testing.T) {
	router := SetupTestRouter()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestSignupValidation(t *testing.T) {
	router := SetupTestRouter()

	// Test cases for signup validation
	testCases := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "Missing name",
			body:       map[string]string{"email": "test@test.com", "password": "password123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Missing email",
			body:       map[string]string{"name": "Test", "password": "password123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid email",
			body:       map[string]string{"name": "Test", "email": "invalid", "password": "password123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Short password",
			body:       map[string]string{"name": "Test", "email": "test@test.com", "password": "123"},
			wantStatus: http.StatusBadRequest,
		},
	}

	// Add a simple validation endpoint for testing
	router.POST("/api/v1/auth/signup", func(c *gin.Context) {
		var req controller.SignupRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "ok"})
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/api/v1/auth/signup", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestLoginValidation(t *testing.T) {
	router := SetupTestRouter()

	router.POST("/api/v1/auth/login", func(c *gin.Context) {
		var req controller.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	testCases := []struct {
		name       string
		body       map[string]string
		wantStatus int
	}{
		{
			name:       "Missing email",
			body:       map[string]string{"password": "password123"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Missing password",
			body:       map[string]string{"email": "test@test.com"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Valid credentials format",
			body:       map[string]string{"email": "test@test.com", "password": "password123"},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestProtectedRouteWithoutAuth(t *testing.T) {
	router := SetupTestRouter()

	// Protected route
	router.GET("/api/v1/todos", middleware.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"todos": []interface{}{}})
	})

	req, _ := http.NewRequest("GET", "/api/v1/todos", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for unauthenticated request, got %d", w.Code)
	}
}

// Placeholder for unused imports
var _ = routes.SetupRoutes
