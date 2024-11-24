package main

import (
	"log"
	_ "music-library/docs" // Подключаем автоматически сгенерированные Swagger-документы

	"music-library/internal/api"
	"music-library/internal/database"
	"music-library/internal/logger"
	"music-library/internal/middleware"
)

func main() {
	// Инициализация логгера
	logger.Init()

	// Подключение базы данных
	database.ConnectDB()
	defer database.DB.Close()

	router := api.SetupRouter()
	router.Use(middleware.RequestLogger())

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
