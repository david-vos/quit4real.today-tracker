package repository

import (
	"project/config"
	"project/models"
)

type FailedRepoController struct {
	DbContr *DatabaseController
}

func (c *FailedRepoController) GetFailsForUser(steamId string) ([]models.Fail, error) {
	query := "SELECT * FROM failed_games WHERE steam_id = ?"
	rows, err := c.DbContr.FetchRowsWithClose(query, steamId)
	if err != nil {
		config.HandleError("failed to parse GetFailsForUser", err)
	}
	defer closeRows(rows)
	var fails []models.Fail
	if rows.Next() {
		fail, err := models.MapFail(rows)
		if err != nil {
			config.HandleError("failed to map GetFailsForUser", err)
			return nil, err
		}
		fails = append(fails, fail)
	}
	return fails, nil
}

func (c *FailedRepoController) CreateFailed(fail models.Fail) error {
	query := "INSERT INTO failed_games (steam_id, game_id, failed_time) VALUES (?, ?, ?)"
	err := c.DbContr.ExecuteQuery(query, fail.SteamId, fail.GameId, fail.FailedTime)
	return err
}
