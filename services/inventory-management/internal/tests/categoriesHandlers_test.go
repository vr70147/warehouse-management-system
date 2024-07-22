package tests_test

import (
	"bytes"
	"encoding/json"
	"inventory-management/internal/api/routes"
	"inventory-management/internal/middleware"
	"inventory-management/internal/model"
	"inventory-management/internal/utils"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware(db))
	ns := utils.NewNotificationService(&MockEmailSender{})
	routes.Routers(r, db, ns)
	return r
}

func createTestToken(userID uint, accountID uint) string {
	claims := jwt.MapClaims{
		"sub":        userID,
		"account_id": accountID,
		"exp":        time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte("test_secret"))
	return tokenString
}

func TestCreateCategory(t *testing.T) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8080")
	os.Setenv("TOKEN_SECRET", "test_secret")

	db, err := gorm.Open(sqlite.Open("test_inventory.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Category{}, &model.User{}, &model.Role{})

	// Create a role and user for testing
	role := model.Role{
		ID: 1,
		// Other fields if required
	}
	db.Create(&role)

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := model.User{
		ID:         1,
		PersonalID: "12345",
		Name:       "Test User",
		Email:      "user@example.com",
		Age:        25,
		BirthDate:  "1999-01-01",
		RoleID:     role.ID,
		Phone:      "1234567890",
		Street:     "Test Street",
		City:       "Test City",
		Password:   string(password),
		IsAdmin:    false,
		AccountID:  1,
		Permission: model.PermissionWorker,
	}
	db.Create(&testUser)

	// Mock user service
	go func() {
		r := gin.Default()
		r.GET("/users/:id", func(c *gin.Context) {
			userID := c.Param("id")
			if userID == "1" {
				c.JSON(http.StatusOK, testUser)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			}
		})
		r.Run(":8080")
	}()

	r := SetupRouter(db)

	t.Run("CreateCategorySuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		category := model.Category{
			Name: "Test Category",
		}
		jsonValue, _ := json.Marshal(category)
		req, _ := http.NewRequest("POST", "/categories", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Category created successfully", response["message"])
	})

	// Clean up the database
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestUpdateCategory(t *testing.T) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8080")
	os.Setenv("TOKEN_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_inventory.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Category{}, &model.User{}, &model.Role{})

	// Create a role and user for testing
	role := model.Role{
		ID: 1,
		// Other fields if required
	}
	db.Create(&role)

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := model.User{
		ID:         1,
		PersonalID: "12345",
		Name:       "Test User",
		Email:      "user@example.com",
		Age:        25,
		BirthDate:  "1999-01-01",
		RoleID:     role.ID,
		Phone:      "1234567890",
		Street:     "Test Street",
		City:       "Test City",
		Password:   string(password),
		IsAdmin:    false,
		AccountID:  1,
		Permission: model.PermissionWorker,
	}
	db.Create(&testUser)

	// Create a category for testing
	category := model.Category{
		ID:        1,
		Name:      "Old Category",
		AccountID: 1,
	}
	db.Create(&category)

	// Mock user service
	go func() {
		r := gin.Default()
		r.GET("/users/:id", func(c *gin.Context) {
			userID := c.Param("id")
			if userID == "1" {
				c.JSON(http.StatusOK, testUser)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			}
		})
		r.Run(":8080")
	}()

	r := SetupRouter(db)

	t.Run("UpdateCategorySuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		updatedCategory := model.Category{
			Name: "Updated Category",
		}
		jsonValue, _ := json.Marshal(updatedCategory)
		req, _ := http.NewRequest("PUT", "/categories/1", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Category updated successfully", response["message"])
	})

	// Clean up the database
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestGetCategories(t *testing.T) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8080")
	os.Setenv("TOKEN_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_inventory.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Category{}, &model.User{}, &model.Role{})

	// Create a role and user for testing
	role := model.Role{
		ID: 1,
		// Other fields if required
	}
	db.Create(&role)

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := model.User{
		ID:         1,
		PersonalID: "12345",
		Name:       "Test User",
		Email:      "user@example.com",
		Age:        25,
		BirthDate:  "1999-01-01",
		RoleID:     role.ID,
		Phone:      "1234567890",
		Street:     "Test Street",
		City:       "Test City",
		Password:   string(password),
		IsAdmin:    false,
		AccountID:  1,
		Permission: model.PermissionWorker,
	}
	db.Create(&testUser)

	// Create categories for testing
	category1 := model.Category{
		ID:        1,
		Name:      "Category 1",
		AccountID: 1,
	}
	category2 := model.Category{
		ID:        2,
		Name:      "Category 2",
		AccountID: 1,
	}
	db.Create(&category1)
	db.Create(&category2)

	// Mock user service
	go func() {
		r := gin.Default()
		r.GET("/users/:id", func(c *gin.Context) {
			userID := c.Param("id")
			if userID == "1" {
				c.JSON(http.StatusOK, testUser)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			}
		})
		r.Run(":8080")
	}()

	r := SetupRouter(db)

	t.Run("GetCategoriesSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("GET", "/categories", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.CategoriesResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Categories retrieved successfully", response.Message)
		assert.Equal(t, 2, len(response.Categories))
	})

	// Clean up the database
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestSoftDeleteCategory(t *testing.T) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8080")
	os.Setenv("TOKEN_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_inventory.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Category{}, &model.Product{}, &model.User{}, &model.Role{})

	// Create a role and user for testing
	role := model.Role{
		ID: 1,
		// Other fields if required
	}
	db.Create(&role)

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := model.User{
		ID:         1,
		PersonalID: "12345",
		Name:       "Test User",
		Email:      "user@example.com",
		Age:        25,
		BirthDate:  "1999-01-01",
		RoleID:     role.ID,
		Phone:      "1234567890",
		Street:     "Test Street",
		City:       "Test City",
		Password:   string(password),
		IsAdmin:    false,
		AccountID:  1,
		Permission: model.PermissionWorker,
	}
	db.Create(&testUser)

	// Create a category for testing
	category := model.Category{
		ID:        1,
		Name:      "Test Category",
		AccountID: 1,
	}
	db.Create(&category)

	// Create a default category for reassignment
	defaultCategory := model.Category{
		ID:        2,
		Name:      "Uncategorized",
		AccountID: 1,
	}
	db.Create(&defaultCategory)

	// Mock user service
	go func() {
		r := gin.Default()
		r.GET("/users/:id", func(c *gin.Context) {
			userID := c.Param("id")
			if userID == "1" {
				c.JSON(http.StatusOK, testUser)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			}
		})
		r.Run(":8080")
	}()

	r := SetupRouter(db)

	t.Run("SoftDeleteCategorySuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("DELETE", "/categories/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Category deleted successfully", response["message"])
	})

	// Clean up the database
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM products")
}

