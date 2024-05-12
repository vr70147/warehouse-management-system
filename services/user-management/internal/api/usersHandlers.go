package api

import (
	"net/http"
	"os"
	"strconv"
	"time"
	"user-management/internal/initializers"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
			"error": "Faild to read body",
		})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	user := model.User{Email: body.Email, Name: body.Name, Age: body.Age, BirthDate: body.BirthDate, Role: body.Role, Phone: body.Phone, Street: body.Street, City: body.City, Password: string(hash)}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create user",
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
			"error": "Faild to read body",
		})
		return
	}
	var user model.User
	initializers.DB.First(&user, "email = ?", body.Email)

	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid email or password",
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
			"error": "Failed to create token",
		})
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
	c.JSON(http.StatusOK, gin.H{})
}
func Validate(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"Message": user,
	})
}

func Logout(c *gin.Context) {
	c.SetCookie("Authorization", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

func UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var user model.User

	result := initializers.DB.First(&user, userID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}
	var updateUser model.User
	if err := c.BindJSON(&updateUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Faild to bind user data",
		})
		return
	}

	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}
	if updateUser.Phone != "" {
		user.Phone = updateUser.Phone
	}
	if updateUser.Street != "" {
		user.Street = updateUser.Street
	}
	if updateUser.City != "" {
		user.City = updateUser.City
	}
	if updateUser.Age != 0 {
		user.Age = updateUser.Age
	}
	if updateUser.Role != "" {
		user.Role = updateUser.Role
	}
	if updateUser.BirthDate != "" {
		user.BirthDate = updateUser.BirthDate
	}
	if updateUser.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Failed to hash password",
			})
			return
		}
		user.BirthDate = string(hash)
	}
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User updated successfuly",
		"user":    user,
	})
}

func GetUsers(c *gin.Context) {
	queryCondition := model.User{}

	var queryExists bool

	if name := c.Query("name"); name != "" {
		queryCondition.Name = name
		queryExists = true
	}
	if email := c.Query("email"); email != "" {
		queryCondition.Email = email
		queryExists = true
	}
	if age := c.Query("age"); age != "" {
		i, _ := strconv.Atoi(age)
		queryCondition.Age = i
		queryExists = true
	}
	if phone := c.Query("phone"); phone != "" {
		queryCondition.Phone = phone
		queryExists = true
	}
	if role := c.Query("role"); role != "" {
		queryCondition.Role = role
		queryExists = true
	}

	var users []model.User
	var result *gorm.DB

	if queryExists || len(c.Params) == 0 {
		result = initializers.DB.Where(&queryCondition).Find(&users)
	} else {
		result = initializers.DB.Find(&users)
	}

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No users found matching the criteria",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	result := initializers.DB.Delete(&model.User{}, userID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No user found with the given ID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfuly"})
}

func RecoverUser(c *gin.Context) {
	userID := c.Param("id")

	result := initializers.DB.Model(&model.User{}).Unscoped().Where("id = ?", userID).Update("deleted_at", nil)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to recover user",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No deleted user found with the given ID",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User recovered successfully"})
}