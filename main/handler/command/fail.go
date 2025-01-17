package command

import (
	"project/main/model"
	"project/main/repository"
	"time"
)

type FailsCommandHandler struct {
	failRepository *repository.FailRepository
}

func (handler *FailsCommandHandler) Add(tracker model.Tracker, newTime int) error {
	fail := model.Fail{
		SteamId:    tracker.SteamId,
		GameId:     tracker.GameId,
		FailedAt:   time.Now(),
		PlayedTime: newTime - tracker.PlayedAmount,
	}
	err := handler.failRepository.Add(fail)
	return err
}
