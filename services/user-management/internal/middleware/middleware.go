package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

func RequireAuth(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			log.Println("Authorization header is missing")
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
			log.Println("Invalid token:", err)
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Invalid token",
			})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Invalid token claims")
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Invalid token claims",
			})
			c.Abort()
			return
		}

		accountID, ok := claims["account_id"].(string)
		if !ok {
			log.Println("Account ID not found in token")
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Account ID not found in token",
			})
			c.Abort()
			return
		}

		userID, ok := claims["sub"].(float64)
		if !ok {
			log.Println("User ID not found in token")
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "User ID not found in token",
			})
			c.Abort()
			return
		}

		accountIDUint, err := strconv.ParseUint(accountID, 10, 32)
		if err != nil {
			log.Println("Invalid Account ID format:", err)
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "Invalid Account ID format",
			})
			c.Abort()
			return
		}

		// Retrieve user from database
		var user model.User

		if err := db.First(&user, uint(userID)).Error; err != nil {
			log.Println("User not found:", err)
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{
				Error: "User not found",
			})
			c.Abort()
			return
		}

		// Set accountID and user to the context
		c.Set("account_id", uint(accountIDUint))
		c.Set("user", user)

		c.Next()
	}
}
