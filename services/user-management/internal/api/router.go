package api

import "github.com/gin-gonic/gin"

func Routers() {
	router := gin.Default()

	//initial Postgres
	router.POST("/signup", Signup)
	router.POST("/login", Login)
	router.GET("/validate", RequireAuth, Validate)
	router.POST("/logout", Logout)
	router.GET("/users", RequireAuth, GetUsers)
	router.PUT("/users/:id", RequireAuth, UpdateUser)
	router.POST("/roles", RequireAuth, CreateRole)
	router.PUT("/roles/:id", RequireAuth, UpdateRole)
	router.GET("/roles", RequireAuth, GetRoles)
	router.DELETE("/roles/:id", RequireAuth, DeleteRole)
	router.PATCH("/roles/recover/:id", RequireAuth, RecoverRole)
	router.Run()
}
