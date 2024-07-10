package api

import (
	"user-management/internal/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	users := r.Group("/products")

	users.POST("/users/signup", Signup(db))
	users.POST("/users/login", Login(db))
	// users.GET("/users/validate", middleware.RequireAuth, Validate(db))
	users.POST("/users/logout", Logout(db))
	users.GET("/users", middleware.RequireAuth, middleware.RequireAdmin, GetUsers(db))
	users.PUT("/users/:id", middleware.RequireAuth, middleware.RequireAdmin, UpdateUser(db))
	users.PATCH("/users/recover/:id", middleware.RequireAuth, middleware.RequireAdmin, RecoverUser(db))
	users.POST("/roles", CreateRole(db))
	users.PUT("/roles/:id", middleware.RequireAuth, middleware.RequireAdmin, UpdateRole(db))
	users.GET("/roles", middleware.RequireAuth, middleware.RequireAdmin, GetRoles(db))
	users.DELETE("/roles/:id", middleware.RequireAuth, middleware.RequireAdmin, DeleteRole(db))
	users.PATCH("/roles/recover/:id", middleware.RequireAuth, middleware.RequireAdmin, RecoverRole(db))
	users.POST("/departments", middleware.RequireAdmin, CreateDepartment(db))
	users.PUT("/departments/:id", middleware.RequireAuth, middleware.RequireAdmin, UpdateDepartment(db))
	users.GET("/departments", middleware.RequireAuth, middleware.RequireAdmin, GetDepartment(db))
	users.DELETE("/departments/:id", middleware.RequireAuth, middleware.RequireAdmin, DeleteDepartment(db))
	users.PATCH("/departments/recover/:id", middleware.RequireAuth, middleware.RequireAdmin, RecoverDepartment(db))
	users.GET("/departments/users", middleware.RequireAuth, middleware.RequireAdmin, GetUsersByDepartment(db))
}
