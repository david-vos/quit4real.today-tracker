package repository

import (
	"database/sql"
	"quit4real.today/logger"
	"quit4real.today/src/model"
)

type SubscriptionRepository struct {
	DatabaseImpl *DatabaseImpl
}

func (repository *SubscriptionRepository) Add(steamId string, gameId string, playerAmount int) error {
	query := "INSERT INTO user_platform_subscriptions (steam_id, game_id, played_amount) VALUES (?, ?, ?);"
	return repository.DatabaseImpl.ExecuteQuery(query, steamId, gameId, playerAmount)
}

func (repository *SubscriptionRepository) Update(steamId string, gameId string, playedAmount int) error {
	query := "UPDATE user_platform_subscriptions SET played_amount = ? WHERE steam_id = ? AND game_id = ?;"
	return repository.DatabaseImpl.ExecuteQuery(query, playedAmount, steamId, gameId)
}

func (repository *SubscriptionRepository) Get(steamId string, gameId string) (model.Subscription, error) {
	query := "SELECT * FROM user_platform_subscriptions WHERE steam_id = ? AND game_id = ?;"
	rows, err := repository.DatabaseImpl.FetchRows(query, steamId, gameId)
	if err != nil {
		return model.Subscription{}, err
	}
	defer func(rows *sql.Rows) {
		err := closeRows(rows)
		if err != nil {
			logger.Fail("failed to close rows" + err.Error())
		}
	}(rows)

	if rows.Next() {
		return model.MapSubscription(rows)
	}
	return model.Subscription{}, nil
}

func (repository *SubscriptionRepository) GetAll(steamId string) ([]model.Subscription, error) {
	query := "SELECT * FROM user_platform_subscriptions WHERE steam_id = ?;"
	rows, err := repository.DatabaseImpl.FetchRows(query, steamId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := closeRows(rows)
		if err != nil {
			logger.Fail("failed to close rows" + err.Error())
		}
	}(rows)

	var subscriptions []model.Subscription
	for rows.Next() {
		subscription, err := model.MapSubscription(rows)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}
