package query

import (
	"project/main/model"
	"project/main/repository"
)

type FailQueryHandler struct {
	FailRepository *repository.FailRepository
}

func (handler *FailQueryHandler) Get(steamId string) ([]model.Fail, error) {
	fails, err := handler.FailRepository.Get(steamId)
	if err != nil {
		return nil, err
	}
	return fails, nil
}

func (handler *FailQueryHandler) GetLeaderBoard() ([]model.Fail, error) {
	failsLeaderBoard, err := handler.FailRepository.GetTopLeaderBoard()
	if err != nil {
		return nil, err
	}
	return failsLeaderBoard, nil
}
