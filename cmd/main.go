package main

import (
	"log"

	"music-library/internal/api"
	"music-library/internal/database"
)

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
