package query

import (
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type SubscriptionQueryHandler struct {
	SubscriptionRepository *repository.SubscriptionRepository
}

func (handler *SubscriptionQueryHandler) GetAllSteam() ([]model.Subscription, error) {
	subscriptions, err := handler.SubscriptionRepository.GetAllSteam()
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (handler *SubscriptionQueryHandler) GetAllView() ([]model.SubscriptionsView, error) {
	subscriptions, err := handler.SubscriptionRepository.GetAllView()
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}
