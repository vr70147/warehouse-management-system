package tests

import (
	"accounts-management/internal/api/handlers"
	"accounts-management/internal/model"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Print Routes for Debugging
func setupRoutes() (*gin.Engine, *gorm.DB) {
	r := gin.Default()
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.Account{}, &model.Plan{}, &model.User{})
	accounts := r.Group("/accounts")
	accounts.POST("/", handlers.CreateAccount(db))
	accounts.GET("/", handlers.GetAccounts(db))
	accounts.PUT("/:id", handlers.UpdateAccount(db))
	accounts.DELETE("/:id", handlers.SoftDeleteAccount(db))
	accounts.DELETE("/hard/:id", handlers.HardDeleteAccount(db))
	accounts.POST("/:id/recover", handlers.RecoverAccount(db))

	for _, route := range r.Routes() {
		fmt.Println(route.Method, route.Path)
	}

	return r, db
}

// Generate JWT token for testing
func generateToken(userID uint, accountID string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":        userID,
		"account_id": accountID,
	})
	tokenString, _ := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))
	return tokenString
}

// Test Create Account Route
func TestCreateAccountRoute(t *testing.T) {
	router, _ := setupRoutes()

	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		Password:    "password",
		PhoneNumber: "1234567890",
	}
	jsonValue, _ := json.Marshal(account)

	req, _ := http.NewRequest("POST", "/accounts/", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// Test Get Accounts Route
func TestGetAccountsRoute(t *testing.T) {
	router, db := setupRoutes()

	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		Password:    "password",
		PhoneNumber: "1234567890",
	}
	db.Create(&account)
	check, _ := http.NewRequest("GET", "/accounts/", nil)
	fmt.Println(check)
	token := generateToken(1, "account_1")
	req, _ := http.NewRequest("GET", "/accounts/", nil) // Added trailing slash
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test Update Account Route
func TestUpdateAccountRoute(t *testing.T) {
	router, db := setupRoutes()

	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		Password:    "password",
		PhoneNumber: "1234567890",
	}
	db.Create(&account)

	updatedAccount := model.Account{
		Email:       "updated@example.com",
		Name:        "Updated User",
		PhoneNumber: "0987654321",
	}
	jsonValue, _ := json.Marshal(updatedAccount)

	token := generateToken(1, "account_1")
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/accounts/%d", account.ID), bytes.NewBuffer(jsonValue))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test Soft Delete Account Route
func TestSoftDeleteAccountRoute(t *testing.T) {
	router, db := setupRoutes()

	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		Password:    "password",
		PhoneNumber: "1234567890",
	}
	db.Create(&account)

	token := generateToken(1, "account_1")
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/accounts/%d", account.ID), nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRecoverAccountRoute(t *testing.T) {
	router, db := setupRoutes()

	// Create an account
	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		Password:    "password",
		PhoneNumber: "1234567890",
	}
	db.Create(&account)
	fmt.Printf("Account created: %+v\n", account) // Debug output

	// Soft delete the account
	db.Delete(&account)
	fmt.Printf("Account soft-deleted: %+v\n", account) // Debug output

	// Verify account is soft-deleted
	var deletedAccount model.Account
	db.Unscoped().Where("id = ?", account.ID).First(&deletedAccount)
	fmt.Printf("Deleted account state: %+v\n", deletedAccount) // Debug output

	req, _ := http.NewRequest("POST", fmt.Sprintf("/accounts/%d/recover", account.ID), bytes.NewBuffer(nil))
	req.Header.Set("Content-Type", "application/json")

	// Print request details for debugging
	fmt.Printf("Request URL: %s\n", req.URL.String())
	fmt.Printf("Request Method: %s\n", req.Method)
	fmt.Printf("Request Headers: %v\n", req.Header)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Print response details for debugging
	fmt.Println("Response Code:", w.Code)
	fmt.Println("Response Body:", w.Body.String())

	assert.Equal(t, http.StatusOK, w.Code)
}

// Test Hard Delete Account Route
func TestHardDeleteAccountRoute(t *testing.T) {
	router, db := setupRoutes()

	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		Password:    "password",
		PhoneNumber: "1234567890",
	}
	db.Create(&account)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/accounts/hard/%d", account.ID), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
