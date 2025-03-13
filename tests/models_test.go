package tests

import (
	"testing"

	"github.com/omega/notes-app/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestUserPasswordHashing(t *testing.T) {
	// Создаем тестового пользователя
	user := models.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// Вызываем метод BeforeSave вручную (в тестах не используется GORM)
	err := user.BeforeSave(nil)
	assert.NoError(t, err)

	// Проверяем, что пароль был хеширован
	assert.NotEqual(t, "password123", user.Password)

	// Проверяем, что хешированный пароль можно проверить
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	assert.NoError(t, err)

	// Проверяем метод ValidatePassword
	err = user.ValidatePassword("password123")
	assert.NoError(t, err)

	// Проверяем метод ValidatePassword с неверным паролем
	err = user.ValidatePassword("wrongpassword")
	assert.Error(t, err)
} 