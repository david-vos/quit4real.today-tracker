package repository

import (
	"database/sql"
	"quit4real.today/logger"
	"quit4real.today/src/model"
)

type SubscriptionRepository struct {
	DatabaseImpl *DatabaseImpl
}

// Add inserts a new subscription for a user into the database.
func (repository *SubscriptionRepository) Add(userId string, gameId string, playedAmount int) error {
	query := "INSERT INTO user_platform_subscriptions (user_id, game_id, played_amount) VALUES (?, ?, ?);"
	return repository.DatabaseImpl.ExecuteQuery(query, userId, gameId, playedAmount)
}

// Update updates the played amount for a specific game subscription for a user.
func (repository *SubscriptionRepository) Update(userId string, gameId string, playedAmount int) error {
	query := "UPDATE user_platform_subscriptions SET played_amount = ? WHERE user_id = ? AND game_id = ?;"
	return repository.DatabaseImpl.ExecuteQuery(query, playedAmount, userId, gameId)
}

// Get retrieves a specific subscription for a user and game.
func (repository *SubscriptionRepository) Get(userId string, gameId string) (model.Subscription, error) {
	query := "SELECT * FROM user_platform_subscriptions WHERE user_id = ? AND game_id = ?;"
	rows, err := repository.DatabaseImpl.FetchRows(query, userId, gameId)
	if err != nil {
		return model.Subscription{}, err
	}
	defer func(rows *sql.Rows) {
		if err := closeRows(rows); err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}(rows)

	if rows.Next() {
		return model.MapSubscription(rows)
	}
	return model.Subscription{}, nil
}

// GetAll retrieves all subscriptions for a specific user.
func (repository *SubscriptionRepository) GetAll(userId string) ([]model.Subscription, error) {
	query := "SELECT * FROM user_platform_subscriptions WHERE user_id = ?;"
	rows, err := repository.DatabaseImpl.FetchRows(query, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := closeRows(rows); err != nil {
			logger.Fail("failed to close rows: " + err.Error())
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
