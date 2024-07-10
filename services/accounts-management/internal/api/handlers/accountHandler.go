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
		if err := db.Find(&accounts).Error; err != nil {
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

		if err := c.ShouldBindJSON(&account); err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: err.Error()})
			return
		}

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

		if err := db.Save(&account).Error; err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, account)
	}
}
