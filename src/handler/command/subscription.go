package command

import (
	"quit4real.today/logger"
	"quit4real.today/src/api"
	"quit4real.today/src/repository"
)

type SubscriptionCommandHandler struct {
	SteamApi               *api.SteamApi
	SubscriptionRepository *repository.SubscriptionRepository
	FailsCommandHandler    *FailsCommandHandler
}

// Add adds a new subscription for a user and retrieves the played time for the game.
//func (handler *SubscriptionCommandHandler) Add(id string, gameId string) error {
//	playedAmount, err := handler.SteamApi.GetRequestedGamePlayedTime(id, gameId)
//	if err != nil {
//		return err
//	}
//	err = handler.SubscriptionRepository.Add(id, gameId, playedAmount)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// UpdateFromSteamApi updates the user's subscriptions based on recent games fetched from the Steam API.
func (handler *SubscriptionCommandHandler) UpdateFromSteamApi(steamId string) {
	apiResponse, err := handler.SteamApi.FetchRecentGames(steamId)
	if err != nil {
		logger.Fail("failed to fetch player information for player: " + steamId + " | ERROR: " + err.Error())
		return
	}

	trackedGamesByUser, err := handler.SubscriptionRepository.GetAll(steamId)
	if err != nil {
		logger.Fail("failed to get all tracked games for player: " + steamId + " | ERROR: " + err.Error())
		return
	}

	failedGames := handler.SteamApi.GetOnlyFailed(apiResponse, trackedGamesByUser)

	for _, failInfo := range failedGames {

		// Update subscription repository
		if err := handler.SubscriptionRepository.Update(steamId, failInfo.DbTrack.GameId, failInfo.SteamApiGame.PlaytimeForever); err != nil {
			logger.Fail("Error updating subscription for user: " + steamId + " | ERROR: " + err.Error())
			return
		}
		logger.Info("A fail from User: " + steamId + " playing game " + failInfo.DbTrack.GameId)

		// Add the failure record using the FailsCommandHandler
		err = handler.FailsCommandHandler.Add(failInfo.DbTrack, failInfo.SteamApiGame.PlaytimeForever)
		if err != nil {
			logger.Fail("Error creating a Fail: " + steamId + " | ERROR: " + err.Error())
		}
	}
}
