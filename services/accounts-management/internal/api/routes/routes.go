package routes

import (
	"accounts-management/internal/api/handlers"
	"accounts-management/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	r.Use(middleware.CORSMiddleware())

	accounts := r.Group("/accounts")
	accounts.POST("/", handlers.CreateAccount(db))
	accounts.GET("/", handlers.GetAccounts(db))

	r.Use(middleware.AuthMiddleware(db))
	accounts.PUT("/:id", handlers.UpdateAccount(db))
	accounts.DELETE("/:id", handlers.SoftDeleteAccount(db))
	accounts.DELETE("hard/:id", handlers.HardDeleteAccount(db))
	accounts.POST("/:id/recover", handlers.RecoverAccount(db))
}
