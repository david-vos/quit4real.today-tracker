package handlers

import (
	"project/models"
	"project/repository"
	"time"
)

type FailsController struct {
	FailRepoContr *repository.FailedRepoController
	UserRepoContr *repository.UserRepoController
}

func (c *FailsController) GetFails(steamId string) ([]models.Fail, error) {
	fails, err := c.FailRepoContr.GetFailsForUser(steamId)
	if err != nil {
		return nil, err
	}
	return fails, nil
}

func (c *FailsController) createFail(tracker models.Tracker, newTime int) error {
	fail := models.Fail{
		SteamId:    tracker.SteamId,
		GameId:     tracker.GameId,
		FailedAt:   time.Now(),
		PlayedTime: newTime - tracker.PlayedAmount,
	}
	err := c.FailRepoContr.CreateFailed(fail)
	return err
}
