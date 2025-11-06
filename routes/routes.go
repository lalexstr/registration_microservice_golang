package routes

import (
	"auth-service/handlers"
	"auth-service/middleware"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	authHandler := handlers.NewAuthHandler()
	userHandler := handlers.NewUserHandler()

	// public
	r.POST("/api/v1/register", authHandler.Register)
	r.POST("/api/v1/login", authHandler.Login)

	// protected
	auth := r.Group("/api/v1")
	auth.Use(middleware.JWTAuthMiddleware())

	{
		auth.GET("/me", userHandler.GetMe)
		// admin-only endpoints: we'll check inside handlers using role
		auth.GET("/users", userHandler.ListUsers)         // admin
		auth.POST("/users", userHandler.CreateUser)       // admin can create admin/user
		auth.GET("/users/:id", userHandler.GetUser)       // admin or owner
		auth.PUT("/users/:id", userHandler.UpdateUser)    // admin or owner
		auth.DELETE("/users/:id", userHandler.DeleteUser) // admin or owner
	}
}
