package repository

import (
	"database/sql"
	"fmt"
	"quit4real.today/logger"
	"quit4real.today/src/model"
)

type FailRepository struct {
	DatabaseImpl *DatabaseImpl
}

// Get retrieves all failure records for a specific user by their user ID.
func (repository *FailRepository) Get(userID string) ([]model.GameFailureRecord, error) {
	query := "SELECT * FROM game_failure_records WHERE user_id = ?"
	rows, err := repository.DatabaseImpl.FetchRows(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse GetFailsForUser: %w", err)
	}

	defer func(rows *sql.Rows) {
		err := closeRows(rows)
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
func (repository *FailRepository) GetTopLeaderBoard() ([]model.GameFailureRecord, error) {
	query := `
	WITH RankedFails AS (
		SELECT
			*,
			ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY timestamp DESC) AS rn
		FROM
			game_failure_records
	)
	SELECT
		id,
		user_id,
		game_id,
		duration_minutes,
		reason,
		timestamp
	FROM
		RankedFails
	WHERE
		rn = 1
	ORDER BY
		timestamp DESC;
	`
	rows, err := repository.DatabaseImpl.FetchRows(query)
	if err != nil {
		return []model.GameFailureRecord{}, fmt.Errorf("failed to parse GetFailsTopLeaderBoard: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := closeRows(rows)
		if err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}(rows)

	// Map the rows to a GameFailureRecord array
	var failures []model.GameFailureRecord
	for rows.Next() {
		failure, err := model.MapGameFailureRecord(rows)
		if err != nil {
			return []model.GameFailureRecord{}, fmt.Errorf("failed to map GetFails: %w", err)
		}
		failures = append(failures, failure)
	}
	return failures, nil
}

// Add inserts a new failure record into the database.
func (repository *FailRepository) Add(failure model.GameFailureRecord) error {
	query := "INSERT INTO game_failure_records (user_id, game_id, duration_minutes, reason, timestamp) VALUES (?, ?, ?, ?, ?)"
	err := repository.DatabaseImpl.ExecuteQuery(query, failure.User.ID, failure.Game.ID, failure.DurationMinutes, failure.Reason, failure.Timestamp)
	return err
}
