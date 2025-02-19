package cron

import (
	"github.com/robfig/cron/v3"
	"quit4real.today/logger"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
	"quit4real.today/src/handler/service"
)

type FailCron struct {
	SubscriptionCommandHandler *command.SubscriptionCommandHandler
	SubscriptionQueryHandler   *query.SubscriptionQueryHandler
	UserQueryHandler           *query.UserQueryHandler
	SteamService               *service.SteamService
}

func (fc *FailCron) Start() {
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

func (fc *FailCron) updateAndSendNotifySteam() {
	allSteamSubscriptions, err := fc.SubscriptionQueryHandler.GetAllSteam()
	if err != nil {
		logger.Fail("Error getting user when running cron jobs: " + err.Error())
		return
	}
	allSteamUsers, err := fc.UserQueryHandler.GetAll()
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
	for platformUserId := range platformUserIds {
		fc.SteamService.UpdateFromSteamApi(platformUserId)
	}
}
