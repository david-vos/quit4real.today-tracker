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

func (handler *FailsCommandHandler) Add(subscription model.Subscription, newTime int) error {
	if handler == nil {
		logger.Fail("FailsCommandHandler is nil")
		return errors.New("FailsCommandHandler is nil")
	}

	fail := model.Fail{
		GameId:     subscription.GameId,
		FailedAt:   time.Now(),
		PlayedTime: newTime - subscription.PlayedAmount,
	}
	logger.Debug("Attempting to add a fail record for subscription: " + subscription.SteamId)
	err := handler.FailRepository.Add(fail)
	if err != nil {
		logger.Fail("Failed to add fail record: " + err.Error())
	}
	return err
}
