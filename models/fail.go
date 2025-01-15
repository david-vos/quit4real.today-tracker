package models

import (
	"database/sql"
	"time"
)

type Fail struct {
	ID         int       `json:"id"`
	SteamId    string    `json:"steam_id"`
	GameId     string    `json:"game_id"`
	FailedAt   time.Time `json:"failed_at"`
	PlayedTime int       `json:"played_time"`
}

func MapFail(rows *sql.Rows) (Fail, error) {
	var fail Fail
	if err := rows.Scan(
		&fail.ID,
		&fail.SteamId,
		&fail.GameId,
		&fail.FailedAt,
		&fail.PlayedTime,
	); err != nil {
		return Fail{}, err
	}
	return fail, nil
}
