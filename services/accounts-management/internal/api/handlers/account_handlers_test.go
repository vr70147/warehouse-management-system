package handlers

import (
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

func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.Account{})
	return db
}

func TestCreateAccount(t *testing.T) {
	db := setupTestDB()
	router := gin.Default()
	router.POST("/accounts", CreateAccount(db))

	account := model.Account{
		Email:       "test@example.com",
		Name:        "Test Account",
		PhoneNumber: "123456789",
	}

	body, _ := json.Marshal(account)
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var responseAccount model.Account
	json.Unmarshal(w.Body.Bytes(), &responseAccount)
	assert.Equal(t, account.Email, responseAccount.Email)
	assert.Equal(t, account.Name, responseAccount.Name)
}

func TestGetAccounts(t *testing.T) {
	db := setupTestDB()
	db.Create(&model.Account{
		Email:          "test@example.com",
		Name:           "Test Account",
		PhoneNumber:    "123456789",
		Password:       "password",
		PlanID:         1,
		IsActive:       true,
		Metadata:       "{}",
		Preferences:    "{}",
		Plan:           model.Plan{ID: 1, Name: "test", Price: 0.0, Description: "something"},
		CompanyName:    "test",
		Address:        "test",
		City:           "test",
		State:          "test",
		PostalCode:     "test",
		Country:        "test",
		BillingEmail:   "test",
		BillingAddress: "test",
	})

	router := gin.Default()
	router.GET("/accounts", GetAccounts(db))

	req, _ := http.NewRequest("GET", "/accounts", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var accounts []model.Account
	json.Unmarshal(w.Body.Bytes(), &accounts)
	assert.Equal(t, 1, len(accounts))
	assert.Equal(t, "test@example.com", accounts[0].Email)
}
