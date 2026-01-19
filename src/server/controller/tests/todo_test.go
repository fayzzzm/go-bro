package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fayzzzm/go-bro/controller"
	"github.com/fayzzzm/go-bro/middleware"
	"github.com/gin-gonic/gin"
)

func TestCreateTodoValidation(t *testing.T) {
	router := SetupTestRouter()

	// Mock auth middleware that always authenticates
	mockAuth := func(c *gin.Context) {
		c.Set(middleware.AuthUserIDKey, 1)
		c.Set(middleware.AuthUserEmailKey, "test@test.com")
		c.Next()
	}

	router.POST("/api/v1/todos", mockAuth, func(c *gin.Context) {
		var req controller.CreateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"todo": gin.H{"id": 1, "title": req.Title}})
	})

	testCases := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name:       "Missing title",
			body:       map[string]interface{}{"description": "Some description"},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Empty title",
			body:       map[string]interface{}{"title": ""},
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Valid todo",
			body:       map[string]interface{}{"title": "Test Todo", "description": "Description"},
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Valid todo without description",
			body:       map[string]interface{}{"title": "Test Todo"},
			wantStatus: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("POST", "/api/v1/todos", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestUpdateTodoValidation(t *testing.T) {
	router := SetupTestRouter()

	mockAuth := func(c *gin.Context) {
		c.Set(middleware.AuthUserIDKey, 1)
		c.Set(middleware.AuthUserEmailKey, "test@test.com")
		c.Next()
	}

	router.PUT("/api/v1/todos/:id", mockAuth, func(c *gin.Context) {
		var req controller.UpdateTodoRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"todo": gin.H{"id": 1, "title": "Updated"}})
	})

	testCases := []struct {
		name       string
		body       map[string]interface{}
		wantStatus int
	}{
		{
			name:       "Update title only",
			body:       map[string]interface{}{"title": "New Title"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Update completed only",
			body:       map[string]interface{}{"completed": true},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Update description only",
			body:       map[string]interface{}{"description": "New description"},
			wantStatus: http.StatusOK,
		},
		{
			name:       "Empty body (valid - no updates)",
			body:       map[string]interface{}{},
			wantStatus: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.body)
			req, _ := http.NewRequest("PUT", "/api/v1/todos/1", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.wantStatus {
				t.Errorf("Expected status %d, got %d. Body: %s", tc.wantStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestTodoRoutesRequireAuth(t *testing.T) {
	router := SetupTestRouter()

	// Routes without auth middleware
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/todos", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"todos": []interface{}{}})
		})
		protected.POST("/todos", func(c *gin.Context) {
			c.JSON(http.StatusCreated, gin.H{"todo": gin.H{}})
		})
		protected.DELETE("/todos/:id", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "deleted"})
		})
	}

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/todos"},
		{"POST", "/api/v1/todos"},
		{"DELETE", "/api/v1/todos/1"},
	}

	for _, ep := range endpoints {
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			req, _ := http.NewRequest(ep.method, ep.path, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("Expected 401 for unauthenticated %s %s, got %d", ep.method, ep.path, w.Code)
			}
		})
	}
}
