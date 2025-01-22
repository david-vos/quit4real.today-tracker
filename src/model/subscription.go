package model

import "database/sql"

type Subscription struct {
	ID           int    `json:"id"`
	UserId       string `json:"user_id"`
	GameId       string `json:"game_id"`
	PlatformId   string `json:"api_key"`
	PlayedAmount int    `json:"played_amount"`
}

func MapSubscription(rows *sql.Rows) (Subscription, error) {
	var subscription Subscription
	if err := rows.Scan(
		&subscription.ID,
		&subscription.UserId,
		&subscription.GameId,
		&subscription.PlatformId,
		&subscription.PlayedAmount,
	); err != nil {
		return Subscription{}, err
	}
	return subscription, nil
}
