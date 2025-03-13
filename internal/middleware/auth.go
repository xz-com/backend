package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/omega/notes-app/internal/auth"
	"github.com/omega/notes-app/internal/database"
	"github.com/omega/notes-app/internal/models"
)

// AuthMiddleware проверяет JWT-токен и аутентифицирует пользователя
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Получаем заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			c.Abort()
			return
		}

		// Проверяем формат заголовка (Bearer token)
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный формат токена"})
			c.Abort()
			return
		}

		// Валидируем токен
		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "недействительный токен: " + err.Error()})
			c.Abort()
			return
		}

		// Получаем пользователя из базы данных
		var user models.User
		result := database.GetDB().First(&user, claims.UserID)
		if result.Error != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не найден"})
			c.Abort()
			return
		}

		// Устанавливаем пользователя в контекст
		c.Set("user", user)
		c.Set("user_id", claims.UserID)

		c.Next()
	}
} 