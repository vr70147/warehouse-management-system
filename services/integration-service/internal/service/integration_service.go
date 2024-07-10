package service

import (
	"fmt"
	"integration-service/internal/client"
	"integration-service/internal/config"
	"integration-service/internal/model"
	"integration-service/kafka"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.InfoLevel)
}

const maxRetries = 3
const retryDelay = 2 * time.Second

func CreateOrder(c *gin.Context) {
	var orderRequest model.OrderRequest
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Invalid request payload")
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	var orderResponse model.OrderResponse
	var err error
	for i := 0; i < maxRetries; i++ {
		orderResponse, err = client.CallOrderService(orderRequest, config.OrderServiceURL)
		if err == nil {
			break
		}
		log.WithFields(log.Fields{
			"attempt": i + 1,
			"error":   err,
		}).Error("Failed to create order")
		time.Sleep(retryDelay)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to create order",
			Details: err.Error(),
		})
		return
	}

	orderEvent := fmt.Sprintf(`{"order_id": %d,"order_status": "%s"}`, orderResponse.OrderID, orderResponse.Status)
	if err := kafka.ProduceMessage(config.OrderEventsTopic, orderEvent); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Failed to publish order event")
	}
	c.JSON(http.StatusOK, orderResponse)
}

func GetOrder(c *gin.Context) {
	orderID := c.Param("id")

	// Assuming an implementation to get order details from the order service
	orderResponse, err := client.CallGetOrderService(orderID, config.OrderServiceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get order",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, orderResponse)
}

func GetSalesReports(c *gin.Context) {
	var request model.SalesReportRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	response, err := client.CallSalesReportService(request, config.ReportingAnalyticsURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Failed to get sales report",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func GetInventoryLevel(c *gin.Context) {
	var request model.InventoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	response, err := client.CallInventoryService(request, config.ReportingAnalyticsURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get inventory level",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func GetShippingStatus(c *gin.Context) {
	var request model.ShippingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	response, err := client.CallShippingService(request, config.ReportingAnalyticsURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get shipping status",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func GetUserActivities(c *gin.Context) {
	var request model.UserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	response, err := client.CallUsersService(request, config.ReportingAnalyticsURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to get user activities",
			Details: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}
