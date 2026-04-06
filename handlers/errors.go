package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SanitizeError returns a safe error message for client consumption
// while logging the full error internally
func SanitizeError(c *gin.Context, fullError error, publicMessage string) {
	// Log the full error for internal debugging
	log.Printf("Error [%s %s]: %v", c.Request.Method, c.Request.URL.Path, fullError)

	// Return safe message to client
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": publicMessage,
	})
}

// SanitizeDBError returns a database-appropriate error message
func SanitizeDBError(c *gin.Context, fullError error) {
	SanitizeError(c, fullError, "database operation failed")
}

// SanitizeCacheError returns a cache-appropriate error message
func SanitizeCacheError(c *gin.Context, fullError error) {
	SanitizeError(c, fullError, "cache operation failed")
}
