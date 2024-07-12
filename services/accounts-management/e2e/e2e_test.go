package e2e

import (
	"accounts-management/internal/api/handlers"
	"accounts-management/internal/model"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupRouter() *gin.Engine {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&model.Account{}, &model.Plan{}) // Include Plan in the migration

	router := gin.Default()
	router.POST("/accounts", handlers.CreateAccount(db))
	router.GET("/accounts", handlers.GetAccounts(db))
	return router
}

func TestE2E_CreateAndGetAccounts(t *testing.T) {
	router := setupRouter()

	// Create Account
	account := model.Account{
		Email:          "test@example.com",
		Name:           "Test Account",
		PhoneNumber:    "123456789",
		CompanyName:    "test",
		BillingEmail:   "billing@example.com",
		BillingAddress: "test",
		PlanID:         1,
		IsActive:       true,
		Metadata:       "{}",
		Preferences:    "{}",
	}
	body, _ := json.Marshal(account)
	req, _ := http.NewRequest("POST", "/accounts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK; got %v", w.Code)
	}

	// Get Accounts
	req, _ = http.NewRequest("GET", "/accounts", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK; got %v", w.Code)
	}

	var accounts []model.Account
	err := json.Unmarshal(w.Body.Bytes(), &accounts)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	expected := []model.Account{
		{
			Email:          "test@example.com",
			Name:           "Test Account",
			PhoneNumber:    "123456789",
			CompanyName:    "test",
			BillingEmail:   "billing@example.com",
			BillingAddress: "test",
			PlanID:         1,
			IsActive:       true,
			Metadata:       "{}",
			Preferences:    "{}",
		},
	}

	if diff := cmp.Diff(expected, accounts, cmpopts.IgnoreFields(model.Account{}, "ID", "CreatedAt", "UpdatedAt", "DeletedAt", "Password", "Plan")); diff != "" {
		t.Fatalf("unexpected accounts (-want +got):\n%s", diff)
	}
}
