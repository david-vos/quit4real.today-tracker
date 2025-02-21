package repository

import "quit4real.today/src/model"

type SubscriptionRepository interface {
	Add(displayName string, platformId string, platformGameId string, platformUserId string, playedAmount int) error
	Update(userId string, gameId string, playedAmount int) error
	Get(userId string, gameId string) (model.Subscription, error)
	GetAllForUser(userId string) ([]model.Subscription, error)
	GetAllSteam() ([]model.Subscription, error)
	GetAllView() ([]model.SubscriptionsView, error)
}
