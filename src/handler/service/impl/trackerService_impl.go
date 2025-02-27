package impl

import (
	"fmt"
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
	"strconv"
	"time"
)

type TrackerServiceImpl struct {
	TrackerRepository repository.TrackerRepository
}

func (service *TrackerServiceImpl) UpdateSteamTrackers(steamId string, steamApiResponse *model.SteamApiResponse) []error {

	var errors []error
	for _, game := range steamApiResponse.Response.Games {
		// Fetch the existing tracker for this game
		tracker, err := service.TrackerRepository.GetLatestTrackerByUserIdAndGameId(steamId, game.AppID)
		if err != nil {
			// Handle the case where tracker might not exist
			errors = append(errors, err)
			continue
		}

		// Calculate the difference in time played
		timeDifference := game.PlaytimeForever - tracker.NewTotalTimePlayed

		// Update tracker fields
		tracker.TimePlayed += timeDifference
		tracker.NewTotalTimePlayed = game.PlaytimeForever
		tracker.AmountOfLogins++

		// Check if the day has changed
		if !tracker.Day.Equal(time.Now().Truncate(24 * time.Hour)) {
			tracker.Day = time.Now().Truncate(24 * time.Hour)
			tracker.AmountOfLogins = 1
			err = service.TrackerRepository.UpdateTracker(tracker)
			if err != nil {
				errors = append(errors, err)
				continue
			}
		}

		// Save the updated tracker
		err = service.TrackerRepository.UpdateTracker(tracker)
		if err != nil {
			errors = append(errors, err)
			continue
		}
	}
	return errors
}

func (service *TrackerServiceImpl) CreateSteamTrackers(steamId string, allGames *model.SteamAPIAllResponse) error {
	// doing it this way might break stuff if a user buys a new game after signing up to Quit4Real.
	// Many ways around this problem though. We can always add a button to the FE asking for a re-index
	userIdInt, err := strconv.Atoi(steamId)
	if err != nil {
		return fmt.Errorf("cannot decode steamId to type Int when trying to create trackers. Error: %s", err.Error())
	}

	for _, game := range allGames.Games {
		tracker := model.Tracker{
			UserID:             userIdInt,
			GameID:             game.Appid,
			PlatformID:         "steam",
			TimePlayed:         0,
			NewTotalTimePlayed: game.PlaytimeForever,
			AmountOfLogins:     0,
			Day:                time.Now().Truncate(24 * time.Hour),
		}
		err := service.TrackerRepository.InsertTracker(tracker)
		if err != nil {
			return fmt.Errorf("cannot insert tracker into database: %s", err.Error())
		}

	}
	return nil
}
