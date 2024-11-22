package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"music-library/internal/database"
	"music-library/internal/models"

	"github.com/gin-gonic/gin"
)

// GetSongs retrieves a list of songs with filtering and pagination
func GetSongs(c *gin.Context) {
	var songs []models.Song

	// Получение параметров фильтрации
	group := c.Query("group")
	song := c.Query("song")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset := (page - 1) * limit

	// Формирование запроса с фильтрацией
	query := "SELECT * FROM songs WHERE 1=1"
	if group != "" {
		query += " AND group_name ILIKE '%' || $1 || '%'"
	}
	if song != "" {
		query += " AND song_name ILIKE '%' || $2 || '%'"
	}
	query += " LIMIT $3 OFFSET $4"

	err = database.DB.Select(&songs, query, group, song, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching songs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs, "page": page, "limit": limit})
}

// GetLyrics retrieves song lyrics with pagination by verses
func GetLyrics(c *gin.Context) {
	idStr := c.Param("id")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "2") // Пагинация по куплетам (например, 2 строки)

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	var song models.Song
	query := "SELECT * FROM songs WHERE id = $1"
	err = database.DB.Get(&song, query, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	verses := strings.Split(song.Lyrics, "\n") // Разбиваем текст на куплеты
	totalVerses := len(verses)
	totalPages := int(math.Ceil(float64(totalVerses) / float64(limit)))

	if page > totalPages {
		c.JSON(http.StatusOK, gin.H{"lyrics": []string{}, "page": page, "total_pages": totalPages})
		return
	}

	start := (page - 1) * limit
	end := start + limit
	if end > totalVerses {
		end = totalVerses
	}

	c.JSON(http.StatusOK, gin.H{
		"lyrics":      verses[start:end],
		"page":        page,
		"total_pages": totalPages,
	})
}

// AddSong adds a new song
func AddSong(c *gin.Context) {
	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Вставка в базу данных
	query := `INSERT INTO songs (group_name, song_name, release_date, lyrics, link) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := database.DB.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Lyrics, song.Link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song"})
		return
	}
	c.JSON(http.StatusCreated, song)
}

// UpdateSong updates the details of a song
func UpdateSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := `UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, lyrics = $4, link = $5 WHERE id = $6`
	_, err = database.DB.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Lyrics, song.Link, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// DeleteSong removes a song from the library
func DeleteSong(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	query := "DELETE FROM songs WHERE id = $1"
	_, err = database.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
