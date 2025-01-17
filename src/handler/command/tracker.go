package command

import (
	"quit4real.today/logger"
	"quit4real.today/src/api"
	"quit4real.today/src/repository"
)

type TrackerCommandHandler struct {
	SteamApi            *api.SteamApi
	TrackerRepository   *repository.TrackerRepository
	FailsCommandHandler *FailsCommandHandler
}

func (handler *TrackerCommandHandler) Add(platformId string, gameId string) error {
	var err = handler.TrackerRepository.Add(platformId, gameId)
	if err != nil {
		return err
	}
	return nil
}

func (handler *TrackerCommandHandler) UpdateFromSteamApi(steamId string) {
	apiResponse, err := handler.SteamApi.FetchApiData(steamId)
	if err != nil {
		logger.Fail("failed get fetch player information for player: " + steamId + " | ERROR: " + err.Error())
		return
	}

	trackedGamesByUser, err := handler.TrackerRepository.GetAll(steamId)
	if err != nil {
		logger.Fail("failed get all tracked games for player: " + steamId)
		return
	}

	failedGames := handler.SteamApi.GetOnlyFailed(steamId, apiResponse, trackedGamesByUser)

	for _, failInfo := range failedGames {
		// Update tracker repository
		if err := handler.TrackerRepository.Update(steamId, failInfo.DbTrack.GameId, failInfo.SteamApiGame.PlaytimeForever); err != nil {
			logger.Fail("Error updating tracker: " + err.Error())
			return
		}
		logger.Info("A fail from User: " + steamId + " playing game " + failInfo.DbTrack.GameId)

		err := handler.FailsCommandHandler.Add(failInfo.DbTrack, failInfo.SteamApiGame.PlaytimeForever)
		if err != nil {
			logger.Fail("Error creating a Fail: " + err.Error())
		}
	}
}
