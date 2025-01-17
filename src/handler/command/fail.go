package command

import (
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
	"time"
)

type FailsCommandHandler struct {
	FailRepository *repository.FailRepository
}

func (handler *FailsCommandHandler) Add(tracker model.Tracker, newTime int) error {
	fail := model.Fail{
		SteamId:    tracker.SteamId,
		GameId:     tracker.GameId,
		FailedAt:   time.Now(),
		PlayedTime: newTime - tracker.PlayedAmount,
	}
	err := handler.FailRepository.Add(fail)
	return err
}
