package command

import (
	"fmt"
	"project/main/model"
	"project/main/repository"
	"strconv"
)

type TrackerCommandHandler struct {
	trackerRepository *repository.TrackerRepository
}

func (handler *TrackerCommandHandler) Add(platformId string, gameId string) error {
	var err = handler.trackerRepository.Add(platformId, gameId)
	if err != nil {
		return err
	}
	return nil
}

func (handler *TrackerCommandHandler) Validate(steamId string, apiResponse *mode.SteamApiResponse) error {
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
