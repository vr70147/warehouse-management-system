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
