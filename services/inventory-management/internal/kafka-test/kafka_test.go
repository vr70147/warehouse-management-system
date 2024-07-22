package kafka_test_test

import (
	"context"
	"encoding/json"
	"inventory-management/internal/initializers"
	kafka "inventory-management/internal/kafka"
	"inventory-management/internal/model"
	"os"
	"testing"
	"time"

	kafkago "github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mock Kafka Service
type MockKafkaService struct {
	messages []kafkago.Message
}

func (m *MockKafkaService) ReadMessage(ctx context.Context) (kafkago.Message, error) {
	if len(m.messages) == 0 {
		return kafkago.Message{}, kafka.ErrEmptyQueue
	}
	msg := m.messages[0]
	m.messages = m.messages[1:]
	return msg, nil
}

func (m *MockKafkaService) WriteMessages(ctx context.Context, msgs ...kafkago.Message) error {
	m.messages = append(m.messages, msgs...)
	return nil
}

func setupTestDatabase() (*gorm.DB, model.User) {
	db, err := gorm.Open(sqlite.Open("test_inventory.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(&model.Product{}, &model.Stock{}, &model.User{}, &model.Role{}, &model.Supplier{}, &model.Order{})

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

	initializers.DB = db

	return db, testUser
}

func TestConsumerOrderEvents(t *testing.T) {
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	os.Setenv("NEW_ORDER_EVENTS_TOPIC", "order-events")

	db, _ := setupTestDatabase()
	defer db.Exec("DELETE FROM orders")

	// Mock Kafka service with a sample message
	mockKafkaService := &MockKafkaService{
		messages: []kafkago.Message{
			{
				Value: []byte(`{"id":1,"quantity":5,"customer_id":1,"product_id":1,"status":"Pending","account_id":1}`),
			},
		},
	}

	// Create a sample order in the database
	order := model.Order{
		ID:         1,
		Quantity:   5,
		CustomerID: 1,
		ProductID:  1,
		Status:     "Pending",
		AccountID:  1,
	}
	db.Create(&order)

	// Override the Kafka service
	kafka.KafkaSvc = mockKafkaService

	// Run the consumer
	go kafka.ConsumerOrderEvents()

	// Give the consumer some time to process the message
	time.Sleep(1 * time.Second)

	// Check the order status in the database
	var updatedOrder model.Order
	db.First(&updatedOrder, order.ID)
	assert.Equal(t, "Ready", updatedOrder.Status)

	// Check the Kafka service for the published message
	assert.Len(t, mockKafkaService.messages, 1)

	var publishedMessage map[string]interface{}
	json.Unmarshal(mockKafkaService.messages[0].Value, &publishedMessage)
	assert.Equal(t, float64(order.ID), publishedMessage["orderID"].(float64))
	assert.Equal(t, "Ready", publishedMessage["status"])
}

func TestPublishOrderStatus(t *testing.T) {
	os.Setenv("KAFKA_BROKERS", "localhost:9092")
	os.Setenv("NEW_ORDER_EVENTS_TOPIC", "order-events")

	// Mock Kafka service
	mockKafkaService := &MockKafkaService{}

	// Override the Kafka service
	kafka.KafkaSvc = mockKafkaService

	// Publish an order status
	err := kafka.PublishOrderStatus(1, "Ready")
	assert.NoError(t, err)

	// Check the Kafka service for the published message
	assert.Len(t, mockKafkaService.messages, 1)

	var publishedMessage map[string]interface{}
	json.Unmarshal(mockKafkaService.messages[0].Value, &publishedMessage)
	assert.Equal(t, float64(1), publishedMessage["orderID"].(float64))
	assert.Equal(t, "Ready", publishedMessage["status"])
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
