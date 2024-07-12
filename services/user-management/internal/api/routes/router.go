package routes

import (
	"user-management/internal/api/handlers"
	"user-management/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	users := r.Group("/users")

	users.POST("/signup", handlers.Signup(db))
	users.POST("/login", handlers.Login(db))
	users.POST("/logout", handlers.Logout(db))

	// Apply RequireAuth middleware to all routes that need authentication
	authenticated := users.Group("/")
	authenticated.Use(middleware.RequireAuth(db))

	authenticated.GET("/", middleware.RequireAdmin(), handlers.GetUsers(db))
	authenticated.PUT("/:id", middleware.RequireAdmin(), handlers.UpdateUser(db))
	authenticated.PATCH("/recover/:id", middleware.RequireAdmin(), handlers.RecoverUser(db))
	authenticated.POST("/roles", middleware.RequireAdmin(), handlers.CreateRole(db))
	authenticated.PUT("/roles/:id", middleware.RequireAdmin(), handlers.UpdateRole(db))
	authenticated.GET("/roles", middleware.RequireAdmin(), handlers.GetRoles(db))
	authenticated.DELETE("/roles/:id/soft", middleware.RequireAdmin(), handlers.SoftDeleteRole(db))
	authenticated.DELETE("/roles/:id/hard", middleware.RequireAdmin(), handlers.HardDeleteRole(db))
	authenticated.PATCH("/roles/recover/:id", middleware.RequireAdmin(), handlers.RecoverRole(db))
	authenticated.POST("/departments", middleware.RequireAdmin(), handlers.CreateDepartment(db))
	authenticated.PUT("/departments/:id", middleware.RequireAdmin(), handlers.UpdateDepartment(db))
	authenticated.GET("/departments", middleware.RequireAdmin(), handlers.GetDepartment(db))
	authenticated.DELETE("/departments/:id/soft", middleware.RequireAdmin(), handlers.SoftDeleteDepartment(db))
	authenticated.DELETE("/departments/:id/hard", middleware.RequireAdmin(), handlers.HardDeleteDepartment(db))
	authenticated.PATCH("/departments/recover/:id", middleware.RequireAdmin(), handlers.RecoverDepartment(db))
	authenticated.GET("/departments/users", middleware.RequireAdmin(), handlers.GetUsersByDepartment(db))
}
