package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"shipping-receiving/internal/api/routes"
	"shipping-receiving/internal/middleware"
	"shipping-receiving/internal/model"
	"shipping-receiving/internal/utils"
	"strconv"
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
	tokenString, _ := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return tokenString
}

func setupTestDB() (*gorm.DB, error) {
	os.Getenv("JWT_SECRET")
	os.Getenv("USER_SERVICE_URL")

	db, err := gorm.Open(sqlite.Open("test_shipping.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&model.Shipping{}, &model.User{}, &model.Role{}, &model.Account{}, &model.Department{})
	// Create a role and user for testing
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

	return db, nil
}

func TestCreateShipping(t *testing.T) {

	// Setup the test database
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_shipping.db")

	// Mock the order service URL
	os.Getenv("ORDER_SERVICE_URL")

	r := SetupRouter(db)

	t.Run("CreateShippingSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		shipping := model.Shipping{
			ReceiverID: 1,
			Status:     "pending",
		}
		jsonValue, _ := json.Marshal(shipping)
		req, _ := http.NewRequest("POST", "/shipping-receiving", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code) // Expecting 200 OK here
		var response model.SuccessResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Shipping created successfully", response.Message)
	})
}

func TestGetShippings(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_shipping.db")

	shipping := model.Shipping{
		ReceiverID: 1,
		Status:     "pending",
		AccountID:  1,
	}
	db.Create(&shipping)

	r := SetupRouter(db)

	t.Run("GetShippingsSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("GET", "/shipping-receiving", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponses
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Shippings retrieved successfully", response.Message)
		assert.Equal(t, 1, len(response.Data))
	})
	db.Exec("DELETE FROM shippings")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM departments")

}

func TestUpdateShipping(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_shipping.db")

	shipping := model.Shipping{
		ReceiverID: 1,
		Status:     "pending",
		AccountID:  1,
	}
	db.Create(&shipping)

	r := SetupRouter(db)

	t.Run("UpdateShippingSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		updatedShipping := model.Shipping{
			Status: "delivered",
		}
		jsonValue, _ := json.Marshal(updatedShipping)
		req, _ := http.NewRequest("PUT", "/shipping-receiving/"+strconv.Itoa(int(shipping.ID)), bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Shipping updated successfully", response.Message)
	})
	db.Exec("DELETE FROM shippings")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM departments")
}

func TestSoftDeleteShipping(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_shipping.db")

	shipping := model.Shipping{
		ReceiverID: 1,
		Status:     "pending",
		AccountID:  1,
	}
	db.Create(&shipping)

	r := SetupRouter(db)

	t.Run("SoftDeleteShippingSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("DELETE", "/shipping-receiving/"+strconv.Itoa(int(shipping.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Shipping soft deleted successfully", response.Message)
	})
	db.Exec("DELETE FROM shippings")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM departments")
}

func TestHardDeleteShipping(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_shipping.db")

	shipping := model.Shipping{
		ReceiverID: 1,
		Status:     "pending",
		AccountID:  1,
	}
	db.Create(&shipping)

	r := SetupRouter(db)

	t.Run("HardDeleteShippingSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("DELETE", "/shipping-receiving/hard/"+strconv.Itoa(int(shipping.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Shipping hard deleted successfully", response.Message)
	})
	db.Exec("DELETE FROM shippings")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM departments")
}

func TestRecoverShipping(t *testing.T) {
	os.Setenv("JWT_SECRET", "test_secret")

	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_shipping.db")

	shipping := model.Shipping{
		ReceiverID: 1,
		Status:     "pending",
		AccountID:  1,
	}
	db.Create(&shipping)

	// Soft delete the shipping
	db.Delete(&shipping)

	r := SetupRouter(db)

	t.Run("RecoverShippingSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("PATCH", "/shipping-receiving/"+strconv.Itoa(int(shipping.ID))+"/recover", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		fmt.Printf("w.Code: %d\n", w.Code)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		t.Log("Response Body:", w.Body.String()) // Add this line for debugging
		assert.Equal(t, "Shipping recovered successfully", response.Message)
	})
	db.Exec("DELETE FROM shippings")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM departments")
}

func TestDeliverShipping(t *testing.T) {
	// Set up the mock email sender
	mockEmailSender := new(MockEmailSender)
	mockEmailSender.On("SendEmail", "customer@example.com", "Your Order is Shipped", "Dear User, Your order has been shipped.").Return(nil)

	db, err := setupTestDB()
	assert.NoError(t, err)
	defer os.Remove("test_shipping.db")

	shipping := model.Shipping{
		ReceiverID: 1,
		Status:     "pending",
		AccountID:  1,
	}
	db.Create(&shipping)

	r := gin.Default()
	ns := utils.NewNotificationService(mockEmailSender)
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware(db))
	routes.Routers(r, db, ns)

	t.Run("DeliverShippingSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		req, _ := http.NewRequest("POST", "/shipping-receiving/"+strconv.Itoa(int(shipping.ID))+"/deliver", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Shipping delivered successfully", response.Message)

		// Verify that the SendEmail method was called
		mockEmailSender.AssertCalled(t, "SendEmail", "customer@example.com", "Your Order is Shipped", "Dear User, Your order has been shipped.")
	})
	db.Exec("DELETE FROM shippings")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM roles")
	db.Exec("DELETE FROM departments")
}
