package query

import (
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type FailQueryHandler struct {
	FailRepository *repository.FailRepository
}

// Get retrieves all failure records for a specific user by their user ID.
func (handler *FailQueryHandler) Get(userID string) ([]model.GameFailureRecord, error) {
	failures, err := handler.FailRepository.Get(userID)
	if err != nil {
		return nil, err
	}
	return failures, nil
}

// GetLeaderBoard retrieves the top failure records for the leaderboard.
func (handler *FailQueryHandler) GetLeaderBoard() ([]model.GameFailureRecord, error) {
	failuresLeaderBoard, err := handler.FailRepository.GetTopLeaderBoard()
	if err != nil {
		return nil, err
	}
	return failuresLeaderBoard, nil
}
