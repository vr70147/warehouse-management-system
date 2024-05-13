package api

import "github.com/gin-gonic/gin"

func Routers() {
	router := gin.Default()

	router.POST("/signup", Signup)
	router.POST("/login", Login)
	router.GET("/validate", RequireAuth, Validate)
	router.POST("/logout", Logout)
	router.GET("/users", RequireAuth, GetUsers)
	router.PUT("/users/:id", RequireAuth, UpdateUser)
	router.PATCH("/users/recover/:id", RequireAuth, RecoverUser)
	router.POST("/roles", CreateRole)
	router.PUT("/roles/:id", RequireAuth, UpdateRole)
	router.GET("/roles", RequireAuth, GetRoles)
	router.DELETE("/roles/:id", RequireAuth, DeleteRole)
	router.PATCH("/roles/recover/:id", RequireAuth, RecoverRole)
	router.POST("/departments", CreateDepartment)
	router.PUT("/departments/:id", RequireAuth, UpdateDepartment)
	router.GET("/departments", RequireAuth, GetDepartment)
	router.DELETE("/departments/:id", RequireAuth, DeleteDepartment)
	router.PATCH("/departments/recover/:id", RequireAuth, RecoverDepartment)
	router.GET("/departments/users", RequireAuth, GetUsersByDepartment)

	router.Run()
}
