package handlers

import (
	"net/http"
	"os"
	"strconv"
	"time"
	"user-management/internal/model"
	"user-management/internal/utils"

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
func Signup(db *gorm.DB, ns *utils.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			PersonalID string           `json:"personal_id" gorm:"unique;not null"`
			Name       string           `json:"name" gorm:"unique;not null"`
			Email      string           `json:"email" gorm:"unique;not null"`
			Age        int              `json:"age" gorm:"not null"`
			BirthDate  string           `json:"birthDate" gorm:"not null"`
			RoleID     uint             `json:"role_id" gorm:"not null"`
			Permission model.Permission `json:"permission" gorm:"not null"`
			Phone      string           `json:"phone" gorm:"unique; not null"`
			Street     string           `json:"street"`
			City       string           `json:"city"`
			Password   string           `json:"password" gorm:"not null"`
			IsAdmin    bool             `json:"is_admin" gorm:"default: false"`
			AccountID  *uint            `json:"account_id,omitempty"` // Optional Account ID
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
		if result := db.Where("id = ?", body.RoleID).First(&role); result.Error != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Role not found",
			})
			return
		}

		accountID := body.AccountID
		if accountID == nil {
			// Try to get account ID from the context if it exists (in case of authenticated admin creating a user)
			if accID, exists := c.Get("account_id"); exists {
				id := accID.(uint)
				accountID = &id
			} else {
				// If no account ID is provided and not found in the context, return an error
				c.JSON(http.StatusBadRequest, model.ErrorResponse{
					Error: "Account ID is required",
				})
				return
			}
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
			AccountID:  *accountID,
			Permission: body.Permission,
		}
		if result := db.Create(&user); result.Error != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Failed to create user: " + result.Error.Error()})
			return
		}

		// Send email to user
		if err := ns.SendUserRegistrationNotification(user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to send notification email"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "User registered successfully", Data: user})
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
func Login(db *gorm.DB, ns *utils.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body struct {
			Email    string `json:"email" gorm:"unique"`
			Password string `json:"password"`
		}

		if c.Bind(&body) != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to read body",
			})
			return
		}

		var user model.User
		if err := db.First(&user, "email = ?", body.Email).Error; err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid email or password",
			})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if err != nil {
			// Notify user of failed login attempt
			if err := ns.SendFailedLoginAttemptNotification(user.Email); err != nil {
				c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to send notification email"})
				return
			}

			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Invalid email or password",
			})
			return
		}

		// Include accountID in the JWT token claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":        user.ID,
			"account_id": user.AccountID,
			"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(),
		})

		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to create token",
			})
			return
		}

		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("Authorization", tokenString, 3600*24*30, "", "", false, true)
		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "User authenticated successfully",
			Data:    tokenString,
		})
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
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		userID := c.Param("id")
		var user model.User

		result := db.Where("id = ? AND account_id = ?", userID, accountID).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "User not found",
			})
			return
		}

		var updateUser model.User
		if err := c.BindJSON(&updateUser); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Error: "Failed to bind user data",
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
			// updating birth date as well
			birthYear := time.Now().Year() - user.Age
			user.BirthDate = strconv.Itoa(birthYear) + "-01-01"
		}
		if updateUser.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(updateUser.Password), 10)
			if err != nil {
				c.JSON(http.StatusBadRequest, model.ErrorResponse{
					Error: "Failed to hash password",
				})
				return
			}
			user.Password = string(hash)
		}
		if err := db.Save(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to update user",
			})
			return
		}

		c.JSON(http.StatusOK, model.UpdateUserResponse{
			SuccessResponse: model.SuccessResponse{
				Message: "User updated successfully",
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
// @Success 200 {object} model.UsersResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users [get]
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		queryCondition := model.User{AccountID: accountID.(uint)}

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
			result = db.Model(&model.User{}).
				Select("users.*, roles.role as role, roles.permission, roles.is_active, departments.name as department").
				Joins("left join roles on roles.id = users.role_id").
				Joins("left join departments on departments.id = roles.department_id").
				Where(&queryCondition).Scan(&users)
		} else {
			result = db.Model(&model.User{}).
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
				Role:       user.RoleID,
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

// SoftDeleteUser godoc
// @Summary Soft delete a user
// @Description Soft delete a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/{id} [delete]
func SoftDeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		userID := c.Param("id")

		var user model.User
		if result := db.Where("id = ? AND account_id = ?", userID, accountID).First(&user); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "User not found",
			})
			return
		}

		if result := db.Delete(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to soft delete user",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "User soft deleted successfully",
		})
	}
}

// HardDeleteUser godoc
// @Summary Hard delete a user
// @Description Hard delete a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /users/hard/{id} [delete]
func HardDeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		userID := c.Param("id")

		var user model.User

		result := db.Where("id = ? AND account_id = ?", userID, accountID).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "User not found",
			})
			return
		}

		if result := db.Unscoped().Where("id = ? AND account_id = ?", userID, accountID).Delete(&model.User{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to hard delete user",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "User hard deleted successfully",
		})
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
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		userID := c.Param("id")

		var user model.User
		if result := db.Unscoped().Where("id = ? AND account_id = ?", userID, accountID).First(&user); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "User not found",
			})
			return
		}

		if result := db.Model(&user).Update("deleted_at", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to recover user",
			})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "User recovered successfully",
		})
	}
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change the password for a user and send a notification email
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body model.ChangePasswordRequest true "Password data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /users/change-password/{id} [post]
func ChangePassword(db *gorm.DB, ns *utils.NotificationService) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		userID := c.Param("id")
		var request model.ChangePasswordRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		var user model.User
		if result := db.Where("id = ? AND account_id = ?", userID, accountID).First(&user); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{
				Error: "User not found",
			})
			return
		}
		// Verify the current password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.CurrentPassword)); err != nil {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Current password is incorrect"})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to hash password"})
			return
		}

		user.Password = string(hashedPassword)
		if result := db.Save(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Error: "Failed to change password",
			})
			return
		}

		// Send notification email
		if err := ns.SendPasswordChangeNotification(user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to send notification email"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{
			Message: "Password changed successfully",
		})
	}
}
