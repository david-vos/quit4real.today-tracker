package command

import (
	"errors"
	"quit4real.today/logger"
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
	"time"
)

type FailsCommandHandler struct {
	FailRepository *repository.FailRepository
}

func (handler *FailsCommandHandler) Add(tracker model.Tracker, newTime int) error {
	if handler == nil {
		logger.Fail("FailsCommandHandler is nil")
		return errors.New("FailsCommandHandler is nil")
	}

	fail := model.Fail{
		SteamId:    tracker.SteamId,
		GameId:     tracker.GameId,
		FailedAt:   time.Now(),
		PlayedTime: newTime - tracker.PlayedAmount,
	}
	logger.Debug("Attempting to add a fail record for tracker: " + tracker.SteamId)
	err := handler.FailRepository.Add(fail)
	if err != nil {
		logger.Fail("Failed to add fail record: " + err.Error())
	}
	return err
}
