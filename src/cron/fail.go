package cron

import (
	"github.com/robfig/cron/v3"
	"quit4real.today/logger"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
)

type FailCron struct {
	SubscriptionCommandHandler *command.SubscriptionCommandHandler
	SubscriptionQueryHandler   *query.SubscriptionQueryHandler
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
	//users, err := fc.UserQueryHandler.GetAll()
	allSteamSubscriptions, err := fc.SubscriptionQueryHandler.GetAllSteam()
	if err != nil {
		logger.Fail("Error getting user when running cron jobs: " + err.Error())
		return
	}
	for _, subscription := range allSteamSubscriptions {
		fc.SubscriptionCommandHandler.UpdateFromSteamApi(subscription.PlatFormUserId)
	}
}
