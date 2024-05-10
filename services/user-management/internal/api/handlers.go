package api

import (
	"net/http"
	"os"
	"time"
	"user-management/internal/initializers"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Signup(c *gin.Context) {
	var body struct {
		ID        int64  `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email" gorm:"unique"`
		Age       int    `json:"age"`
		BirthDate string `json:"birthDate"`
		Role      string `json:"role"`
		Phone     string `json:"phone" gorm:"unique"`
		Street    string `json:"street"`
		City      string `json:"city"`
		Password  string `json:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Faild to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to hash password",
		})
		return
	}

	user := model.User{Email: body.Email, Age: body.Age, BirthDate: body.BirthDate, Role: body.Role, Phone: body.Phone, Street: body.Street, City: body.City, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email" gorm:"unique"`
		Password string `json:"password"`
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Faild to read body",
		})
		return
	}
	var user model.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Invalid email or password",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": "Failed to create token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Token": tokenString,
	})
}
