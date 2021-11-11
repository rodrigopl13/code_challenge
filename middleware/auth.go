package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthToken(validateToken func(string) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		value := c.GetHeader("authorization")
		var token string
		if value == "" {
			token = c.Query("bearer")
		} else {
			token = extractToken(value)
		}
		err := validateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid jwt token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func extractToken(value string) string {
	words := strings.Fields(value)
	if len(words) != 2 {
		return ""
	}
	return words[1]
}
