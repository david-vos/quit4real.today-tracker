package model

import "database/sql"

type Trakcer struct {
	ID                   int    `json:"id"`      // Auto-incrementing ID
	UserId               int    `json:"user_id"` // Internal User ID
	PlatformId           string `json:"platform_id"`
	PlatformGameId       string `json:"platform_game_id"`
	PlatformUserId       string `json:"platform_user_id"`
	DurationMinutes      int    `json:"duration_minutes"`
	AmountOfTimesStarted int    `json:"amount_of_times_started"`
	Day                  string `json:"day"`
	GameName             string `json:"game_name"`
}

func MapTrakcer(rows *sql.Rows) (Trakcer, error) {
	var tracker Trakcer
	if err := rows.Scan(
		&tracker.ID,
		&tracker.UserId,
		&tracker.PlatformId,
		&tracker.PlatformGameId,
		&tracker.PlatformUserId,
		&tracker.DurationMinutes,
		&tracker.AmountOfTimesStarted,
		&tracker.Day,
		&tracker.GameName,
	); err != nil {
		return tracker, err
	}
	return tracker, nil
}
