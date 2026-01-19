package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fayzzzm/go-bro/middleware"
	"github.com/fayzzzm/go-bro/pkg/auth"
	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Helper to generate a valid token
	validToken, _ := auth.GenerateToken(1, "test@example.com")

	tests := []struct {
		name           string
		setupRequest   func(req *http.Request)
		expectedStatus int
		expectedEmail  string
	}{
		{
			name: "Valid token in Cookie",
			setupRequest: func(req *http.Request) {
				req.AddCookie(&http.Cookie{
					Name:  middleware.AuthCookieName,
					Value: validToken,
				})
			},
			expectedStatus: http.StatusOK,
			expectedEmail:  "test@example.com",
		},
		{
			name: "Valid token in Header",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", validToken))
			},
			expectedStatus: http.StatusOK,
			expectedEmail:  "test@example.com",
		},
		{
			name:           "Missing token",
			setupRequest:   func(req *http.Request) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Invalid token",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid-token")
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "Empty Bearer token",
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer ")
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := gin.New()
			r.Use(middleware.AuthMiddleware())
			r.GET("/test", func(c *gin.Context) {
				email, _ := middleware.GetUserEmail(c)
				c.JSON(http.StatusOK, gin.H{"email": email})
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			tc.setupRequest(req)

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				var resp map[string]string
				importJson(w.Body.Bytes(), &resp)
				if resp["email"] != tc.expectedEmail {
					t.Errorf("Expected email %s, got %s", tc.expectedEmail, resp["email"])
				}
			}
		})
	}
}

// Simple helper since encoding/json is common
func importJson(data []byte, v interface{}) {
	json.Unmarshal(data, v)
}
