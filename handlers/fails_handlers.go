package handlers

import (
	"project/models"
	"project/repository"
)

type FailsController struct {
	FailRepoContr repository.FailedRepoController
	UserRepoContr repository.UserRepoController
}

func (c *FailsController) GetFails(steamId string) ([]models.Fail, error) {
	fails, err := c.FailRepoContr.GetFailsForUser(steamId)
	if err != nil {
		return nil, err
	}
	return fails, nil
}
