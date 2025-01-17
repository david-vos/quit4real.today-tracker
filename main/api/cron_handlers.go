// TODO: I have to go to sleep so I will handle the cron job part of this later
// These should not be under API, instead I think these should be under event. Events being extranal to and service
// outside of this app
package api

import (
	"encoding/json"
	"fmt"
	"github.com/robfig/cron/v3"
	"io"
	"net/http"
	"project/config"
	"project/main/model"
	"project/repository"
	"strconv"
)

type CronController struct {
	UserRepoContr    *repository.UserRepoController
	TrackerRepoContr *repository.TrackerRepoController
	FailsContr       *FailsController
}

func (c *CronController) updateTrackerForGames(steamId string, apiResponse *model.ApiResponse) {
	trackedGamesByUser, err := c.TrackerRepoContr.GetUserTracker(steamId)
	if err != nil {
		config.HandleError("Error getting user tracker", err)
		return // Early return on error
	}

	// Create a map for quick lookup of tracked game IDs
	trackedGameMap := make(map[string]model.Tracker)
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
		// This will always fail on the first iteration when a new user is created it does not actually set the
		// PlayedAmount in the beginning making it playtimeForever > 0 which is always true
		if game.PlaytimeForever > trackedGame.PlayedAmount {
			// Update the tracker
			if err := c.TrackerRepoContr.UpdateTracker(steamId, gameID, game.PlaytimeForever); err != nil {
				config.HandleError("Error updating tracker", err)
				return // Early return on error
			}

			fmt.Println("A fail from User: " + steamId + " playing game " + gameID)
			err := c.FailsContr.createFail(trackedGame, game.PlaytimeForever)
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
	cronJob := cron.New()
	// ((24*60)/10)*694 ~= 100.000 the STEAM API limit
	// 694 -> max amount of users :thinking per API key
	_, err := cronJob.AddFunc("@every 10m", func() {
		c.updateAndSendNotify()
	})
	if err != nil {
		config.HandleError("Error adding cron job", err)
		return
	}

	cronJob.Start()
}
