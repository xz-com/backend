package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/omega/notes-app/internal/database"
	"github.com/omega/notes-app/internal/models"
)

// NoteRequest представляет данные для создания или обновления заметки
type NoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content"`
}

// CreateNote обрабатывает запрос на создание новой заметки
func CreateNote(c *gin.Context) {
	var req NoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	// Создаем новую заметку
	note := models.Note{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID.(uint),
	}

	// Сохраняем заметку в базе данных
	if err := database.GetDB().Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при создании заметки"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "заметка успешно создана",
		"note":    note,
	})
}

// GetNotes возвращает все заметки пользователя
func GetNotes(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	var notes []models.Note
	if err := database.GetDB().Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при получении заметок"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notes": notes,
	})
}

// GetNote возвращает заметку по ID
func GetNote(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	// Получаем ID заметки из URL
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID заметки"})
		return
	}

	var note models.Note
	result := database.GetDB().Where("id = ? AND user_id = ?", noteID, userID).First(&note)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "заметка не найдена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"note": note,
	})
}

// UpdateNote обновляет заметку по ID
func UpdateNote(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	// Получаем ID заметки из URL
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID заметки"})
		return
	}

	// Проверяем, существует ли заметка и принадлежит ли она пользователю
	var note models.Note
	result := database.GetDB().Where("id = ? AND user_id = ?", noteID, userID).First(&note)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "заметка не найдена"})
		return
	}

	// Получаем данные для обновления
	var req NoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Обновляем заметку
	note.Title = req.Title
	note.Content = req.Content

	if err := database.GetDB().Save(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при обновлении заметки"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "заметка успешно обновлена",
		"note":    note,
	})
}

// DeleteNote удаляет заметку по ID
func DeleteNote(c *gin.Context) {
	// Получаем ID пользователя из контекста
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "пользователь не аутентифицирован"})
		return
	}

	// Получаем ID заметки из URL
	noteID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный ID заметки"})
		return
	}

	// Проверяем, существует ли заметка и принадлежит ли она пользователю
	var note models.Note
	result := database.GetDB().Where("id = ? AND user_id = ?", noteID, userID).First(&note)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "заметка не найдена"})
		return
	}

	// Удаляем заметку
	if err := database.GetDB().Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка при удалении заметки"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "заметка успешно удалена",
	})
} 