package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"music-library/internal/database"
	"music-library/internal/logger"
	"music-library/internal/models"
	"music-library/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// GetSongs godoc
// @Summary      Получение песен с фильтрацией и пагинацией
// @Description  Возвращает список песен с фильтрацией по названию группы и песни, а также поддерживает пагинацию
// @Tags         Songs
// @Param        group   query   string  false  "Название группы"
// @Param        song    query   string  false  "Название песни"
// @Param        page    query   int     false  "Номер страницы" default(1)
// @Param        limit   query   int     false  "Количество элементов на странице" default(10)
// @Success      200     {object}  []models.Song
// @Failure      400     {object}  map[string]string
// @Failure      500     {object}  map[string]string
// @Router       /songs [get]
func GetSongs(c *gin.Context) {
	logger.Log.Debug("Entering GetSongs handler")

	var songs []models.Song
	group := c.Query("group")
	song := c.Query("song")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		logger.Log.WithFields(logrus.Fields{"page": pageStr}).Debug("Invalid page parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		logger.Log.WithFields(logrus.Fields{"limit": limitStr}).Debug("Invalid limit parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	offset := (page - 1) * limit
	query := "SELECT * FROM songs WHERE 1=1"
	if group != "" {
		query += " AND group_name ILIKE '%' || $1 || '%'"
	}
	if song != "" {
		query += " AND song_name ILIKE '%' || $2 || '%'"
	}
	query += " LIMIT $3 OFFSET $4"

	logger.Log.WithFields(logrus.Fields{"group": group, "song": song, "page": page, "limit": limit}).Info("Fetching songs with filters")

	err = database.DB.Select(&songs, query, group, song, limit, offset)
	if err != nil {
		logger.Log.WithError(err).Debug("Error fetching songs from the database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching songs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"songs": songs, "page": page, "limit": limit})
	logger.Log.Info("Songs fetched successfully")
}

// GetLyrics godoc
// @Summary      Получение текста песни с пагинацией
// @Description  Возвращает текст песни, разделённый на куплеты, с поддержкой пагинации
// @Tags         Songs
// @Param        id     path     int  true   "ID песни"
// @Param        page   query    int  false  "Номер страницы" default(1)
// @Param        limit  query    int  false  "Количество строк на странице" default(2)
// @Success      200    {object} map[string]interface{}
// @Failure      400    {object}  map[string]string
// @Failure      404    {object}  map[string]string
// @Router       /songs/{id}/lyrics [get]
func GetLyrics(c *gin.Context) {
	logger.Log.Debug("Entering GetLyrics handler")

	idStr := c.Param("id")
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "2")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"song_id": idStr}).Debug("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		logger.Log.WithFields(logrus.Fields{"page": pageStr}).Debug("Invalid page parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
		return
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		logger.Log.WithFields(logrus.Fields{"limit": limitStr}).Debug("Invalid limit parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}

	var song models.Song
	query := "SELECT * FROM songs WHERE id = $1"
	err = database.DB.Get(&song, query, id)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"song_id": id}).Debug("Song not found in database")
		c.JSON(http.StatusNotFound, gin.H{"error": "Song not found"})
		return
	}

	verses := strings.Split(song.Lyrics, "\n")
	totalVerses := len(verses)
	totalPages := int(math.Ceil(float64(totalVerses) / float64(limit)))

	start := (page - 1) * limit
	end := start + limit
	if end > totalVerses {
		end = totalVerses
	}

	logger.Log.WithFields(logrus.Fields{"song_id": id, "page": page}).Info("Fetching song lyrics with pagination")
	c.JSON(http.StatusOK, gin.H{"lyrics": verses[start:end], "page": page, "total_pages": totalPages})
	logger.Log.Info("Lyrics fetched successfully")
}

// UpdateSong godoc
// @Summary      Изменение данных песни
// @Description  Обновляет информацию о песне
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        id    path      int         true  "ID песни"
// @Param        song  body      models.Song true  "Новые данные песни"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /songs/{id} [put]
func UpdateSong(c *gin.Context) {
	logger.Log.Debug("Entering UpdateSong handler")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"song_id": idStr}).Debug("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		logger.Log.WithError(err).Debug("Invalid input for updating song")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	query := `UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, lyrics = $4, link = $5 WHERE id = $6`
	_, err = database.DB.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Lyrics, song.Link, id)
	if err != nil {
		logger.Log.WithError(err).Debug("Failed to update song in database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	logger.Log.WithFields(logrus.Fields{"song_id": id}).Info("Song updated successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// DeleteSong godoc
// @Summary      Удаление песни
// @Description  Удаляет песню из библиотеки
// @Tags         Songs
// @Param        id   path      int  true  "ID песни"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /songs/{id} [delete]
func DeleteSong(c *gin.Context) {
	logger.Log.Debug("Entering DeleteSong handler")

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Log.WithFields(logrus.Fields{"song_id": idStr}).Debug("Invalid song ID")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	query := "DELETE FROM songs WHERE id = $1"
	_, err = database.DB.Exec(query, id)
	if err != nil {
		logger.Log.WithError(err).Debug("Failed to delete song from database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	logger.Log.WithFields(logrus.Fields{"song_id": id}).Info("Song deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}

// AddSong godoc
// @Summary      Добавление новой песни
// @Description  Добавляет новую песню в библиотеку, запрашивая данные из внешнего API
// @Tags         Songs
// @Accept       json
// @Produce      json
// @Param        song  body      models.Song  true  "Данные песни"
// @Success      201   {object}  models.Song
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /songs [post]
func AddSong(c *gin.Context) {
	var songInput models.Song
	if err := c.ShouldBindJSON(&songInput); err != nil {
		logger.Log.WithError(err).Debug("Invalid input data for adding a new song")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	logger.Log.WithFields(logrus.Fields{
		"group_name": songInput.GroupName,
		"song_name":  songInput.SongName,
	}).Info("Adding a new song")

	apiData, err := services.FetchExternalSong(songInput.GroupName, songInput.SongName)
	if err != nil {
		logger.Log.WithError(err).Debug("Failed to fetch song details from external API")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch song details: " + err.Error()})
		return
	}

	songInput.ReleaseDate = apiData.ReleaseDate
	songInput.Lyrics = apiData.Lyrics
	songInput.Link = apiData.Link

	query := `INSERT INTO songs (group_name, song_name, release_date, lyrics, link) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err = database.DB.Exec(query, songInput.GroupName, songInput.SongName, songInput.ReleaseDate, songInput.Lyrics, songInput.Link)
	if err != nil {
		logger.Log.WithError(err).Debug("Failed to insert song into database")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save song"})
		return
	}

	c.JSON(http.StatusCreated, songInput)
	logger.Log.Info("Song added successfully")
}
