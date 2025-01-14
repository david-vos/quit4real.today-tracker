package models

import (
	"database/sql"
	"time"
)

type Fail struct {
	ID         int       `json:"id"`
	SteamId    string    `json:"steam_id"`
	GameId     string    `json:"game_id"`
	FailedTime time.Time `json:"failed_time"`
}

func MapFail(rows *sql.Rows) (Fail, error) {
	var fail Fail
	if err := rows.Scan(
		&fail.ID,
		&fail.SteamId,
		&fail.GameId,
		&fail.FailedTime,
	); err != nil {
		return Fail{}, err
	}
	return fail, nil
}
