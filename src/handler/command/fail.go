package command

import (
	"errors"
	"quit4real.today/logger"
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
	"time"
)

type FailsCommandHandler struct {
	FailRepository *repository.FailRepository
}

// Add adds a new failure record based on the subscription and the new time.
func (handler *FailsCommandHandler) Add(subscription model.Subscription, newTime int) error {
	if handler == nil {
		logger.Fail("FailsCommandHandler is nil")
		return errors.New("FailsCommandHandler is nil")
	}

	// Create a new GameFailureRecord based on the subscription
	failure := model.GameFailureRecord{
		DisplayName:     subscription.DisplayName,
		PlatformId:      subscription.PlatformId,
		PlatformGameId:  subscription.PlatformGameId,
		PlatformUserId:  subscription.PlatFormUserId,
		DurationMinutes: newTime - subscription.PlayedAmount, // Calculate the duration of the failure
		Reason:          "Game failed due to user action",    // Customize this reason as needed
		Timestamp:       time.Now().Format(time.RFC3339),     // Format the timestamp as needed
	}

	logger.Debug("Attempting to add a fail record for subscription: " + subscription.PlatformId)
	err := handler.FailRepository.Add(failure) // Add the failure record to the repository
	if err != nil {
		logger.Fail("Failed to add fail record: " + err.Error())
		return err
	}
	return nil
}
