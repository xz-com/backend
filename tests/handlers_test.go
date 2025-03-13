package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/omega/notes-app/internal/handlers"
	"github.com/omega/notes-app/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Мок для базы данных
type MockDB struct {
	users []models.User
	notes []models.Note
}

// Мок для GetDB
func mockGetDB() *gorm.DB {
	return nil // В тестах мы не будем использовать реальную базу данных
}

// Настройка тестового окружения
func setupTestRouter() (*gin.Engine, *MockDB) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	mockDB := &MockDB{
		users: []models.User{},
		notes: []models.Note{},
	}
	return router, mockDB
}

// Тест для регистрации пользователя
func TestRegisterHandlerMock(t *testing.T) {
	// Этот тест просто проверяет, что обработчик принимает запрос в правильном формате
	// Для полноценного тестирования нужно использовать моки или тестовую базу данных
	router := gin.Default()
	router.POST("/register", func(c *gin.Context) {
		var req handlers.RegisterRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// Имитируем успешную регистрацию
		c.JSON(http.StatusCreated, gin.H{
			"message": "пользователь успешно зарегистрирован",
			"user": gin.H{
				"id":       1,
				"username": req.Username,
				"email":    req.Email,
			},
			"token": "mock_token",
		})
	})

	// Создаем тестовые данные
	registerData := handlers.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(registerData)

	// Создаем тестовый запрос
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Выполняем запрос
	router.ServeHTTP(resp, req)

	// Проверяем результат
	assert.Equal(t, http.StatusCreated, resp.Code)

	// Парсим ответ
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Проверяем, что ответ содержит ожидаемые поля
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "user")
	assert.Contains(t, response, "token")
}

// Тест для входа пользователя
func TestLoginHandlerMock(t *testing.T) {
	router := gin.Default()
	router.POST("/login", func(c *gin.Context) {
		var req handlers.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		// Проверяем тестовые учетные данные
		if req.Email == "test@example.com" && req.Password == "password123" {
			c.JSON(http.StatusOK, gin.H{
				"message": "вход выполнен успешно",
				"user": gin.H{
					"id":       1,
					"username": "testuser",
					"email":    req.Email,
				},
				"token": "mock_token",
			})
			return
		}
		
		// Имитируем неудачный вход
		c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный email или пароль"})
	})

	// Создаем тестовые данные
	loginData := handlers.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonData, _ := json.Marshal(loginData)

	// Создаем тестовый запрос
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Выполняем запрос
	router.ServeHTTP(resp, req)

	// Проверяем результат (должен быть 200, так как мы имитируем успешный вход)
	assert.Equal(t, http.StatusOK, resp.Code)
	
	// Проверяем неверный пароль
	loginData.Password = "wrongpassword"
	jsonData, _ = json.Marshal(loginData)
	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusUnauthorized, resp.Code)
} 