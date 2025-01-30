package model

import "database/sql"

type Subscription struct {
	ID             int64
	DisplayName    string `json:"display_name"` // Ensure this is a string
	PlatformId     string `json:"platform_id"`  // Ensure this is a string
	PlatformGameId string `json:"platform_game_id"`
	PlatFormUserId string `json:"platform_user_id"`
	PlayedAmount   int    `json:"played_amount"` // This can remain an int
}

type SubscriptionsView struct {
	ID             int64
	DisplayName    string `json:"display_name"` // Ensure this is a string
	PlatformId     string `json:"platform_id"`  // Ensure this is a string
	PlatformGameId string `json:"platform_game_id"`
	PlatFormUserId string `json:"platform_user_id"`
	PlayedAmount   int    `json:"played_amount"` // This can remain an int
	GameName       string `json:"game_name"`
}

func MapSubscription(rows *sql.Rows) (Subscription, error) {
	var subscription Subscription
	if err := rows.Scan(
		&subscription.ID,
		&subscription.DisplayName,
		&subscription.PlatformId,
		&subscription.PlatformGameId,
		&subscription.PlatFormUserId,
		&subscription.PlayedAmount,
	); err != nil {
		return Subscription{}, err
	}
	return subscription, nil
}

func MapSubscriptionsView(rows *sql.Rows) (SubscriptionsView, error) {
	var subscriptionsView SubscriptionsView
	if err := rows.Scan(
		&subscriptionsView.ID,
		&subscriptionsView.DisplayName,
		&subscriptionsView.PlatformId,
		&subscriptionsView.PlatformGameId,
		&subscriptionsView.PlatFormUserId,
		&subscriptionsView.PlayedAmount,
		&subscriptionsView.GameName,
	); err != nil {
		return SubscriptionsView{}, err
	}
	return subscriptionsView, nil
}
