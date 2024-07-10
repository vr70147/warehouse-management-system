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

		c.JSON(http.StatusOK, account)
	}
}
