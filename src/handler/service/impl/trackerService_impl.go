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

func (service *TrackerServiceImpl) UpdateSteamTracker(steamId string, steamApiResponse *model.SteamApiResponse) {

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
