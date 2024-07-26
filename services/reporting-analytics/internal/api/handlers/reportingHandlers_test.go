package handlers_test

import (
	"os"
	"reporting-analytics/internal/api/routes"
	"reporting-analytics/internal/middleware"
	"reporting-analytics/internal/model"
	"testing"
	"time"

	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() (*gorm.DB, error) {
	os.Setenv("TOKEN_SECRET", "test_secret")

	db, err := gorm.Open(sqlite.Open("test_reporting.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.SalesReport{}, &model.InventoryLevel{}, &model.ShippingStatus{}, &model.UserActivity{}, &model.User{}, &model.Role{}, &model.Account{})

	// Create test data
	role := model.Role{
		ID: 1,
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

	account := model.Account{
		ID:   1,
		Name: "Test Account",
	}
	db.Create(&account)

	go func() {
		r := gin.Default()
		r.GET("/users/:id", func(c *gin.Context) {
			userID := c.Param("id")
			if userID == "1" {
				c.JSON(200, testUser)
			} else {
				c.JSON(404, gin.H{"error": "User not found"})
			}
		})
		r.Run(":8080")
	}()

	return db, nil
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

func setupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware(db))
	routes.Routers(r, db)
	return r
}

func TestGetSalesReports(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_reporting.db")

	// Create test data
	salesReport := model.SalesReport{
		AccountID:  1,
		TotalSales: 100.0,
	}
	db.Create(&salesReport)

	r := setupRouter(db)

	t.Run("GetSalesReportsSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("GET", "/reports/sales", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.SalesReport
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, 100.0, response[0].TotalSales)
	})
	db.Exec("DELETE FROM sales_reports")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM inventory_levels")
	db.Exec("DELETE FROM shipping_statuses")
	db.Exec("DELETE FROM user_activities")
}

func TestGetInventoryLevels(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_reporting.db")

	// Create test data
	inventoryLevel := model.InventoryLevel{
		AccountID: 1,
		ProductID: 1,
		Quantity:  50,
	}
	db.Create(&inventoryLevel)

	r := setupRouter(db)

	t.Run("GetInventoryLevelsSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("GET", "/reports/inventory", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.InventoryLevel
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, 50, response[0].Quantity)
	})
	db.Exec("DELETE FROM sales_reports")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM inventory_levels")
	db.Exec("DELETE FROM shipping_statuses")
	db.Exec("DELETE FROM user_activities")
}

func TestGetShippingStatuses(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_reporting.db")

	// Create test data
	shippingStatus := model.ShippingStatus{
		AccountID: 1,
		OrderID:   1,
		Status:    "shipped",
	}
	db.Create(&shippingStatus)

	r := setupRouter(db)

	t.Run("GetShippingStatusesSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("GET", "/reports/shipping", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.ShippingStatus
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, "shipped", response[0].Status)
	})
	db.Exec("DELETE FROM sales_reports")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM inventory_levels")
	db.Exec("DELETE FROM shipping_statuses")
	db.Exec("DELETE FROM user_activities")
}

func TestGetUserActivities(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_reporting.db")

	// Create test data
	userActivity := model.UserActivity{
		AccountID: 1,
		UserID:    1,
		Action:    "login",
	}
	db.Create(&userActivity)

	r := setupRouter(db)

	t.Run("GetUserActivitiesSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("GET", "/reports/user-activity", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response []model.UserActivity
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, 1, len(response))
		assert.Equal(t, "login", response[0].Action)
	})
	db.Exec("DELETE FROM sales_reports")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM inventory_levels")
	db.Exec("DELETE FROM shipping_statuses")
	db.Exec("DELETE FROM user_activities")
}
