package command

import (
	"fmt"
	"quit4real.today/logger"
	"quit4real.today/src/handler/service"
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type SubscriptionCommandHandler struct {
	SteamApi               *service.SteamService
	SubscriptionRepository *repository.SubscriptionRepository
	FailsCommandHandler    *FailsCommandHandler
	GameCommandHandler     *GameCommandHandler
}

// Add adds a new subscription for a user and retrieves the played time for the game.
func (handler *SubscriptionCommandHandler) Add(subscription model.Subscription) error {
	if subscription.PlatformId == "steam" {
		game, err := handler.SteamApi.GetRequestedGame(
			subscription.PlatFormUserId, subscription.PlatformGameId)
		if err != nil {
			return err
		}
		// Adds a new game to the game table only if it does not exist already.
		err = handler.GameCommandHandler.Add(subscription.PlatformGameId, game.Name, "steam")
		if err != nil {
			return fmt.Errorf("failed to add Game to game table: %v", err)
		}
		err = handler.SubscriptionRepository.Add(
			subscription.DisplayName,
			subscription.PlatformId,
			subscription.PlatformGameId,
			subscription.PlatFormUserId,
			game.PlaytimeForever)
		if err != nil {
			return fmt.Errorf("subscription most likely already exists: %v", err)
		}
		return nil
	}
	return fmt.Errorf("we currently only support steam as a valid platform")
}

// UpdateFromSteamApi updates the user's subscriptions based on recent games fetched from the Steam API.
func (handler *SubscriptionCommandHandler) UpdateFromSteamApi(steamId string) {
	apiResponse, err := handler.SteamApi.FetchRecentGames(steamId)
	if err != nil {
		logger.Fail("failed to fetch player information for player: " + steamId + " | ERROR: " + err.Error())
		return
	}

	trackedGamesByUser, err := handler.SubscriptionRepository.GetAllForUser(steamId)
	if err != nil {
		logger.Fail("failed to get all tracked games for player: " + steamId + " | ERROR: " + err.Error())
		return
	}

	failedGames := handler.SteamApi.GetOnlyFailed(apiResponse, trackedGamesByUser)

	for _, failInfo := range failedGames {

		// Update subscription repository
		if err := handler.SubscriptionRepository.Update(steamId, failInfo.DbTrack.PlatformGameId, failInfo.SteamApiGame.PlaytimeForever); err != nil {
			logger.Fail("Error updating subscription for user: " + steamId + " | ERROR: " + err.Error())
			return
		}
		logger.Info("A fail from User: " + steamId + " playing game " + failInfo.DbTrack.PlatformGameId)

		// Add the failure record using the FailsCommandHandler
		err = handler.FailsCommandHandler.Add(failInfo.DbTrack, failInfo.SteamApiGame.PlaytimeForever)
		if err != nil {
			logger.Fail("Error creating a Fail: " + steamId + " | ERROR: " + err.Error())
		}
	}
}
