package impl

import (
	"quit4real.today/logger"
	"quit4real.today/src/handler/query"
	"quit4real.today/src/handler/service"
	"quit4real.today/src/model"
)

type SubscriptionServiceImpl struct {
	SubscriptionQueryHandler query.SubscriptionQueryHandlerImpl
	SteamService             service.SteamService
}

func NewSubscriptionServiceImpl(subscriptionQueryHandler query.SubscriptionQueryHandlerImpl, steamService service.SteamService) *SubscriptionServiceImpl {
	return &SubscriptionServiceImpl{
		SubscriptionQueryHandler: subscriptionQueryHandler,
		SteamService:             steamService,
	}
}

// GetOnlyFailedSteam gets all the failed games for a steam user
func (service *SubscriptionServiceImpl) GetOnlyFailedSteam(steamId string, apiResponse *model.SteamApiResponse) []model.MatchedDbGameToSteamGameInfo {
	trackedGamesByUser, err := service.SubscriptionQueryHandler.GetAllUser(steamId)
	if err != nil {
		logger.Fail("Failed to get all tracked games for player: " + steamId + " | ERROR: " + err.Error())
		return nil
	}

	failedGames := service.SteamService.GetOnlyFailed(apiResponse, trackedGamesByUser)
	return failedGames
}
