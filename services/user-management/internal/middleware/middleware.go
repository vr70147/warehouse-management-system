package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

// RequireAuth checks if the request contains a valid JWT token
func RequireAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Authorization header is required",
			})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Invalid token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Invalid token claims",
			})
			c.Abort()
			return
		}

		accountID, ok := claims["account_id"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Account ID not found in token",
			})
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "User ID not found in token",
			})
			c.Abort()
			return
		}

		// Retrieve user from database
		var user model.User
		if err := db.First(&user, uint(userID)).Error; err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "User not found",
			})
			c.Abort()
			return
		}

		// Set accountID and user to the context
		c.Set("account_id", accountID)
		c.Set("user", user)

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "User not found in context",
			})
			c.Abort()
			return
		}

		u, ok := user.(model.User)
		if !ok || !u.IsAdmin {
			c.JSON(http.StatusForbidden, model.ErrorResponse{
				Error: "Admin privileges required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
