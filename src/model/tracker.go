package model

import (
	"database/sql"
	"time"
)

// Tracker represents the tracker table.
type Tracker struct {
	ID                 int       `db:"id"`
	UserID             int       `db:"user_id"`
	GameID             int       `db:"game_id"`
	PlatformID         string    `db:"platform_id"`
	TimePlayed         int       `db:"time_played"`
	NewTotalTimePlayed int       `db:"new_total_time_played"`
	AmountOfLogins     int       `db:"amount_of_logins"`
	Day                time.Time `db:"day"`
}

func MapTracker(rows *sql.Rows) (Tracker, error) {
	var tracker Tracker
	if err := rows.Scan(
		&tracker.ID,
		&tracker.UserID,
		&tracker.GameID,
		&tracker.PlatformID,
		&tracker.TimePlayed,
		&tracker.NewTotalTimePlayed,
		&tracker.AmountOfLogins,
		&tracker.Day,
	); err != nil {
		return tracker, err
	}
	return tracker, nil
}
