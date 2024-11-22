package api

import (
	"music-library/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/songs", handlers.GetSongs)
	r.GET("/songs/:id/lyrics", handlers.GetLyrics)
	r.POST("/songs", handlers.AddSong)
	r.PUT("/songs/:id", handlers.UpdateSong)
	r.DELETE("/songs/:id", handlers.DeleteSong)

	return r
}
