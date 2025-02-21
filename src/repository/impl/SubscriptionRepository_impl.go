package impl

import (
	"database/sql"
	"quit4real.today/logger"
	"quit4real.today/src/model"
)

type SubscriptionRepositoryImpl struct {
	DatabaseImpl *DatabaseImpl
}

// Add inserts a new subscription for a user into the database.
func (repository *SubscriptionRepositoryImpl) Add(displayName string, platformId string, platformGameId string, platformUserId string, playedAmount int) error {
	query := "INSERT INTO platform_subscriptions (display_name, platform_id, platform_game_id, platform_user_id, played_amount) VALUES (?, ?, ?, ?, ?);"
	return repository.DatabaseImpl.ExecuteQuery(query, displayName, platformId, platformGameId, platformUserId, playedAmount)
}

// Update updates the played amount for a specific game subscription for a user.
func (repository *SubscriptionRepositoryImpl) Update(userId string, gameId string, playedAmount int) error {
	query := "UPDATE platform_subscriptions SET played_amount = ? WHERE platform_user_id = ? AND platform_game_id = ?;"
	return repository.DatabaseImpl.ExecuteQuery(query, playedAmount, userId, gameId)
}

// Get retrieves a specific subscription for a user and game.
func (repository *SubscriptionRepositoryImpl) Get(userId string, gameId string) (model.Subscription, error) {
	query := "SELECT * FROM platform_subscriptions WHERE platform_user_id = ? AND platform_game_id = ?;"
	rows, err := repository.DatabaseImpl.FetchRows(query, userId, gameId)
	if err != nil {
		return model.Subscription{}, err
	}
	defer func(rows *sql.Rows) {
		if err := repository.DatabaseImpl.CloseRows(rows); err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}(rows)

	if rows.Next() {
		return model.MapSubscription(rows)
	}
	return model.Subscription{}, nil
}

// GetAllForUser retrieves all subscriptions for a specific user.
func (repository *SubscriptionRepositoryImpl) GetAllForUser(userId string) ([]model.Subscription, error) {
	query := "SELECT * FROM platform_subscriptions WHERE platform_user_id = ?;"
	rows, err := repository.DatabaseImpl.FetchRows(query, userId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := repository.DatabaseImpl.CloseRows(rows); err != nil {
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

func (repository *SubscriptionRepositoryImpl) GetAllSteam() ([]model.Subscription, error) {
	query := "SELECT * FROM platform_subscriptions WHERE platform_id = ?;"
	rows, err := repository.DatabaseImpl.FetchRows(query, "steam")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := repository.DatabaseImpl.CloseRows(rows); err != nil {
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

func (repository *SubscriptionRepositoryImpl) GetAllView() ([]model.SubscriptionsView, error) {
	query := "SELECT PS.*, g.name FROM platform_subscriptions as PS LEFT JOIN games g on g.id = PS.platform_game_id"
	rows, err := repository.DatabaseImpl.FetchRows(query)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		if err := repository.DatabaseImpl.CloseRows(rows); err != nil {
			logger.Fail("failed to close rows: " + err.Error())
		}
	}(rows)

	var subscriptions []model.SubscriptionsView
	for rows.Next() {
		subscription, err := model.MapSubscriptionsView(rows)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}
