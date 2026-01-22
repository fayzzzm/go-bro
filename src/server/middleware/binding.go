package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const BodyContextKey = "request_body"

// BindJSON is a generic middleware that binds the request body to a struct of type T.
// If binding fails, it aborts the request with a 400 Bad Request error.
func BindJSON[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input T
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			c.Abort()
			return
		}
		c.Set(BodyContextKey, input)
		c.Next()
	}
}

// GetBody retrieves the bound body from the context and casts it to type T.
// It will panic if the body is not present or is of a different type,
// so it should only be used in handlers where BindJSON[T] was applied.
func GetBody[T any](c *gin.Context) T {
	val, exists := c.Get(BodyContextKey)
	if !exists {
		panic("GetBody called on handler without BindJSON middleware")
	}
	return val.(T)
}
