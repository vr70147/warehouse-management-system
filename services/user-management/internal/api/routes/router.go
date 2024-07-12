package routes

import (
	"user-management/internal/api/handlers"
	"user-management/internal/middleware"
	"user-management/internal/model"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	users := r.Group("/users")

	users.POST("/signup", middleware.RequirePermission(model.PermissionManager), handlers.Signup(db))
	users.POST("/login", handlers.Login(db))
	users.POST("/logout", handlers.Logout(db))

	users.Use(middleware.RequireAuth(db))

	users.GET("/", middleware.RequirePermission(model.PermissionManager), handlers.GetUsers(db))
	users.PUT("/:id", middleware.RequirePermission(model.PermissionManager), handlers.UpdateUser(db))
	users.PATCH("/recover/:id", middleware.RequirePermission(model.PermissionManager), handlers.RecoverUser(db))
	users.DELETE("/:id", middleware.RequirePermission(model.PermissionManager), handlers.SoftDeleteUser(db))
	users.DELETE("/hard/:id", middleware.RequirePermission(model.PermissionManager), handlers.HardDeleteUser(db))

	users.POST("/roles", middleware.RequirePermission(model.PermissionManager), handlers.CreateRole(db))
	users.PUT("/roles/:id", middleware.RequirePermission(model.PermissionManager), handlers.UpdateRole(db))
	users.GET("/roles", middleware.RequirePermission(model.PermissionManager), handlers.GetRoles(db))
	users.DELETE("/roles/:id", middleware.RequirePermission(model.PermissionManager), handlers.SoftDeleteRole(db))
	users.DELETE("/roles/hard/:id", middleware.RequirePermission(model.PermissionManager), handlers.HardDeleteRole(db))
	users.PATCH("/roles/recover/:id", middleware.RequirePermission(model.PermissionManager), handlers.RecoverRole(db))

	users.POST("/departments", middleware.RequirePermission(model.PermissionManager), handlers.CreateDepartment(db))
	users.PUT("/departments/:id", middleware.RequirePermission(model.PermissionManager), handlers.UpdateDepartment(db))
	users.GET("/departments", middleware.RequirePermission(model.PermissionManager), handlers.GetDepartment(db))
	users.DELETE("/departments/:id", middleware.RequirePermission(model.PermissionManager), handlers.SoftDeleteDepartment(db))
	users.DELETE("/departments/hard/:id", middleware.RequirePermission(model.PermissionManager), handlers.HardDeleteDepartment(db))
	users.PATCH("/departments/recover/:id", middleware.RequirePermission(model.PermissionManager), handlers.RecoverDepartment(db))
	users.GET("/departments/users", middleware.RequirePermission(model.PermissionManager), handlers.GetUsersByDepartment(db))
}
