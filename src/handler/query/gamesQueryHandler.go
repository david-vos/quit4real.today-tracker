package query

import (
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type GameQueryHandler struct {
	GameRepository *repository.GameRepository
}

func (handler *GameQueryHandler) Search(searchParam string, platformId string) ([]model.Game, error) {
	games, err := handler.GameRepository.Search(searchParam, platformId)
	if err != nil {
		return []model.Game{}, err
	}
	return games, nil
}
