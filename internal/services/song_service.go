package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"music-library/internal/models"
)

func FetchExternalSong(group, song string) (*models.Song, error) {
	apiURL := os.Getenv("API_URL")
	reqURL := fmt.Sprintf("%s/info?group=%s&song=%s", apiURL, group, song)

	resp, err := http.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external API error: %v", resp.Status)
	}

	var songDetails models.Song
	if err := json.NewDecoder(resp.Body).Decode(&songDetails); err != nil {
		return nil, err
	}
	return &songDetails, nil
}
