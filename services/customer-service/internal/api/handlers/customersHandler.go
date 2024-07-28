package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"customer-service/internal/model"
)

// GetCustomers godoc
// @Summary Get customers
// @Description Retrieve a list of customers with optional query parameters
// @Tags customers
// @Produce json
// @Param name query string false "Customer Name"
// @Param email query string false "Customer Email"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} model.SuccessResponses
// @Failure 401 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /customers [get]
func GetCustomers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var customers []model.Customer
		query := db.Where("account_id = ?", accountID)

		if name := c.Query("name"); name != "" {
			query = query.Where("name LIKE ?", "%"+name+"%")
		}

		if email := c.Query("email"); email != "" {
			query = query.Where("email = ?", email)
		}

		if limit := c.Query("limit"); limit != "" {
			if limitInt, err := strconv.Atoi(limit); err == nil {
				query = query.Limit(limitInt)
			}
		}

		if offset := c.Query("offset"); offset != "" {
			if offsetInt, err := strconv.Atoi(offset); err == nil {
				query = query.Offset(offsetInt)
			}
		}

		if result := query.Find(&customers); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponses{Message: "Customers found", Customers: customers})
	}
}

// GetCustomerByID godoc
// @Summary Get customer by ID
// @Description Retrieve customer details by customer ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} model.Customer
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /customers/{id} [get]
func GetCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		customerID := c.Param("id")

		var customer model.Customer
		if err := db.Where("id = ?", customerID).First(&customer).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Customer not found"})
			return
		}

		c.JSON(http.StatusOK, customer)
	}
}

// CreateCustomer godoc
// @Summary Create a new customer
// @Description Create a new customer in the system
// @Tags customers
// @Accept json
// @Produce json
// @Param body body model.Customer true "Customer data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Router /customers [post]
func CreateCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var customerRequest model.Customer
		if err := c.ShouldBindJSON(&customerRequest); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid customer data"})
			return
		}

		customer := model.Customer{
			AccountID:  accountID.(uint),
			Name:       customerRequest.Name,
			Email:      customerRequest.Email,
			Phone:      customerRequest.Phone,
			Address:    customerRequest.Address,
			City:       customerRequest.City,
			State:      customerRequest.State,
			PostalCode: customerRequest.PostalCode,
			Country:    customerRequest.Country,
		}

		if err := db.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to create customer"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Customer created successfully", Customer: customer})
	}
}

// UpdateCustomer godoc
// @Summary Update a customer
// @Description Update a customer by ID
// @Tags customers
// @Accept json
// @Produce json
// @Param id path int true "Customer ID"
// @Param body body model.Customer true "Customer data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Router /customers/{id} [put]
func UpdateCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var customerUpdate model.Customer
		if err := c.ShouldBindJSON(&customerUpdate); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid customer data"})
			return
		}

		var currentCustomer model.Customer
		if err := db.Where("id = ? AND account_id = ?", c.Param("id"), accountID).First(&currentCustomer).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Customer not found"})
			return
		}

		if err := db.Model(&currentCustomer).Updates(customerUpdate).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to update customer"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Customer updated successfully", Customer: currentCustomer})
	}
}

// SoftDeleteCustomer godoc
// @Summary Soft delete a customer
// @Description Soft delete a customer by ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /customers/{id} [delete]
func SoftDeleteCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		if result := db.Where("id = ? AND account_id = ?", id, accountID).Delete(&model.Customer{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete customer"})
			return
		}
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Customer deleted successfully"})
	}
}

// HardDeleteCustomer godoc
// @Summary Hard delete a customer
// @Description Permanently delete a customer by ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /customers/{id}/hard [delete]
func HardDeleteCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		if result := db.Unscoped().Where("id = ? AND account_id = ?", id, accountID).Delete(&model.Customer{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to delete customer"})
			return
		}
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Customer deleted permanently"})
	}
}

// RecoverCustomer godoc
// @Summary Recover a deleted customer
// @Description Recover a soft-deleted customer by ID
// @Tags customers
// @Produce json
// @Param id path int true "Customer ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /customers/{id}/recover [post]
func RecoverCustomer(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		var customer model.Customer

		if result := db.Unscoped().Where("id = ? AND account_id = ?", id, accountID).First(&customer); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Customer not found"})
			return
		}

		if result := db.Model(&customer).Update("DeletedAt", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to recover customer"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Customer recovered successfully"})
	}
}
