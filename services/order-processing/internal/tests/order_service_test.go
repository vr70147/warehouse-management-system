package tests_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"order-processing/internal/api/routes"
	"order-processing/internal/middleware"
	"order-processing/internal/model"
	"order-processing/internal/utils"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	kafka_go "github.com/segmentio/kafka-go"
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

type MockKafkaWriter struct {
	mock.Mock
}

func (m *MockKafkaWriter) WriteMessages(ctx context.Context, msgs ...kafka_go.Message) error {
	args := m.Called(ctx, msgs)
	return args.Error(0)
}

func (m *MockKafkaWriter) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestGetOrders(t *testing.T) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8081") // Set the USER_SERVICE_URL environment variable
	os.Setenv("TOKEN_SECRET", "test_secret")               // Set the TOKEN_SECRET environment variable

	db, err := gorm.Open(sqlite.Open("test_order.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Order{}, &model.User{}, &model.Role{})

	// Create a role and user for testing login
	role := model.Role{
		ID: 1,
		// Other fields if required
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

	// Mock user service
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
		r.Run(":8081")
	}()

	r := SetupRouter()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware(db))

	// Mock email service
	ns := &utils.NotificationService{}

	routes.Routers(r, db, ns)

	// Seed some orders for testing
	db.Create(&model.Order{AccountID: 1, CustomerID: 1, Quantity: 1, ProductID: 1, Status: "Pending"})
	db.Create(&model.Order{AccountID: 1, CustomerID: 2, Quantity: 2, ProductID: 2, Status: "Shipped"})

	t.Run("GetOrdersSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)
		req, _ := http.NewRequest("GET", "/orders", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponses
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Orders found", response.Message)
		assert.Greater(t, len(response.Orders), 0)
	})

	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestCreateOrder(t *testing.T) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8080") // Set the USER_SERVICE_URL environment variable
	os.Setenv("TOKEN_SECRET", "test_secret")               // Set the TOKEN_SECRET environment variable
	os.Setenv("KAFKA_BROKERS", "localhost:9092")           // Set Kafka brokers
	os.Setenv("ORDER_EVENTS_TOPIC", "order-events")        // Set Kafka topic

	db, err := gorm.Open(sqlite.Open("test_order.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Order{}, &model.User{}, &model.Role{})

	// Clean up the database before and after the test
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")

	// Create a role and user for testing
	role := model.Role{
		// Other fields if required
	}
	db.Create(&role)

	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
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

	// Mock user service
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

	r := SetupRouter()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.AuthMiddleware(db))

	mockEmailSender := new(MockEmailSender)
	mockEmailSender.On("SendEmail", "customer@example.com", "Your Order is Complete", "Dear User, Your order has been completed successfully.").Return(nil)
	ns := utils.NewNotificationService(mockEmailSender)

	// Mock Kafka writer
	mockKafkaWriter := new(MockKafkaWriter)
	mockKafkaWriter.On("WriteMessages", mock.Anything, mock.AnythingOfType("kafka.Message")).Return(nil)
	mockKafkaWriter.On("Close").Return(nil)
	// kafka.KafkaWriterInstance = mockKafkaWriter

	routes.Routers(r, db, ns)

	t.Run("CreateOrderSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)

		order := model.Order{
			CustomerID: 1,
			Quantity:   1,
			ProductID:  1,
		}
		jsonValue, _ := json.Marshal(order)
		req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Order created successfully", response["message"])
	})

	// Clean up the database
	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}

func TestUpdateOrder(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test_order.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Order{})

	r := SetupRouter()
	r.Use(middleware.CORSMiddleware())

	ns := &utils.NotificationService{}

	routes.Routers(r, db, ns)

	// Seed an order for testing
	order := model.Order{AccountID: 1, CustomerID: 1, Quantity: 1, ProductID: 1, Status: "Pending", Version: 1}
	db.Create(&order)

	t.Run("UpdateOrderSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)
		orderUpdate := model.Order{
			ID:      order.ID,
			Status:  "Shipped",
			Version: order.Version,
		}
		jsonValue, _ := json.Marshal(orderUpdate)
		req, _ := http.NewRequest("PUT", "/orders/"+strconv.Itoa(int(order.ID)), bytes.NewBuffer(jsonValue))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Order updated successfully", response.Message)
	})

	db.Exec("DELETE FROM orders")
}

