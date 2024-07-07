package service

import (
	"fmt"
	"integration-service/internal/client"
	"integration-service/internal/config"
	"integration-service/internal/model"
	"integration-service/kafka"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var orderRequest model.OrderRequest
	if err := c.ShouldBindBodyWithJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request payload",
			Details: err.Error(),
		})
		return
	}

	orderResponse, err := client.CallOrderService(orderRequest, config.OrderServiceURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "failed to create order")
		return
	}

	orderEvent := fmt.Sprintf(`{"order_id": %d,"order_status": "%s"}`, orderResponse.OrderID, orderResponse.Status)
	if err := kafka.ProduceMessage(config.OrderEventsTopic, orderEvent); err != nil {
		log.Printf("failed to publish order event: %v", err)
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
