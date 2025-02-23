package cronJobs

import (
	"github.com/robfig/cron/v3"
	"quit4real.today/logger"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
	"quit4real.today/src/handler/service"
)

type FailCronImpl struct {
	SubscriptionCommandHandler *command.SubscriptionCommandHandlerImpl
	SubscriptionQueryHandler   *query.SubscriptionQueryHandlerImpl
	UserQueryHandler           *query.UserQueryHandlerImpl
	SteamService               service.SteamService
	SubscriptionService        service.SubscriptionService
}

func (fc *FailCronImpl) Start() {
	// If in the future I will add another cron job then it's better to have a simulacra structure to how the Endpoints work
	cronJob := cron.New()
	// ((24*60)/10)*694 ~= 100.000 the STEAM API limit
	// 694 -> max amount of users :thinking per API key
	// Currently set to 1 minute for dev, I will increase this back to 10 minutes when going live
	_, err := cronJob.AddFunc("@every 1m", func() {
		fc.updateAndSendNotifySteam()
	})
	if err != nil {
		logger.Fail("Error adding cron job: " + err.Error())
		return
	}

	cronJob.Start()
}

func (fc *FailCronImpl) updateAndSendNotifySteam() {
	allSteamSubscriptions, err := fc.SubscriptionQueryHandler.GetAllSteam()
	if err != nil {
		logger.Fail("Error getting user when running cron jobs: " + err.Error())
		return
	}
	allSteamUsers, err := fc.UserQueryHandler.GetAllSteamVerified()
	if err != nil {
		logger.Fail("Error getting user when running cron jobs: " + err.Error())
	}

	platformUserIds := make(map[string]bool)

	// Populate the map with all Steam subscriptions
	for _, subscription := range allSteamSubscriptions {
		platformUserId := subscription.PlatFormUserId
		platformUserIds[platformUserId] = true
	}

	// Populate the map with all Steam subscriptions
	for _, user := range allSteamUsers {
		userId := user.SteamID
		if !platformUserIds[userId] {
			platformUserIds[userId] = true
		}
	}

	// Iterate over each unique platformUserId and update
	// This should all be abstracted away so you don't need a different service with a different platform
	for platformUserId := range platformUserIds {
		println("Updating Steam subscription for user: " + platformUserId)
		apiResponse, err := fc.SteamService.FetchRecentGames(platformUserId)
		if err != nil {
			logger.Fail("failed to fetch player information for player: " + platformUserId + " | ERROR: " + err.Error())
			return
		}
		failedGamesByUser := fc.SubscriptionService.UpdateSteamSubscription(platformUserId, apiResponse)
		fc.SubscriptionCommandHandler.UpdateSubscriptions(platformUserId, failedGamesByUser)
	}
}
