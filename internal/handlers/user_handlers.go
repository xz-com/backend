package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/omega/notes-app/internal/auth"
	"github.com/omega/notes-app/internal/database"
	"github.com/omega/notes-app/internal/models"
)

// RegisterRequest представляет данные для регистрации пользователя
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest представляет данные для входа пользователя
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register обрабатывает запрос на регистрацию нового пользователя
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, существует ли пользователь с таким email
	var existingUser models.User
	result := database.GetDB().Where("email = ?", req.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пользователь с таким email уже существует"})
		return
	}

	// Проверяем, существует ли пользователь с таким username
	result = database.GetDB().Where("username = ?", req.Username).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "пользователь с таким именем уже существует"})
		return
	}

	// Создаем нового пользователя
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// Сохраняем пользователя в базе данных
	if err := database.GetDB().Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании пользователя"})
		return
	}

	// Генерируем JWT-токен
	token, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании токена"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "пользователь успешно зарегистрирован",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
		"token": token,
	})
}

// Login обрабатывает запрос на вход пользователя
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ищем пользователя по email
	var user models.User
	result := database.GetDB().Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный email или пароль"})
		return
	}

	// Проверяем пароль
	if err := user.ValidatePassword(req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный email или пароль"})
		return
	}

	// Генерируем JWT-токен
	token, err := auth.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании токена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "вход выполнен успешно",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
		"token": token,
	})
}

// GetProfile возвращает профиль текущего пользователя
func GetProfile(c *gin.Context) {
	// Получаем пользователя из контекста (установленного middleware)
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
} 