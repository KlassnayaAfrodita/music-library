package main

import (
	"fmt"
	_ "music-library/docs" // Подключаем автоматически сгенерированные Swagger-документы

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"music-library/internal/api"
	"music-library/internal/logger"
	"music-library/internal/middleware"
)

func main() {
	// Инициализация логгера
	logger.Init()

	// Подключение базы данных
	// database.ConnectDB()
	// defer database.DB.Close()

	router := api.SetupRouter()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Use(middleware.RequestLogger())

	if err := router.Run(":8080"); err != nil {
		fmt.Errorf("start router error")
	}
}
