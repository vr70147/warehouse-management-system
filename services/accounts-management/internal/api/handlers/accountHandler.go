package handlers

import (
	"accounts-management/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

		c.JSON(http.StatusOK, account)
	}
}

func GetAccounts(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accounts []model.Account
		if err := db.Preload("Plan").Find(&accounts).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, accounts)
	}
}

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

		c.JSON(http.StatusOK, account)
	}
}

func SoftDeleteAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account model.Account
		if err := db.Where("id = ?", c.Param("id")).First(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		if err := db.Delete(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, account)
	}
}

func HardDeleteAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account model.Account
		if err := db.Where("id = ?", c.Param("id")).First(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		if err := db.Unscoped().Delete(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, account)
	}
}

func RecoverAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var account model.Account
		if err := db.Unscoped().Where("id = ?", c.Param("id")).First(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		account.DeletedAt = nil
		if err := db.Save(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, account)
	}
}
