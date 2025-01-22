package model

import "database/sql"

type Subscription struct {
	ID             int64
	UserId         string `json:"user_id"`          // Ensure this is a string
	PlatformId     string `json:"platform_id"`      // Ensure this is a string
	GameId         string `json:"platform_game_id"` // Ensure this is a string
	PlatFormUserId string `json:"platform_user_id"`
	PlayedAmount   int    `json:"played_amount"` // This can remain an int
}

func MapSubscription(rows *sql.Rows) (Subscription, error) {
	var subscription Subscription
	if err := rows.Scan(
		&subscription.UserId,
		&subscription.UserId,
		&subscription.GameId,
		&subscription.PlatformId,
		&subscription.PlatFormUserId,
		&subscription.PlayedAmount,
	); err != nil {
		return Subscription{}, err
	}
	return subscription, nil
}
