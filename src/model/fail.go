package model

import (
	"database/sql"
)

// GameFailureRecord represents a failure record for a game.
type GameFailureRecord struct {
	ID              int    `json:"id"` // Auto-incrementing ID
	DisplayName     string `json:"display_name"`
	PlatformId      string `json:"platform_id"`
	PlatformGameId  string `json:"platform_game_id"`
	PlatformUserId  string `json:"platform_user_id"`
	DurationMinutes int    `json:"duration_minutes"`
	Reason          string `json:"reason"`
	Timestamp       string `json:"timestamp"`
}

// MapGameFailureRecord maps SQL rows to a GameFailureRecord struct.
func MapGameFailureRecord(rows *sql.Rows) (GameFailureRecord, error) {
	var failureRecord GameFailureRecord
	if err := rows.Scan(
		&failureRecord.ID,
		&failureRecord.DisplayName,
		&failureRecord.PlatformId,
		&failureRecord.PlatformGameId,
		&failureRecord.PlatformUserId,
		&failureRecord.DurationMinutes,
		&failureRecord.Reason,
		&failureRecord.Timestamp,
	); err != nil {
		return GameFailureRecord{}, err
	}
	return failureRecord, nil
}