func TestHardDeleteCategory(t *testing.T) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8080")
	os.Setenv("TOKEN_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_inventory.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Category{}, &model.Product{}, &model.User{}, &model.Role{})

	// Create a role and user for testing
	role := model.Role{
		ID: 1,
		// Other fields if required
	}
	db.Create(&role)

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := model.User{
		ID:         1,
		PersonalID: "12345",
		Name:       "Test User",
		Email:      "user@example.com",
		Age:        25,
		BirthDate:  "1999-01-01",
		RoleID:     role.ID,
		Phone:      "1234567890",
		Street:     "Test Street",
		City:       "Test City",
		Password:   string(password),
		IsAdmin:    false,
		AccountID:  1,
		Permission: model.PermissionWorker,
	}
	db.Create(&testUser)

	// Create a category for testing
	category := model.Category{
		ID:        1,
		Name:      "Test Category",
		AccountID: 1,
	}
	db.Create(&category)

	// Create a default category for reassignment
	defaultCategory := model.Category{
		ID:        2,
		Name:      "Uncategorized",
		AccountID: 1,
	}
	db.Create(&defaultCategory)

	// Mock user service
	go func() {
		r := gin.Default()
		r.GET("/users/:id", func(c *gin.Context) {
			userID := c.Param("id")
			if userID == "1" {
				c.JSON(http.StatusOK, testUser)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			}
		})
		r.Run(":8080")
	}()

	r := SetupRouter(db)

	t.Run("HardDeleteCategorySuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("DELETE", "/categories/hard/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		t.Log("Response Body:", w.Body.String())

		assert.Equal(t, "Category deleted successfully", response["message"])
	})

	// Clean up the database
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestRecoverCategory(t *testing.T) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8080")
	os.Setenv("TOKEN_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_inventory.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Category{}, &model.Product{}, &model.User{}, &model.Role{})

	// Create a role and user for testing
	role := model.Role{
		ID: 1,
		// Other fields if required
	}
	db.Create(&role)

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	testUser := model.User{
		ID:         1,
		PersonalID: "12345",
		Name:       "Test User",
		Email:      "user@example.com",
		Age:        25,
		BirthDate:  "1999-01-01",
		RoleID:     role.ID,
		Phone:      "1234567890",
		Street:     "Test Street",
		City:       "Test City",
		Password:   string(password),
		IsAdmin:    false,
		AccountID:  1,
		Permission: model.PermissionWorker,
	}
	db.Create(&testUser)

	// Create a category for testing
	category := model.Category{
		ID:        1,
		Name:      "Test Category",
		AccountID: 1,
	}
	db.Create(&category)

	// Soft delete the category
	db.Delete(&category)

	go func() {
		r := gin.Default()
		r.GET("/users/:id", func(c *gin.Context) {
			userID := c.Param("id")
			if userID == "1" {
				c.JSON(http.StatusOK, testUser)
			} else {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			}
		})
		r.Run(":8080")
	}()

	r := SetupRouter(db)

	t.Run("RecoverCategorySuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("PATCH", "/categories/recover/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		t.Log("Response Body:", w.Body.String())

		assert.Equal(t, "Category recovered successfully", response["message"])
	})

	// Clean up the database
	db.Exec("DELETE FROM categories")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM products")
}
