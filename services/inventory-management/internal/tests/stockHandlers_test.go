package tests_test

import (
	"bytes"
	"encoding/json"
	"inventory-management/internal/model"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestEnvironment() (*gorm.DB, string, *model.User) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8080")
	os.Setenv("TOKEN_SECRET", "test_secret")

	db, err := gorm.Open(sqlite.Open("test_inventory.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.Stock{}, &model.Product{}, &model.User{}, &model.Role{})

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

	token := createTestToken(testUser.ID, testUser.AccountID)

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

	return db, token, &testUser
}

func TestCreateStock(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	t.Run("CreateStockSuccess", func(t *testing.T) {
		product := model.Product{
			ID:        1,
			Name:      "Test Product",
			AccountID: testUser.AccountID,
		}
		db.Create(&product)

		stock := model.Stock{
			ProductID: product.ID,
			Quantity:  100,
			Location:  "Warehouse 1",
			AccountID: testUser.AccountID,
		}
		jsonValue, _ := json.Marshal(stock)
		req, _ := http.NewRequest("POST", "/stocks", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		log.Printf("Request Payload: %s", jsonValue)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		log.Printf("Response Code: %d", w.Code)
		log.Printf("Response Body: %s", w.Body.String())

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.Stock
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, stock.ProductID, response.ProductID)
		assert.Equal(t, stock.Quantity, response.Quantity)
		assert.Equal(t, stock.Location, response.Location)
	})

	// Clean up the database
	db.Exec("DELETE FROM stocks")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestGetStocks(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	// Create a product for testing
	product := model.Product{
		ID:        1,
		Name:      "Test Product",
		AccountID: testUser.AccountID,
	}
	db.Create(&product)

	// Create stock items for testing
	stock1 := model.Stock{
		ID:        1,
		ProductID: product.ID,
		Quantity:  100,
		Location:  "Warehouse 1",
		AccountID: testUser.AccountID,
	}
	stock2 := model.Stock{
		ID:        2,
		ProductID: product.ID,
		Quantity:  200,
		Location:  "Warehouse 2",
		AccountID: testUser.AccountID,
	}
	db.Create(&stock1)
	db.Create(&stock2)

	t.Run("GetStocksSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/stocks", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.StocksResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Stocks retrieved successfully", response.Message)
		assert.Equal(t, 2, len(response.Stocks))
	})

	// Clean up the database
	db.Exec("DELETE FROM stocks")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestSoftDeleteStock(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	// Create a product for testing
	product := model.Product{
		ID:        1,
		Name:      "Test Product",
		AccountID: testUser.AccountID,
	}
	db.Create(&product)

	// Create a stock item for testing
	stock := model.Stock{
		ID:        1,
		ProductID: product.ID,
		Quantity:  100,
		Location:  "Warehouse 1",
		AccountID: testUser.AccountID,
	}
	db.Create(&stock)

	t.Run("SoftDeleteStockSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/stocks/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Stock deleted successfully", response["message"])
	})

	// Clean up the database
	db.Exec("DELETE FROM stocks")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestHardDeleteStock(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	// Create a product for testing
	product := model.Product{
		ID:        1,
		Name:      "Test Product",
		AccountID: testUser.AccountID,
	}
	db.Create(&product)

	// Create a stock item for testing
	stock := model.Stock{
		ID:        1,
		ProductID: product.ID,
		Quantity:  100,
		Location:  "Warehouse 1",
		AccountID: testUser.AccountID,
	}
	db.Create(&stock)

	t.Run("HardDeleteStockSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/stocks/hard/1", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Stock deleted permanently", response["message"])
	})

	// Clean up the database
	db.Exec("DELETE FROM stocks")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestRecoverStock(t *testing.T) {
	db, token, testUser := setupTestEnvironment()
	r := SetupRouter(db)

	// Create a product for testing
	product := model.Product{
		ID:        1,
		Name:      "Test Product",
		AccountID: testUser.AccountID,
	}
	db.Create(&product)

	// Create a stock item for testing
	stock := model.Stock{
		ID:        1,
		ProductID: product.ID,
		Quantity:  100,
		Location:  "Warehouse 1",
		AccountID: testUser.AccountID,
	}
	db.Create(&stock)

	// Soft delete the stock item
	db.Delete(&stock)

	t.Run("RecoverStockSuccess", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/stocks/1/recover", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Stock recovered successfully", response["message"])
	})

	// Clean up the database
	db.Exec("DELETE FROM stocks")
	db.Exec("DELETE FROM products")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}
