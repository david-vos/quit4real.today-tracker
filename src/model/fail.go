package model

import (
	"database/sql"
)

// GameFailureRecord represents a failure record for a game.
type GameFailureRecord struct {
	ID              int    `json:"id"`               // Auto-incrementing ID
	User            User   `json:"user"`             // Embedded User object
	Game            Game   `json:"game"`             // Embedded Game object
	DurationMinutes int    `json:"duration_minutes"` // Duration of the failure in minutes
	Reason          string `json:"reason"`           // Reason for the failure
	Timestamp       string `json:"timestamp"`        // When the failure was recorded
}

// MapGameFailureRecord maps SQL rows to a GameFailureRecord struct.
func MapGameFailureRecord(rows *sql.Rows) (GameFailureRecord, error) {
	var failureRecord GameFailureRecord
	var gameID string
	if err := rows.Scan(
		&failureRecord.ID,
		&failureRecord.User.ID,
		&gameID,
		&failureRecord.DurationMinutes,
		&failureRecord.Reason,
		&failureRecord.Timestamp,
	); err != nil {
		return GameFailureRecord{}, err
	}

	// Create a Game object based on the game ID
	failureRecord.Game = Game{ID: gameID}
	return failureRecord, nil
}
