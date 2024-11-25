package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token in Authorization header
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		// Bearer token format: "Bearer <token>"
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader { // No "Bearer " prefix
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization format must be Bearer <token>"})
			c.Abort()
			return
		}

		// Parse and validate the token
		userID, err := ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Attach userID to the context
		c.Set("userID", userID)

		// Continue to the next handler
		c.Next()
	}
}

// RequestLoggerMiddleware logs the request and response details
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Log request details
		start := time.Now()
		Logger.Infof("Started %s %s", c.Request.Method, c.Request.URL.Path)

		// Process the request
		c.Next()

		// Log response details
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		Logger.Infof("Completed %s %s with status %d in %v", c.Request.Method, c.Request.URL.Path, statusCode, duration)
	}
}
