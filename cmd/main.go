package main

import (
	"log"
	_ "music-library/docs" // Подключаем автоматически сгенерированные Swagger-документы

	"music-library/internal/api"
	"music-library/internal/database"
)

// @title Example API
// @version 1.0
// @description This is a sample server for Swagger.
// @host localhost:8080
// @BasePath /
func main() {
	// Подключение базы данных
	database.ConnectDB()
	defer database.DB.Close()

	// Настройка роутера
	router := api.SetupRouter()

	// Запуск сервера
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
