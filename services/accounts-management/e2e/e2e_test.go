package handlers_test

import (
	"accounts-management/internal/api/routes"
	"accounts-management/internal/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	routes.Routers(r, db)
	return r
}

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	db.AutoMigrate(&model.Account{}, &model.Plan{})
	return db
}

func TestCreateAccount(t *testing.T) {
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
		Metadata:       "{}",
		Preferences:    "{}",
	}

	body, _ := json.Marshal(account)
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response model.SuccessResponse
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Account created successfully", response.Message)
}

func TestGetAccounts(t *testing.T) {
	db := setupTestDB()
	router := setupRouter(db)

	// Create an account to retrieve
	db.Create(&model.Account{
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
		Metadata:       "{}",
		Preferences:    "{}",
	})

	req, _ := http.NewRequest("GET", "/accounts", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response model.SuccessResponse
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Accounts retrieved successfully", response.Message)
	assert.NotEmpty(t, response.Data)
}

func TestUpdateAccount(t *testing.T) {
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
		Metadata:       "{}",
		Preferences:    "{}",
	}

	db.Create(&account)
	account.Name = "Updated User"
	body, _ := json.Marshal(account)
	req, _ := http.NewRequest("PUT", "/accounts/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response model.SuccessResponse
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Account updated successfully", response.Message)
	assert.Equal(t, "Updated User", response.Data.(map[string]interface{})["name"])
}

func TestSoftDeleteAccount(t *testing.T) {
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
		Metadata:       "{}",
		Preferences:    "{}",
	}

	db.Create(&account)
	req, _ := http.NewRequest("DELETE", "/accounts/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response model.SuccessResponse
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Account soft deleted successfully", response.Message)
}

func TestHardDeleteAccount(t *testing.T) {
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
		Metadata:       "{}",
		Preferences:    "{}",
	}

	db.Create(&account)
	req, _ := http.NewRequest("DELETE", "/accounts/hard/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response model.SuccessResponse
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Account hard deleted successfully", response.Message)
}

func TestRecoverAccount(t *testing.T) {
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
		Metadata:       "{}",
		Preferences:    "{}",
	}

	db.Create(&account)
	db.Delete(&account) // Soft delete the account

	req, _ := http.NewRequest("POST", "/accounts/1/recover", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	var response model.SuccessResponse
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Equal(t, "Account recovered successfully", response.Message)
}
