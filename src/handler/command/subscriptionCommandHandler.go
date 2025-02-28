package command

import (
	"fmt"
	"quit4real.today/logger"
	"quit4real.today/src/handler/service"
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type SubscriptionCommandHandlerImpl struct {
	SubscriptionRepository repository.SubscriptionRepository
	// The service layer should take care of all of this...
	// The service layer is only allowed to talk to other service layers and downwards to a command/query handler.
	// the commmand/query handlers are only allowed to talk to their repository/projection
	SteamService        service.SteamService
	FailsCommandHandler *FailsCommandHandlerImpl
	GameCommandHandler  *GameCommandHandlerImpl
}

// Add adds a new subscription for a user and retrieves the played time for the game.
func (handler *SubscriptionCommandHandlerImpl) Add(subscription model.Subscription) error {
	if subscription.PlatformId == "steam" {
		game, err := handler.SteamService.GetRequestedGame(
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

func (handler *SubscriptionCommandHandlerImpl) UpdateSubscriptions(steamId string, failedGames []model.MatchedDbGameToSteamGameInfo) {
	for _, failInfo := range failedGames {
		// Update subscription repository
		if err := handler.SubscriptionRepository.Update(steamId, failInfo.DbTrack.PlatformGameId, failInfo.SteamApiGame.PlaytimeForever); err != nil {
			logger.Fail("Error updating subscription for user: " + steamId + " | ERROR: " + err.Error())
			return
		}
		logger.Info("A fail from User: " + steamId + " playing game " + failInfo.DbTrack.PlatformGameId)

		// Add the failure record using the FailsCommandHandler
		err := handler.FailsCommandHandler.Add(failInfo.DbTrack, failInfo.SteamApiGame.PlaytimeForever)
		if err != nil {
			logger.Fail("Error creating a Fail: " + steamId + " | ERROR: " + err.Error())
		}
	}
}
