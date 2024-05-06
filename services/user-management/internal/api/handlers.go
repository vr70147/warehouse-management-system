package api

import (
	"net/http"
	"user-management/internal/model"
	"user-management/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(conn *pgx.Conn) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user model.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid request body",
			})
			return
		}

		if user.Email == "" || user.Password == "" || user.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Name, email and password are required",
			})
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Faild to hash the password",
			})
			return
		}

		user.Password = string(hashedPassword)

		if err := service.CreateUser(conn, user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Faild to create user",
			})
			c.JSON(http.StatusOK, gin.H{})
		}
	}
}
