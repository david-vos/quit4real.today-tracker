package query

import (
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type FailQueryHandlerImpl struct {
	FailRepository repository.FailRepository
}

// GetLeaderBoard retrieves the top failure records for the leaderboard.
func (handler *FailQueryHandlerImpl) GetLeaderBoard() ([]model.FailResponse, error) {
	failuresLeaderBoard, err := handler.FailRepository.GetTopLeaderBoard()
	if err != nil {
		return nil, err
	}
	return failuresLeaderBoard, nil
}
