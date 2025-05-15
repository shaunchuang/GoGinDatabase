package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Logger is a middleware that logs the incoming requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log the request
		// You can use a logging library or just print to console
		fmt.Printf("Request: %s %s\n", c.Request.Method, c.Request.URL.Path)
		c.Next() // Call the next middleware/handler
	}
}

// Recovery is a middleware that recovers from panics and writes a 500 if there was one
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}()
		c.Next()
	}
}
