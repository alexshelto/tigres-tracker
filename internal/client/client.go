package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/alexshelto/tigres-tracker/config"
	"github.com/alexshelto/tigres-tracker/dto"
)

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(cfg config.ClientConfig) *APIClient {
	return &APIClient{
		BaseURL:    cfg.BaseURL,
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *APIClient) PostSongPlay(userID, songName, guildID string) ([]byte, error) {
	// Create the payload data
	playData := map[string]interface{}{
		"user_id":   userID,
		"song_name": songName,
		"guild_id":  guildID,
	}

	// Convert payload to JSON
	data, err := json.Marshal(playData)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %w", err)
	}

	// Format the URL for the POST request (for example, "/song/play")
	url := fmt.Sprintf("%s/song", c.BaseURL)

	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Check for non-200 status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

func (c *APIClient) GetTopSongsInGuild(guildID string, limit int) ([]dto.SongRequestCountDTO, error) {
	url := fmt.Sprintf("%s/song/top?guild_id=%s&limit=%d", c.BaseURL, guildID, limit)

	// Make the GET request
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return nil, errors.New("unexpected error")
	}
	// Parse response body

	var songs []dto.SongRequestCountDTO
	err = parseResponseBody(resp.Body, &songs)
	if err != nil {
		fmt.Println("Error parsing response body:", err)
		return nil, errors.New("could not parse body")
	}

	return songs, nil
}

func (c *APIClient) GetTopSongsByUserInGuild(userID string, guildID string, limit int) ([]dto.SongRequestCountDTO, error) {
	url := fmt.Sprintf("%s/user/%s/song/top?guild_id=%s&limit=%d", c.BaseURL, userID, guildID, limit)

	// Make the GET request
	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return nil, errors.New("unexpected error")
	}
	// Parse response body

	var songs []dto.SongRequestCountDTO
	err = parseResponseBody(resp.Body, &songs)
	if err != nil {
		fmt.Println("Error parsing response body:", err)
		return nil, errors.New("could not parse body")
	}

	return songs, nil
}

func (c *APIClient) GetTotalSongPlaysInGuild(guildID string) (dto.TotalSongPlayDTO, error) {
	url := fmt.Sprintf("%s/song/count?guild_id=%s", c.BaseURL, guildID)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return dto.TotalSongPlayDTO{}, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return dto.TotalSongPlayDTO{}, errors.New("unexpected error")
	}

	var total dto.TotalSongPlayDTO
	err = parseResponseBody(resp.Body, &total)

	if err != nil {
		fmt.Println("Error parsing response body:", err)
		return total, errors.New("could not parse body")
	}

	return total, nil
}

func (c *APIClient) GetTotalUserSongPlaysInGuild(userID string, guildID string) (dto.TotalSongPlayDTO, error) {
	url := fmt.Sprintf("%s/user/%s/song/count?guild_id=%s", c.BaseURL, userID, guildID)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return dto.TotalSongPlayDTO{}, fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	// Check for non-2xx status codes
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return dto.TotalSongPlayDTO{}, errors.New("unexpected error")
	}

	var total dto.TotalSongPlayDTO
	err = parseResponseBody(resp.Body, &total)

	if err != nil {
		fmt.Println("Error parsing response body:", err)
		return total, errors.New("could not parse body")
	}

	return total, nil
}

func parseResponseBody(body io.Reader, target interface{}) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(target)
}
