package model

import "database/sql"

type Subscription struct {
	UserId       string `json:"user_id"`       // Ensure this is a string
	GameId       string `json:"game_id"`       // Ensure this is a string
	PlatformId   string `json:"platform_id"`   // Ensure this is a string
	PlayedAmount int    `json:"played_amount"` // This can remain an int
}

func MapSubscription(rows *sql.Rows) (Subscription, error) {
	var subscription Subscription
	if err := rows.Scan(
		&subscription.UserId,
		&subscription.UserId,
		&subscription.GameId,
		&subscription.PlatformId,
		&subscription.PlayedAmount,
	); err != nil {
		return Subscription{}, err
	}
	return subscription, nil
}
