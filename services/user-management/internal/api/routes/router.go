package routes

import (
	"user-management/internal/api/handlers"
	"user-management/internal/middleware"
	"user-management/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Routers(r *gin.Engine, db *gorm.DB) {
	users := r.Group("/users")

	users.POST("/signup", handlers.Signup(db, &utils.NotificationService{}))
	users.POST("/login", handlers.Login(db, &utils.NotificationService{}))
	users.POST("/logout", handlers.Logout(db))

	users.Use(middleware.RequireAuth(db))

	users.GET("/", handlers.GetUsers(db))
	users.PUT("/:id", handlers.UpdateUser(db))
	users.PATCH("/recover/:id", handlers.RecoverUser(db))
	users.DELETE("/:id", handlers.SoftDeleteUser(db))
	users.DELETE("/hard/:id", handlers.HardDeleteUser(db))
	users.POST("/users/change-password/:id", handlers.ChangePassword(db, &utils.NotificationService{}))

	users.POST("/roles", handlers.CreateRole(db))
	users.PUT("/roles/:id", handlers.UpdateRole(db))
	users.GET("/roles", handlers.GetRoles(db))
	users.DELETE("/roles/:id", handlers.SoftDeleteRole(db))
	users.DELETE("/roles/hard/:id", handlers.HardDeleteRole(db))
	users.PATCH("/roles/recover/:id", handlers.RecoverRole(db))

	users.POST("/departments", handlers.CreateDepartment(db))
	users.PUT("/departments/:id", handlers.UpdateDepartment(db))
	users.GET("/departments", handlers.GetDepartments(db))
	users.DELETE("/departments/:id", handlers.SoftDeleteDepartment(db))
	users.DELETE("/departments/hard/:id", handlers.HardDeleteDepartment(db))
	users.PATCH("/departments/recover/:id", handlers.RecoverDepartment(db))
	users.GET("/departments/users", handlers.GetUsersByDepartment(db))
}
