package handlers

import (
	"net/http"
	"shipping-receiving/internal/model"
	"shipping-receiving/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateShipping godoc
// @Summary Create a new Shipping
// @Description Create a new Shipping
// @Tags Shippings
// @Accept json
// @Produce json
// @Param Shipping body model.Shipping true "Shipping"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /shippings [post]
func CreateShipping(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		var shipping model.Shipping
		if err := c.ShouldBindJSON(&shipping); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		shipping.AccountID = accountID.(uint)

		if result := db.Create(&shipping); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Shipping created successfully", Data: shipping})
	}
}

// GetShippings godoc
// @Summary Get Shippings
// @Description Get Shippings
// @Tags Shippings
// @Produce json
// @Param id query string false "Shipping ID"
// @Param status query string false "Shipping Status"
// @Param receiver_id query string false "Receiver ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {array} model.Shipping
// @Failure 500 {object} model.ErrorResponse
// @Router /shippings [get]
func GetShippings(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Query("id")

		if id != "" {
			// Fetch a single Shipping by ID
			var shipping model.Shipping
			if result := db.Where("id = ? AND account_id = ?", id, accountID).First(&shipping); result.Error != nil {
				c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Shipping not found"})
				return
			}
			c.JSON(http.StatusOK, model.SuccessResponse{Message: "Shipping retrieved successfully", Data: shipping})
			return
		}

		// Fetch list of Shippings with optional query parameters
		var shippings []model.Shipping
		query := db.Where("account_id = ?", accountID)

		if status := c.Query("status"); status != "" {
			query = query.Where("status = ?", status)
		}

		if receiverID := c.Query("receiver_id"); receiverID != "" {
			query = query.Where("receiver_id = ?", receiverID)
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

		if result := query.Find(&shippings); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponses{Message: "Shippings retrieved successfully", Data: shippings})
	}
}

// UpdateShipping godoc
// @Summary Update a Shipping
// @Description Update a Shipping
// @Tags Shippings
// @Accept json
// @Produce json
// @Param id path string true "Shipping ID"
// @Param Shipping body model.Shipping true "Shipping"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /shippings/{id} [put]
func UpdateShipping(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		var shipping model.Shipping
		if result := db.Where("id = ? AND account_id = ?", id, accountID).First(&shipping); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Shipping not found"})
			return
		}

		if err := c.ShouldBindJSON(&shipping); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Invalid request data"})
			return
		}

		shipping.AccountID = accountID.(uint)

		if result := db.Save(&shipping); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Shipping updated successfully", Data: shipping})
	}
}

// SoftDeleteShipping godoc
// @Summary Soft delete a Shipping
// @Description Soft delete a Shipping
// @Tags Shippings
// @Produce json
// @Param id path string true "Shipping ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /shippings/{id} [delete]
func SoftDeleteShipping(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		var shipping model.Shipping
		if result := db.Where("id = ? AND account_id = ?", id, accountID).First(&shipping); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Shipping not found"})
			return
		}

		if result := db.Delete(&shipping); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Shipping soft deleted successfully", Data: shipping})
	}
}

// HardDeleteShipping godoc
// @Summary Hard delete a Shipping
// @Description Hard delete a Shipping
// @Tags Shippings
// @Produce json
// @Param id path string true "Shipping ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /shippings/hard/{id} [delete]
func HardDeleteShipping(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		if result := db.Unscoped().Where("id = ? AND account_id = ?", id, accountID).Delete(&model.Shipping{}); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Shipping hard deleted successfully"})
	}
}

// RecoverShipping godoc
// @Summary Recover a soft-deleted Shipping
// @Description Recover a soft-deleted Shipping
// @Tags Shippings
// @Produce json
// @Param id path string true "Shipping ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /shippings/recover/{id} [put]
func RecoverShipping(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		id := c.Param("id")
		var shipping model.Shipping

		// Find the soft-deleted Shipping
		if result := db.Unscoped().Where("id = ? AND account_id = ?", id, accountID).First(&shipping); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Shipping not found"})
			return
		}

		// Recover the Shipping by setting DeletedAt to NULL
		if result := db.Model(&shipping).Update("deleted_at", nil); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Shipping recovered successfully", Data: shipping})
	}
}

// DeliverShipping godoc
// @Summary Deliver a Shipping
// @Description Mark a Shipping as delivered and send a notification email
// @Tags Shippings
// @Accept json
// @Produce json
// @Param id path string true "Shipping ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /shippings/deliver/{id} [post]
func DeliverShipping(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountID, exists := c.Get("account_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, model.ErrorResponse{Error: "Account ID not found"})
			return
		}

		shippingID := c.Param("id")
		var shipping model.Shipping

		if result := db.Where("id = ? AND accound_id = ?", shippingID, accountID).First(&shipping); result.Error != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Shipping not found"})
			return
		}

		shipping.Status = "delivered"
		if result := db.Save(&shipping); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to update shipping status"})
			return
		}

		// Send notification email
		emailSubject := "Your Order Has Been Delivered"
		emailBody := "Your order with Shipping ID " + strconv.Itoa(int(shipping.ID)) + " has been delivered successfully."
		if err := utils.SendEmail("customer@example.com", emailSubject, emailBody); err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to send notification email"})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Shipping delivered successfully"})
	}
}
