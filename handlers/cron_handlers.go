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

func GetSteamApiData(db *sql.DB, steamId string) {
	apiKey := config.GetSteamApiKey()
	url := fmt.Sprintf("http://api.steampowered.com/IPlayerService/GetRecentlyPlayedGames/v0001/?key=%s&steamid=%s&format=json", apiKey, steamId)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making HTTP request: %v\n", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing response body: %v\n", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("HTTP request failed with status: %s\n", resp.Status)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	var apiResponse models.ApiResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	trackedGamesByUser, err := repository.GetUserTracker(db, steamId)
	if err != nil {
		fmt.Printf("Error getting user tracker: %v\n", err)
	}

	// main logic
	for _, game := range apiResponse.Response.Games {
		for _, value := range trackedGamesByUser {
			if value.GameId == strconv.Itoa(game.AppID) {
				// These are now only games that are followed by the User
				err := repository.UpdateTracker(db, steamId, strconv.Itoa(game.AppID), game.PlaytimeForever)
				if err != nil {
					return
				}
				if value.PlayedAmount > game.PlaytimeForever {
					// THis should actually send a post or somehting to someone that will update anything idk yet.
					fmt.Println("YOU PLAYED THE GAME NO YOU DIE")
				}
			}
		}
	}
}

func updateAndSendNotify(db *sql.DB) {
	users := repository.GetAllUsers(db)
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
		fmt.Printf("Error adding cron job: %v\n", err)
		return
	}

	c.Start()
	defer c.Stop()
}
