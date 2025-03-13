package database

import (
	"fmt"
	"log"
	"os"

	"github.com/omega/notes-app/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB инициализирует соединение с базой данных
func InitDB() (*gorm.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	config := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, err
	}

	// Миграция моделей
	err = DB.AutoMigrate(&models.User{}, &models.Note{})
	if err != nil {
		return nil, err
	}

	log.Println("База данных успешно подключена и мигрирована")
	return DB, nil
}

// GetDB возвращает экземпляр соединения с базой данных
func GetDB() *gorm.DB {
	return DB
} 