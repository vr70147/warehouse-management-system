package handlers_test

import (
	"accounts-management/internal/api/handlers"
	"accounts-management/internal/model"

	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	db.Migrator().DropTable(&model.Account{})
	db.AutoMigrate(&model.Account{})
	return db
}

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.POST("/accounts", handlers.CreateAccount(db))
	return r
}

func TestCreateAccountHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()
	router := setupRouter(db)

	account := model.Account{
		Email:          "test@example.com",
		Name:           "Test User",
		PhoneNumber:    "123456789",
		CompanyName:    "Test Company",
		Address:        "123 Test St",
		City:           "Test City",
		State:          "Test State",
		PostalCode:     "12345",
		Country:        "Test Country",
		BillingEmail:   "billing@example.com",
		BillingAddress: "123 Billing St",
		PlanID:         1,
		IsActive:       true,
		Metadata:       "Test Metadata",
		Preferences:    "Test Preferences",
	}

	jsonValue, _ := json.Marshal(account)
	req, _ := http.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err, "Error unmarshalling response")

	assert.NotNil(t, response["data"], "Response data should not be nil")

	data, ok := response["data"].(map[string]interface{})
	fmt.Println(data)
	assert.True(t, ok, "Response data should be a map")
	assert.NotNil(t, data["id"], "ID field should not be nil in response data")
}

func TestGetAccounts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.Default()
	router.GET("/accounts", handlers.GetAccounts(db))

	// Create a test account
	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		PhoneNumber: "123456789",
		CompanyName: "Test Company",
		Address:     "123 Test St",
		City:        "Test City",
		State:       "Test State",
		PostalCode:  "12345",
		Country:     "Test Country",
		IsActive:    true,
	}
	db.Create(&account)

	req, _ := http.NewRequest(http.MethodGet, "/accounts", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Accounts retrieved successfully", response["message"])
	assert.NotNil(t, response["data"])
}

func TestUpdateAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.Default()
	router.PUT("/accounts/:id", handlers.UpdateAccount(db))

	// Create a test account
	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		PhoneNumber: "123456789",
		CompanyName: "Test Company",
		Address:     "123 Test St",
		City:        "Test City",
		State:       "Test State",
		PostalCode:  "12345",
		Country:     "Test Country",
		IsActive:    true,
	}
	db.Create(&account)

	updatedData := model.Account{
		Email:       "updated@example.com",
		Name:        "Updated User",
		PhoneNumber: "987654321",
		CompanyName: "Updated Company",
		Address:     "456 Updated St",
		City:        "Updated City",
		State:       "Updated State",
		PostalCode:  "54321",
		Country:     "Updated Country",
		IsActive:    false,
	}
	jsonValue, _ := json.Marshal(updatedData)
	req, _ := http.NewRequest(http.MethodPut, "/accounts/"+strconv.Itoa(int(account.ID)), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Account updated successfully", response["message"])
	assert.NotNil(t, response["data"])
}

func TestSoftDeleteAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.Default()
	router.DELETE("/accounts/:id", handlers.SoftDeleteAccount(db))

	// Create a test account
	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		PhoneNumber: "123456789",
		CompanyName: "Test Company",
		Address:     "123 Test St",
		City:        "Test City",
		State:       "Test State",
		PostalCode:  "12345",
		Country:     "Test Country",
		IsActive:    true,
	}
	db.Create(&account)

	req, _ := http.NewRequest(http.MethodDelete, "/accounts/"+strconv.Itoa(int(account.ID)), nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Account soft deleted successfully", response["message"])
	assert.NotNil(t, response["data"])
}

func TestHardDeleteAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.Default()
	router.DELETE("/accounts/hard/:id", handlers.HardDeleteAccount(db))

	// Create a test account
	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		PhoneNumber: "123456789",
		CompanyName: "Test Company",
		Address:     "123 Test St",
		City:        "Test City",
		State:       "Test State",
		PostalCode:  "12345",
		Country:     "Test Country",
		IsActive:    true,
	}
	db.Create(&account)

	req, _ := http.NewRequest(http.MethodDelete, "/accounts/hard/"+strconv.Itoa(int(account.ID)), nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Account hard deleted successfully", response["message"])
	assert.NotNil(t, response["data"])
}

func TestRecoverAccount(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db := setupTestDB()

	router := gin.Default()
	router.POST("/accounts/:id/recover", handlers.RecoverAccount(db))

	// Create a test account
	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test User",
		PhoneNumber: "123456789",
		CompanyName: "Test Company",
		Address:     "123 Test St",
		City:        "Test City",
		State:       "Test State",
		PostalCode:  "12345",
		Country:     "Test Country",
		IsActive:    true,
	}
	db.Create(&account)

	// Soft delete the account
	db.Delete(&account)

	req, _ := http.NewRequest(http.MethodPost, "/accounts/"+strconv.Itoa(int(account.ID))+"/recover", nil)

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Account recovered successfully", response["message"])
	assert.NotNil(t, response["data"])
}
