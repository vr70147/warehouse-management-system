package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	users := r.Group("/products")

	users.POST("/users/signup", Signup(db))
	users.POST("/users/login", Login(db))
	// users.GET("/users/validate", RequireAuth, Validate(db))
	users.POST("/users/logout", Logout(db))
	users.GET("/users", RequireAuth, RequireAdmin, GetUsers(db))
	users.PUT("/users/:id", RequireAuth, RequireAdmin, UpdateUser(db))
	users.PATCH("/users/recover/:id", RequireAuth, RequireAdmin, RecoverUser(db))
	users.POST("/roles", CreateRole(db))
	users.PUT("/roles/:id", RequireAuth, RequireAdmin, UpdateRole(db))
	users.GET("/roles", RequireAuth, RequireAdmin, GetRoles(db))
	users.DELETE("/roles/:id", RequireAuth, RequireAdmin, DeleteRole(db))
	users.PATCH("/roles/recover/:id", RequireAuth, RequireAdmin, RecoverRole(db))
	users.POST("/departments", RequireAdmin, CreateDepartment(db))
	users.PUT("/departments/:id", RequireAuth, RequireAdmin, UpdateDepartment(db))
	users.GET("/departments", RequireAuth, RequireAdmin, GetDepartment(db))
	users.DELETE("/departments/:id", RequireAuth, RequireAdmin, DeleteDepartment(db))
	users.PATCH("/departments/recover/:id", RequireAuth, RequireAdmin, RecoverDepartment(db))
	users.GET("/departments/users", RequireAuth, RequireAdmin, GetUsersByDepartment(db))
}
