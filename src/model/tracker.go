package model

import "database/sql"

type Tracker struct {
	ID           int    `json:"id"`
	SteamId      string `json:"steam_id"`
	GameId       string `json:"api_key"`
	PlayedAmount int    `json:"played_amount"`
}

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
