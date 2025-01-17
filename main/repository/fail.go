package repository

import (
	"database/sql"
	"fmt"
	"project/logger"
	"project/main/model"
)

type FailRepository struct {
	DatabaseImpl *DatabaseImpl
}

func (repository *FailRepository) Get(steamId string) ([]model.Fail, error) {
	query := "SELECT * FROM failed_games WHERE steam_id = ?"
	rows, err := repository.DatabaseImpl.FetchRows(query, steamId)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GetFailsForUser: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := closeRows(rows)
		if err != nil {
			logger.Debug("failed to close rows: " + err.Error())
		}
	}(rows)

	// Map the rows to Fail array
	var fails []model.Fail
	if rows.Next() {
		fail, err := model.MapFail(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to map GetFailsForUser: %w", err)
		}
		fails = append(fails, fail)
	}
	return fails, nil
}

func (repository *FailRepository) GetTopLeaderBoard() ([]model.Fail, error) {
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
	rows, err := repository.DatabaseImpl.FetchRows(query)
	if err != nil {
		return []model.Fail{}, fmt.Errorf("failed to parse GetFailsTopLeaderBoard: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := closeRows(rows)
		if err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}(rows)

	// Map the row to a Fail array
	var fails []model.Fail
	for rows.Next() {
		fail, err := model.MapFail(rows)
		if err != nil {
			return []model.Fail{}, fmt.Errorf("failed to map GetFails: %w", err)
		}
		fails = append(fails, fail)
	}
	return fails, nil
}

func (repository *FailRepository) Add(fail model.Fail) error {
	query := "INSERT INTO failed_games (steam_id, game_id, failed_at, played_time) VALUES (?, ?, ?, ?)"
	err := repository.DatabaseImpl.ExecuteQuery(query, fail.SteamId, fail.GameId, fail.FailedAt, fail.PlayedTime)
	return err
}