func TestSoftDeleteOrder(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test_order.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Order{})

	r := SetupRouter()
	r.Use(middleware.CORSMiddleware())

	ns := &utils.NotificationService{}

	routes.Routers(r, db, ns)

	// Seed an order for testing
	order := model.Order{AccountID: 1, CustomerID: 1, Quantity: 1, ProductID: 1, Status: "Pending"}
	db.Create(&order)

	t.Run("SoftDeleteOrderSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)
		req, _ := http.NewRequest("DELETE", "/orders/"+strconv.Itoa(int(order.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Order deleted successfully", response.Message)
	})

	db.Exec("DELETE FROM orders")
}

func TestHardDeleteOrder(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test_order.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Order{})

	r := SetupRouter()
	r.Use(middleware.CORSMiddleware())

	ns := &utils.NotificationService{}

	routes.Routers(r, db, ns)

	// Seed an order for testing
	order := model.Order{AccountID: 1, CustomerID: 1, Quantity: 1, ProductID: 1, Status: "Pending"}
	db.Create(&order)

	t.Run("HardDeleteOrderSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)
		req, _ := http.NewRequest("DELETE", "/orders/hard/"+strconv.Itoa(int(order.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Order deleted permanently", response.Message)
	})

	db.Exec("DELETE FROM orders")
}

func TestRecoverOrder(t *testing.T) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8081") // Set the USER_SERVICE_URL environment variable
	os.Setenv("TOKEN_SECRET", "test_secret")               // Set the TOKEN_SECRET environment variable

	db, err := gorm.Open(sqlite.Open("test_order.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.Order{}, &model.User{}, &model.Role{})

	// Create a role and user for testing
	role := model.Role{
		ID: 1,
		// Other fields if required
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

	// Mock user service
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
		r.Run(":8081")
	}()

	r := SetupRouter()
	r.Use(middleware.CORSMiddleware())

	ns := &utils.NotificationService{}

	routes.Routers(r, db, ns)

	// Seed a soft-deleted order for testing
	order := model.Order{AccountID: 1, CustomerID: 1, Quantity: 1, ProductID: 1, Status: "Pending"}
	db.Create(&order)
	db.Delete(&order)

	t.Run("RecoverOrderSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)
		req, _ := http.NewRequest("POST", "/orders/recover/"+strconv.Itoa(int(order.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Order recovered successfully", response.Message)
	})

	db.Exec("DELETE FROM orders")
}

func TestCancelOrder(t *testing.T) {
	os.Setenv("USER_SERVICE_URL", "http://localhost:8080") // Set the USER_SERVICE_URL environment variable
	os.Setenv("TOKEN_SECRET", "test_secret")               // Set the TOKEN_SECRET environment variable
	db, err := gorm.Open(sqlite.Open("test_order.db"), &gorm.Config{})
	assert.NoError(t, err)

	db.AutoMigrate(&model.User{}, &model.Role{}, &model.Order{})

	// Create a role for the user
	role := model.Role{
		ID: 1,
		// Populate other fields as necessary
	}
	db.Create(&role)

	// Create a user for testing
	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
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

	// Mock user service
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

	// Create a mock email sender
	mockEmailSender := new(MockEmailSender)
	ns := utils.NewNotificationService(mockEmailSender)

	r := SetupRouter()
	r.Use(middleware.CORSMiddleware())
	routes.Routers(r, db, ns)

	// Seed an order for testing
	order := model.Order{AccountID: 1, CustomerID: 1, Quantity: 1, ProductID: 1, Status: "Pending"}
	db.Create(&order)

	// Set up the expected email details
	expectedEmail := "customer@example.com"
	expectedSubject := "Your Order Has Been Cancelled"
	expectedBody := fmt.Sprintf("Your order with Order ID %d has been cancelled.", order.ID)

	// Set up the mock to expect the email with the correct order ID
	mockEmailSender.On("SendEmail", expectedEmail, expectedSubject, expectedBody).Return(nil)

	t.Run("CancelOrderSuccess", func(t *testing.T) {
		token := createTestToken(1, 1)
		req, _ := http.NewRequest("POST", "/orders/cancel/"+strconv.Itoa(int(order.ID)), nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var response model.SuccessResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Order cancelled successfully", response.Message)

		// Assert that the mock email sender was called with the expected arguments
		mockEmailSender.AssertExpectations(t)
	})

	db.Exec("DELETE FROM orders")
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM roles")
}
