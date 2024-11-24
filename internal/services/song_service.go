package services

import (
	"fmt"
	"net/http"
	"os"

	"music-library/internal/logger"
	"music-library/internal/models"

	"encoding/json"

	"github.com/sirupsen/logrus"
)

func FetchExternalSong(group, song string) (*models.Song, error) {
	apiURL := os.Getenv("API_URL")
	reqURL := fmt.Sprintf("%s/info?group=%s&song=%s", apiURL, group, song)

	logger.Log.WithFields(logrus.Fields{
		"group": group,
		"song":  song,
		"url":   reqURL,
	}).Debug("Sending request to external API")

	resp, err := http.Get(reqURL)
	if err != nil {
		logger.Log.WithError(err).Error("Failed to send request to external API")
		return nil, err
	}
	defer resp.Body.Close()

	// Логируем статус ответа
	logger.Log.WithFields(logrus.Fields{
		"status": resp.StatusCode,
		"url":    reqURL,
	}).Info("Received response from external API")

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("external API error: %v", resp.Status)
		logger.Log.WithError(err).Error("Invalid response from external API")
		return nil, err
	}

	var songDetails models.Song
	if err := json.NewDecoder(resp.Body).Decode(&songDetails); err != nil {
		logger.Log.WithError(err).Error("Failed to decode response from external API")
		return nil, err
	}

	logger.Log.WithFields(logrus.Fields{
		"song_details": songDetails,
	}).Debug("Successfully fetched song details from external API")

	return &songDetails, nil
}
