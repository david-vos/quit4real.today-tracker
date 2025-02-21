package impl

import (
	"database/sql"
	"fmt"
	"quit4real.today/logger"
	"quit4real.today/src/model"
)

type FailRepositoryImpl struct {
	DatabaseImpl *DatabaseImpl
}

// Get retrieves all failure records for a specific user by their user ID.
func (repository *FailRepositoryImpl) Get(userID string) ([]model.GameFailureRecord, error) {
	query := "SELECT * FROM game_failure_records WHERE platform_user_id = ?"
	rows, err := repository.DatabaseImpl.FetchRows(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GetFailsForUser: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := repository.DatabaseImpl.CloseRows(rows)
		if err != nil {
			logger.Debug("failed to close rows: " + err.Error())
		}
	}(rows)

	// Map the rows to GameFailureRecord array
	var failures []model.GameFailureRecord
	for rows.Next() {
		failure, err := model.MapGameFailureRecord(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to map GetFailsForUser: %w", err)
		}
		failures = append(failures, failure)
	}
	return failures, nil
}

// GetTopLeaderBoard retrieves the top failure records for the leaderboard.
func (repository *FailRepositoryImpl) GetTopLeaderBoard() ([]model.FailResponse, error) {
	query := `WITH RankedFails AS (
		SELECT
	gfr.id,
		gfr.display_name,
		gfr.platform_id,
		gfr.platform_game_id,
		gfr.platform_user_id,
		gfr.duration_minutes,
		gfr.reason,
		gfr.timestamp,
		g.name AS game_name,
	ROW_NUMBER() OVER (PARTITION BY gfr.platform_user_id, gfr.platform_game_id ORDER BY gfr.timestamp DESC) AS rn
	FROM
	game_failure_records gfr
	JOIN
	games g ON gfr.platform_game_id = g.id AND gfr.platform_id = g.platform_id
	)
	SELECT
	id,
		display_name,
		platform_id,
		platform_game_id,
		platform_user_id,
		duration_minutes,
		reason,
		timestamp,
		game_name
	FROM
	RankedFails
	WHERE
	rn = 1 
	ORDER BY
	timestamp DESC;`

	rows, err := repository.DatabaseImpl.FetchRows(query)
	if err != nil {
		return []model.FailResponse{}, fmt.Errorf("failed to parse GetFailsTopLeaderBoard: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := repository.DatabaseImpl.CloseRows(rows)
		if err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}(rows)

	// Map the rows to a FailResponse array
	var failures []model.FailResponse
	for rows.Next() {
		failure, err := model.MapFailResponse(rows)
		if err != nil {
			return []model.FailResponse{}, fmt.Errorf("failed to map GetFails: %w", err)
		}
		failures = append(failures, failure)
	}
	return failures, nil
}

// Add inserts a new failure record into the database.
func (repository *FailRepositoryImpl) Add(failure model.GameFailureRecord) error {
	query := "INSERT INTO game_failure_records (display_name, platform_id, platform_game_id, platform_user_id, duration_minutes, reason, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?)"
	err := repository.DatabaseImpl.ExecuteQuery(query, failure.DisplayName, failure.PlatformId, failure.PlatformGameId, failure.PlatformUserId, failure.DurationMinutes, failure.Reason, failure.Timestamp)
	return err
}
