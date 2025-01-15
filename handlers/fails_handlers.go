package handlers

import (
	"encoding/json"
	"net/http"
	"project/models"
	"project/repository"
	"time"
)

type FailsController struct {
	FailRepoContr *repository.FailedRepoController
	UserRepoContr *repository.UserRepoController
}

func (c *FailsController) GetFailsForUser(steamId string) ([]models.Fail, error) {
	fails, err := c.FailRepoContr.GetFailsForUser(steamId)
	if err != nil {
		return nil, err
	}
	return fails, nil
}

func (c *FailsController) GetFailsLeaderBoard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		failsLeaderBoard, err := c.FailRepoContr.GetFailsTopLeaderBoard()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return // Ensure to return after writing the header
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(failsLeaderBoard); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
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
