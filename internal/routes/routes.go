package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/omega/notes-app/internal/handlers"
	"github.com/omega/notes-app/internal/middleware"
)

// SetupRoutes настраивает все маршруты API
func SetupRoutes(router *gin.Engine) {
	// Группа маршрутов для API
	api := router.Group("/api")
	{
		// Маршруты для аутентификации (без middleware)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// Маршруты, требующие аутентификации
		user := api.Group("/user")
		user.Use(middleware.AuthMiddleware())
		{
			user.GET("/profile", handlers.GetProfile)
		}

		// Маршруты для заметок (требуют аутентификации)
		notes := api.Group("/notes")
		notes.Use(middleware.AuthMiddleware())
		{
			notes.POST("", handlers.CreateNote)
			notes.GET("", handlers.GetNotes)
			notes.GET("/:id", handlers.GetNote)
			notes.PUT("/:id", handlers.UpdateNote)
			notes.DELETE("/:id", handlers.DeleteNote)
		}
	}
} 