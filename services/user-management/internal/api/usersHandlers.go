package api

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
	"user-management/internal/initializers"

	// "user-management/internal/kafka"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Signup godoc
// @Summary Signup a new user
// @Description Register a new user with personal details and role
// @Tags users
// @Accept json
// @Produce json
// @Param body body model.User true "User data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /signup [post]
func Signup(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			PersonalID string `json:"personal_id" gorm:"unique;not null"`
			Name       string `json:"name" gorm:"unique;not null"`
			Email      string `json:"email" gorm:"unique;not null"`
			Age        int    `json:"age" gorm:"not null"`
			BirthDate  string `json:"birthDate" gorm:"not null"`
			RoleID     uint   `json:"role_id" gorm:"not null"`
			Role       string `json:"role" gorm:"foreignKey:RoleID"`
			Phone      string `json:"phone" gorm:"unique; not null"`
			Street     string `json:"street"`
			City       string `json:"city"`
			Password   string `json:"password" gorm:"not null"`
			IsAdmin    bool   `json:"is_admin" gorm:"default: false"`
		}

		if err := c.Bind(&body); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to read body",
			})
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to hash password",
			})
			return
		}

		var role model.Role
		fmt.Println(body)
		if result := initializers.DB.Where("id = ?", body.RoleID).First(&role); result.Error != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Role not found",
			})
			return
		}

		user := model.User{
			PersonalID: body.PersonalID,
			Email:      body.Email,
			Name:       body.Name,
			Age:        body.Age,
			BirthDate:  body.BirthDate,
			RoleID:     role.ID,
			Phone:      body.Phone,
			Street:     body.Street,
			City:       body.City,
			Password:   string(hash),
			IsAdmin:    body.IsAdmin,
		}
		if result := initializers.DB.Create(&user); result.Error != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Failed to create user: " + result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "User registered successfully"})
	}

}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return a JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param body body model.User true "Login credentials"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /login [post]
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {

		var body struct {
			Email    string `json:"email" gorm:"unique"`
			Password string `json:"password"`
		}

		if c.Bind(&body) != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Faild to read body",
			})
			return
		}
		var user model.User
		initializers.DB.First(&user, "email = ?", body.Email)

		if user.ID == 0 {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid email or password",
			})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid email or password",
			})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.ID,
			"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		// Sign and get the complete encoded token as a string using the secret
		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		fmt.Println(tokenString)

		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to create token",
			})
			return
		}
		// event := model.UserEvent{
		// 	EventType: "UserLoggedIn",
		// 	User:      user,
		// }
		// kafka.ProducerUserEvent(event)

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "User authenticated successfully",
		})
	}
}

// CheckUserAuthentication godoc
// @Summary Check user authentication
// @Description Verify if the user is authenticated
// @Tags users
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /check-auth [get]
func CheckUserAuthentication(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}

// Logout godoc
// @Summary Logout user
// @Description Logout user and clear the authentication cookie
// @Tags users
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Router /logout [post]
func Logout(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.SetCookie("Authorization", "", -1, "/", "", false, true)

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Successfully logged out",
		})
	}
}

// UpdateUser godoc
// @Summary Update user details
// @Description Update user information by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body model.User true "User data"
// @Success 200 {object} model.UpdateUserResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/{id} [put]
func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")
		var user model.User

		result := initializers.DB.First(&user, userID)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		var updateUser model.User
		if err := c.BindJSON(&updateUser); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Faild to bind user data",
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
		if updateUser.BirthDate != "" {
			user.BirthDate = updateUser.BirthDate
		}
		if updateUser.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 10)
			if err != nil {
				c.JSON(http.StatusBadRequest, model.ErrorResponse{
					Error: "Failed to hash password",
				})
				return
			}
			user.BirthDate = string(hash)
		}
		if err := initializers.DB.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to update user",
			})
			return
		}

		c.JSON(http.StatusOK, model.UpdateUserResponse{
			SuccessResponse: model.SuccessResponse{
				Message: "User updated successfuly",
			},
			User: user,
		})
	}
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve all users or filter by query parameters
// @Tags users
// @Produce json
// @Param name query string false "Name"
// @Param email query string false "Email"
// @Param age query int false "Age"
// @Param phone query string false "Phone"
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users [get]
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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

		var users []struct {
			model.User
			RoleID     string `gorm:"column:role_id"`
			Permission string `gorm:"column:permission"`
			IsActive   bool   `gorm:"column:is_active"`
			Department string `gorm:"column:department"`
		}

		var result *gorm.DB

		if queryExists || len(c.Params) == 0 {
			result = initializers.DB.Model(&model.User{}).
				Select("users.*, roles.role as role, roles.permission, roles.is_active, departments.name as department").
				Joins("left join roles on roles.id = users.role_id").
				Joins("left join departments on departments.id = roles.department_id").
				Where(&queryCondition).Scan(&users)
		} else {
			result = initializers.DB.Model(&model.User{}).
				Select("users.*, roles.role as role, roles.permission, roles.is_active, departments.name as department").
				Joins("left join roles on roles.id = users.role_id").
				Joins("left join departments on departments.id = roles.department_id").Scan(&users)
		}

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to retrieve users",
			})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "No users found matching the criteria",
			})
			return
		}

		var userResponses []model.UserResponse
		for _, user := range users {
			userResponses = append(userResponses, model.UserResponse{
				User:       user.User,
				RoleID:     user.RoleID,
				Permission: user.Permission,
				IsActive:   user.IsActive,
				Department: user.Department,
			})
		}

		c.JSON(http.StatusOK, model.UsersResponse{
			Message: "Users retrieved successfully",
			Users:   userResponses,
		})
	}
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/{id} [delete]
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

		result := initializers.DB.Delete(&model.User{}, userID)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to delete user",
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "No user found with the given ID",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "User deleted successfuly"})
	}
}

// RecoverUser godoc
// @Summary Recover a deleted user
// @Description Recover a soft-deleted user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/{id}/recover [post]
func RecoverUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

		result := initializers.DB.Model(&model.User{}).Unscoped().Where("id = ?", userID).Update("deleted_at", nil)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to recover user",
			})
			return
		}

		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "No deleted user found with the given ID",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "User recovered successfully"})
	}
}
