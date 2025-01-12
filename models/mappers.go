package models

import (
	"database/sql"
)

func MapTracker(rows *sql.Rows) (Tracker, error) {
	var tracker Tracker
	if err := rows.Scan(
		&tracker.ID,
		&tracker.SteamId,
		&tracker.GameId,
		&tracker.PlayedAmount,
	); err != nil {
		return Tracker{}, err
	}
	return tracker, nil
}

func MapUser(rows *sql.Rows) (User, error) {
	var user User
	if err := rows.Scan(
		&user.ID,
		&user.Name,
		&user.SteamId,
		&user.ApiKey,
	); err != nil {
		return User{}, err
	}
	return user, nil
}
