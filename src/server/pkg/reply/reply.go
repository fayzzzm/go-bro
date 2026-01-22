package reply

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error sends a JSON error response and returns true if an error exists.
func Error(c *gin.Context, code int, message string, err error) bool {
	if err != nil {
		c.JSON(code, gin.H{"error": message})
		return true
	}
	return false
}

// NotFound is a specialized version of Error for 404 responses.
func NotFound(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
		return true
	}
	return false
}

// InternalError is a specialized version of Error for 500 responses.
func InternalError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return true
	}
	return false
}

// OK sends a 200 OK response.
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// Created sends a 201 Created response.
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}
