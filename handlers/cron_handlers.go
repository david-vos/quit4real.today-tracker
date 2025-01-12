package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"net/http"
	"project/config"
	"project/models"
	"project/repository"
	"strconv"
)

func closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		config.HandleError("Error closing response body: %v\n", err)
	}
}

func fetchSteamApiData(steamId string) (*models.ApiResponse, error) {
	apiKey := config.GetSteamApiKey()
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetRecentlyPlayedGames/v0001/?key=%s&steamid=%s&format=json", apiKey, steamId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var apiResponse models.ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return &apiResponse, nil
}

func updateTrackerForGames(db *sql.DB, steamId string,
	apiResponse *models.ApiResponse) {
	trackedGamesByUser, err := repository.GetUserTracker(db, steamId)
	if err != nil {
		config.HandleError("Error getting user tracker", err)
		return // Early return on error
	}

	// Create a map for quick lookup of tracked game IDs
	trackedGameMap := make(map[string]models.Tracker)
	for _, trackedGame := range trackedGamesByUser {
		trackedGameMap[trackedGame.GameId] = trackedGame
	}

	for _, game := range apiResponse.Response.Games {
		gameID := strconv.Itoa(game.AppID)
		trackedGame, exists := trackedGameMap[gameID]
		if !exists {
			continue // Skip if the game is not tracked
		}

		// Update the tracker
		if err := repository.UpdateTracker(db, steamId, gameID,
			game.PlaytimeForever); err != nil {
			config.HandleError("Error updating tracker", err)
			return // Early return on error
		}

		// Check if the played amount is greater than the current playtime
		if trackedGame.PlayedAmount > game.PlaytimeForever {
			fmt.Println("YOU PLAYED THE GAME NO YOU DIE")
			return // Early return after printing the message
		}
	}
}

func GetSteamApiData(db *sql.DB, steamId string) {
	apiResponse, err := fetchSteamApiData(steamId)
	if err != nil {
		config.HandleError("Failed to fetch Steam API data", err)
		return
	}

	updateTrackerForGames(db, steamId, apiResponse)
}

func updateAndSendNotify(db *sql.DB) {
	users, err := repository.GetAllUsers(db)
	if err != nil {
		config.HandleError("Error getting all users", err)
		return
	}
	for _, user := range users {
		GetSteamApiData(db, user.SteamId)
	}
}

func SetupCronJobs(db *sql.DB) {
	c := cron.New()

	_, err := c.AddFunc("@hourly", func() {
		updateAndSendNotify(db)
	})
	if err != nil {
		config.HandleError("Error adding cron job", err)
		return
	}

	c.Start()
	defer c.Stop()
}
