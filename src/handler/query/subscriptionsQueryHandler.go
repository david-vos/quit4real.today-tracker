package query

import (
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type SubscriptionQueryHandlerImpl struct {
	SubscriptionRepository repository.SubscriptionRepository
}

func (handler *SubscriptionQueryHandlerImpl) GetAllSteam() ([]model.Subscription, error) {
	subscriptions, err := handler.SubscriptionRepository.GetAllSteam()
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (handler *SubscriptionQueryHandlerImpl) GetAllView() ([]model.SubscriptionsView, error) {
	subscriptions, err := handler.SubscriptionRepository.GetAllView()
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}

func (handler *SubscriptionQueryHandlerImpl) GetAllUser(userId string) ([]model.Subscription, error) {
	subscriptions, err := handler.SubscriptionRepository.GetAllForUser(userId)
	if err != nil {
		return nil, err
	}
	return subscriptions, nil
}
