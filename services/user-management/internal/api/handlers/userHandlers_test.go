package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"user-management/internal/api/handlers"
	"user-management/internal/model"
	"user-management/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	return r
}

type MockEmailSender struct {
	mock.Mock
}

func (m *MockEmailSender) SendEmail(to, subject, body string) error {
	args := m.Called(to, subject, body)
	return args.Error(0)
}

func TestSignup(t *testing.T) {
	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_user.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.User{}, &model.Role{})

	mockEmailSender := new(MockEmailSender)
	ns := utils.NewNotificationService(mockEmailSender)

	r := SetupRouter()
	r.POST("/signup", handlers.Signup(db, ns))

	t.Run("SignupSuccess", func(t *testing.T) {
		role := model.Role{
			ID: 1,
		}
		db.Create(&role)

		user := model.User{
			PersonalID: "12345",
			Name:       "Test User",
			Email:      "user@example.com",
			Age:        25,
			BirthDate:  "1999-01-01",
			RoleID:     role.ID,
			Phone:      "1234567890",
			Street:     "Test Street",
			City:       "Test City",
			Password:   "password123",
			IsAdmin:    false,
			AccountID:  1, // Ensure this is set
			Permission: model.PermissionWorker,
		}

		jsonValue, _ := json.Marshal(user)
		fmt.Println("Test Signup Request Body:", string(jsonValue)) // Log the outgoing request

		req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User registered successfully", response["message"])
	})

	// Clean up the database
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestLogin(t *testing.T) {
	// Set up the environment variable for JWT secret
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_user.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.User{})

	// Create a user for testing login
	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	role := model.Role{
		ID: 1,
		// Other fields if required
	}
	db.Create(&role)

	testUser := model.User{
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

	mockEmailSender := new(MockEmailSender)
	ns := utils.NewNotificationService(mockEmailSender)

	r := SetupRouter()
	r.POST("/login", handlers.Login(db, ns))

	// Define the test case
	t.Run("LoginSuccess", func(t *testing.T) {
		loginCredentials := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    "user@example.com",
			Password: "password123",
		}
		jsonValue, _ := json.Marshal(loginCredentials)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		t.Log("Response Body:", w.Body.String()) // Add this line for debugging

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User authenticated successfully", response["message"])
		assert.NotEmpty(t, response["data"]) // JWT token should be present in the response
	})

	// Define a test case for invalid credentials
	t.Run("LoginInvalidCredentials", func(t *testing.T) {
		loginCredentials := struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}{
			Email:    "user@example.com",
			Password: "wrongpassword",
		}
		jsonValue, _ := json.Marshal(loginCredentials)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		// Set up mock expectation
		mockEmailSender.On("SendEmail", "user@example.com", "Failed Login Attempt", "Dear User, There was a failed login attempt on your account.").Return(nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid email or password", response["error"])

		// Assert that the expected method was called
		// mockEmailSender.AssertExpectations(t)
	})

	// Clean up the database
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestUpdateUser(t *testing.T) {
	// Set up the environment variable for JWT secret
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_user.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.User{})

	// Create a user for testing login
	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	role := model.Role{
		ID: 1,
		// Other fields if required
	}
	db.Create(&role)

	testUser := model.User{
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

	r := SetupRouter()
	r.PUT("/users/:id", func(c *gin.Context) {
		c.Set("account_id", uint(1)) // Set the account ID in the context
		handlers.UpdateUser(db)(c)
	})

	// Define the test case for successful update
	t.Run("UpdateUserSuccess", func(t *testing.T) {
		updateData := struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Street   string `json:"street"`
			City     string `json:"city"`
			Age      int    `json:"age"`
			Password string `json:"password"`
		}{
			Name:     "Updated User",
			Email:    "updated@example.com",
			Phone:    "0987654321",
			Street:   "New Street",
			City:     "New City",
			Age:      30,
			Password: "newpassword123",
		}
		jsonValue, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", "/users/"+strconv.Itoa(int(testUser.ID)), bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User updated successfully", response["message"])
	})

	// Define the test case for invalid user ID
	t.Run("UpdateUserInvalidID", func(t *testing.T) {
		updateData := struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Street   string `json:"street"`
			City     string `json:"city"`
			Age      int    `json:"age"`
			Password string `json:"password"`
		}{
			Name:     "Updated User",
			Email:    "updated@example.com",
			Phone:    "0987654321",
			Street:   "New Street",
			City:     "New City",
			Age:      30,
			Password: "newpassword123",
		}
		jsonValue, _ := json.Marshal(updateData)
		req, _ := http.NewRequest("PUT", "/users/9999", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User not found", response["error"])
	})

	// Clean up the database
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestSoftDeleteUser(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_user.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(model.User{})

	password, _ := bcrypt.GenerateFromPassword(([]byte("password123")), bcrypt.DefaultCost)

	role := model.Role{
		ID: 1,
	}

	db.Create(&role)

	testUser := model.User{
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

	r := SetupRouter()
	r.DELETE("/users/:id", func(c *gin.Context) {
		c.Set("account_id", uint(1)) // Set the account ID in the context
		handlers.SoftDeleteUser(db)(c)
	})

	// Define the test case
	t.Run("SoftDeleteUserSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/"+strconv.Itoa(int(testUser.ID)), nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User soft deleted successfully", response["message"])

		// Check if the user is deleted from the database
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM roles")
	})

	// Define the test case for invalid user ID
	t.Run("SoftDeleteUserInvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/9999", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User not found", response["error"])

		// Check if the user is deleted from the database
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM roles")
	})
}

func TestHardDeleteUser(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_user.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(model.User{})

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	role := model.Role{
		ID: 1,
	}

	db.Create(&role)

	testUser := model.User{
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

	r := SetupRouter()
	r.DELETE("/users/hard/:id", func(c *gin.Context) {
		c.Set("account_id", uint(1)) // Set the account ID in the context
		handlers.HardDeleteUser(db)(c)
	})

	// Define the test case
	t.Run("HardDeleteUserSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/hard/"+strconv.Itoa(int(testUser.ID)), nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User hard deleted successfully", response["message"])

		// Check if the user is deleted from the database
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM roles")
	})

	// Define the test case for invalid user ID
	t.Run("HardDeleteUserInvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/users/hard/9999", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		t.Log("Response Body:", w.Body.String())

		assert.Equal(t, "User not found", response["error"])

		// Check if the user is deleted from the database
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM roles")
	})
}

func TestRecoverUser(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")
	// Setup the database
	db, err := gorm.Open(sqlite.Open("test_user.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(model.User{})

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	role := model.Role{
		ID: 1,
	}

	db.Create(&role)

	testUser := model.User{
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

	r := SetupRouter()
	r.PATCH("/users/recover/:id", func(c *gin.Context) {
		c.Set("account_id", uint(1)) // Set the account ID in the context
		handlers.RecoverUser(db)(c)
	})

	// Define the test case
	t.Run("RecoverUserSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/users/recover/"+strconv.Itoa(int(testUser.ID)), nil)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User recovered successfully", response["message"])

		// Check if the user is deleted from the database
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM roles")
	})

	// Define the test case for invalid user ID
	t.Run("RecoverUserInvalidID", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/users/recover/9999", nil)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		t.Log("Response Body:", w.Body.String())

		assert.Equal(t, "User not found", response["error"])

		// Check if the user is deleted from the database
		db.Exec("DELETE FROM users")
		db.Exec("DELETE FROM roles")
	})
}
