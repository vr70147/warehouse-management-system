package handlers

import (
	"accounts-management/internal/model"
	"accounts-management/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateAccount godoc
// @Summary Create a new account
// @Description Create a new account
// @Tags accounts
// @Accept json
// @Produce json
// @Param body body model.Account true "Account data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /accounts [post]
func CreateAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account model.Account
		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		if err := db.Create(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		// Log the created account ID
		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Account created successfully", Data: account})
	}
}

// GetAccounts godoc
// @Summary Get all accounts
// @Description Retrieve all accounts
// @Tags accounts
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /accounts [get]
func GetAccounts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accounts []model.Account
		if err := db.Preload("Plan").Find(&accounts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Accounts retrieved successfully", Data: accounts})
	}
}

// UpdateAccount godoc
// @Summary Update an account
// @Description Update account by ID
// @Tags accounts
// @Accept json
// @Produce json
// @Param id path int true "Account ID"
// @Param body body model.Account true "Account data"
// @Success 200 {object} model.SuccessResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /accounts/{id} [put]
func UpdateAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account model.Account
		if err := db.Where("id = ?", c.Param("id")).First(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		var updateData model.Account
		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

		account.Email = updateData.Email
		account.Name = updateData.Name
		account.PhoneNumber = updateData.PhoneNumber
		account.CompanyName = updateData.CompanyName
		account.Address = updateData.Address
		account.City = updateData.City
		account.State = updateData.State
		account.PostalCode = updateData.PostalCode
		account.Country = updateData.Country
		account.BillingEmail = updateData.BillingEmail
		account.BillingAddress = updateData.BillingAddress
		account.PlanID = updateData.PlanID
		account.IsActive = updateData.IsActive
		account.Metadata = updateData.Metadata
		account.Preferences = updateData.Preferences

		if err := db.Save(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Account updated successfully", Data: account})
	}
}

// SoftDeleteAccount godoc
// @Summary Soft delete an account
// @Description Soft delete account by ID
// @Tags accounts
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /accounts/{id} [delete]
func SoftDeleteAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account model.Account
		if err := db.Where("id = ?", c.Param("id")).First(&account).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: err.Error()})
			return
		}

		if err := db.Delete(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Account soft deleted successfully", Data: account})
	}
}

// HardDeleteAccount godoc
// @Summary Hard delete an account
// @Description Hard delete account by ID
// @Tags accounts
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /accounts/hard/{id} [delete]
func HardDeleteAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account model.Account
		if err := db.Where("id = ?", c.Param("id")).First(&account).Error; err != nil {
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: err.Error()})
			return
		}

		if err := db.Unscoped().Delete(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Account hard deleted successfully", Data: account})
	}
}

// RecoverAccount godoc
// @Summary Recover a deleted account
// @Description Recover a soft-deleted account by ID
// @Tags accounts
// @Produce json
// @Param id path int true "Account ID"
// @Success 200 {object} model.SuccessResponse
// @Failure 404 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /accounts/{id}/recover [post]
func RecoverAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account model.Account
		id := c.Param("id")
		if err := db.Unscoped().Where("id = ?", id).First(&account).Error; err != nil {
			fmt.Println("Error finding account:", id)
			c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "Account not found"})
			return
		}

		if !account.DeletedAt.Valid {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "Account is not deleted"})
			return
		}

		account.DeletedAt = gorm.DeletedAt{}
		if err := db.Save(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Account recovered successfully", Data: account})
	}
}

// CheckSubscriptionRenewal godoc
// @Summary Check subscription renewal
// @Description Check if account subscriptions are nearing expiration and send reminder emails
// @Tags accounts
// @Produce json
// @Success 200 {object} model.SuccessResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /accounts/check-renewal [get]
func CheckSubscription(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accounts []model.Account
		now := time.Now()
		renewalThreshold := now.AddDate(0, 0, 7)

		if result := db.Where("subscription_expires_at <= ?", renewalThreshold).Find(&accounts); result.Error != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to retrieve accounts"})
			return
		}

		for _, account := range accounts {
			// Send email notification
			emailSubject := "Subscription Renewal Reminder"
			emailBody := "Dear " + account.Name + ",\n\nYour subscription is set to expire on " + account.SubscriptionExpiresAt.Format("2006-01-02") + ". Please renew your subscription to avoid service interruption."
			if err := utils.SendEmail(account.Email, emailSubject, emailBody); err != nil {
				c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "Failed to send notification email"})
				return
			}
		}

		c.JSON(http.StatusOK, model.SuccessResponse{Message: "Subscription renewal reminders sent successfully"})
	}
}
