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

func (c *FailedRepoController) GetFailsTopLeaderBoard() ([]models.Fail, error) {
	query := `
	WITH RankedFails AS (
		SELECT
			*,
			ROW_NUMBER() OVER (PARTITION BY steam_id ORDER BY failed_at DESC) AS rn
		FROM
			failed_games
	)
	SELECT
		id,
		steam_id,
		game_id,
		failed_at,
		played_time
	FROM
		RankedFails
	WHERE
		rn = 1
	ORDER BY
		failed_at DESC;
	`
	rows, err := c.DbContr.FetchRowsWithClose(query)
	if err != nil {
		config.HandleError("failed to parse GetFails", err)
		return []models.Fail{}, nil
	}
	defer closeRows(rows)
	var fails []models.Fail
	for rows.Next() {
		fail, err := models.MapFail(rows)
		if err != nil {
			config.HandleError("failed to map GetFails", err)
			return []models.Fail{}, nil
		}
		fails = append(fails, fail)
	}
	return fails, nil
}

func (c *FailedRepoController) CreateFailed(fail models.Fail) error {
	query := "INSERT INTO failed_games (steam_id, game_id, failed_at, played_time) VALUES (?, ?, ?, ?)"
	err := c.DbContr.ExecuteQuery(query, fail.SteamId, fail.GameId, fail.FailedAt, fail.PlayedTime)
	return err
}
