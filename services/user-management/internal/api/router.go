package api

import "github.com/gin-gonic/gin"

func Routers() {
	router := gin.Default()

	router.POST("/signup", Signup)
	router.POST("/login", Login)
	router.GET("/validate", RequireAuth, Validate)
	router.POST("/logout", Logout)
	router.GET("/users", RequireAuth, RequireAdmin, GetUsers)
	router.PUT("/users/:id", RequireAuth, RequireAdmin, UpdateUser)
	router.PATCH("/users/recover/:id", RequireAuth, RequireAdmin, RecoverUser)
	router.POST("/roles", RequireAdmin, CreateRole)
	router.PUT("/roles/:id", RequireAuth, RequireAdmin, UpdateRole)
	router.GET("/roles", RequireAuth, RequireAdmin, GetRoles)
	router.DELETE("/roles/:id", RequireAuth, RequireAdmin, DeleteRole)
	router.PATCH("/roles/recover/:id", RequireAuth, RequireAdmin, RecoverRole)
	router.POST("/departments", RequireAdmin, CreateDepartment)
	router.PUT("/departments/:id", RequireAuth, RequireAdmin, UpdateDepartment)
	router.GET("/departments", RequireAuth, RequireAdmin, GetDepartment)
	router.DELETE("/departments/:id", RequireAuth, RequireAdmin, DeleteDepartment)
	router.PATCH("/departments/recover/:id", RequireAuth, RequireAdmin, RecoverDepartment)
	router.GET("/departments/users", RequireAuth, RequireAdmin, GetUsersByDepartment)

	router.Run()
}
