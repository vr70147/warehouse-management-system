package utils

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// ExtractToken retrieves the JWT token from the request context.
func ExtractToken(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		return "", fmt.Errorf("authorization header is required")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	if tokenString == "" {
		return "", fmt.Errorf("invalid token format")
	}

	return tokenString, nil
}
