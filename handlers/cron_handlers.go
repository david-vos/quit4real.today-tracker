package handlers

import (
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

type CronController struct {
	UserRepoContr    *repository.UserRepoController
	TrackerRepoContr *repository.TrackerRepoController
	FailsContr       *FailsController
}

func closeBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		config.HandleError("Error closing response body: %v\n", err)
	}
}

func (c *CronController) fetchSteamApiData(steamId string) (*models.ApiResponse, error) {
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

func (c *CronController) updateTrackerForGames(steamId string, apiResponse *models.ApiResponse) {
	trackedGamesByUser, err := c.TrackerRepoContr.GetUserTracker(steamId)
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
		// Check if the played amount is greater than the current playtime
		if game.PlaytimeForever > trackedGame.PlayedAmount {
			// Update the tracker
			if err := c.TrackerRepoContr.UpdateTracker(steamId, gameID, game.PlaytimeForever); err != nil {
				config.HandleError("Error updating tracker", err)
				return // Early return on error
			}

			fmt.Println("A fail from User: " + steamId + " playing game " + gameID)
			err := c.FailsContr.createFail(trackedGame)
			if err != nil {
				config.HandleError("Error creating a Fail", err)
			}
			//TODO: This should send out a message to another service that lisons to these events kafka 0_o
			return
		}
	}
}

func (c *CronController) GetSteamApiData(steamId string) {
	apiResponse, err := c.fetchSteamApiData(steamId)
	if err != nil {
		config.HandleError("Failed to fetch Steam API data", err)
		return
	}

	c.updateTrackerForGames(steamId, apiResponse)
}

func (c *CronController) updateAndSendNotify() {
	users, err := c.UserRepoContr.GetAllUsers()
	if err != nil {
		config.HandleError("Error getting all users", err)
		return
	}
	for _, user := range users {
		c.GetSteamApiData(user.SteamId)
	}
}

func (c *CronController) SetupCronJobs() {
	c.updateAndSendNotify()

	cronJob := cron.New()
	_, err := cronJob.AddFunc("@every 2m", func() {
		c.updateAndSendNotify()
	})
	if err != nil {
		config.HandleError("Error adding cron job", err)
		return
	}

	cronJob.Start()
}
