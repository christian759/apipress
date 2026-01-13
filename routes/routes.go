package routes

import (
	"apipress/controllers"
	"apipress/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		// Public Post Routes
		api.GET("/posts", controllers.GetPosts)
		api.GET("/posts/:slug", controllers.GetPostBySlug)

		// Protected Post Routes
		posts := api.Group("/posts")
		posts.Use(middleware.AuthMiddleware())
		{
			posts.POST("", controllers.CreatePost)
			posts.PUT("/:id", controllers.UpdatePost)
			posts.DELETE("/:id", controllers.DeletePost)
		}
	}
}
