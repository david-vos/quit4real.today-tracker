package cron

import (
	"github.com/robfig/cron/v3"
	"quit4real.today/logger"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
)

type FailCron struct {
	UserQueryHandler      *query.UserQueryHandler
	TrackerCommandHandler *command.TrackerCommandHandler
}

func (fc *FailCron) Start() {
	// If in the future I will add another cron job then it's better to have a simulacra structure to how the Endpoints work
	cronJob := cron.New()
	// ((24*60)/10)*694 ~= 100.000 the STEAM API limit
	// 694 -> max amount of users :thinking per API key
	_, err := cronJob.AddFunc("@every 1m", func() {
		fc.updateAndSendNotify()
	})
	if err != nil {
		logger.Fail("Error adding cron job: " + err.Error())
		return
	}

	cronJob.Start()
}

func (fc *FailCron) updateAndSendNotify() {
	users, err := fc.UserQueryHandler.GetAll()
	if err != nil {
		logger.Fail("Error getting user when running cron jobs: " + err.Error())
		return
	}
	for _, user := range users {
		fc.TrackerCommandHandler.UpdateFromSteamApi(user.SteamId)
	}
}
