package repository

import "project/models"

type FailedRepoController struct {
	DbContr *DatabaseController
}

func (c *FailedRepoController) CreateFailed(fail models.Fail) error {
	query := "INSERT INTO failed_games (steam_id, game_id, failed_time) VALUES (?, ?, ?)"
	err := c.DbContr.ExecuteQuery(query, fail.SteamId, fail.GameId, fail.FailedTime)
	return err
}
