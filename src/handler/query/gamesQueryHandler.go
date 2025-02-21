package query

import (
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type GameQueryHandlerImpl struct {
	GameRepository repository.GameRepository
}

func (handler *GameQueryHandlerImpl) Search(searchParam string, platformId string) ([]model.Game, error) {
	games, err := handler.GameRepository.Search(searchParam, platformId)
	if err != nil {
		return []model.Game{}, err
	}
	return games, nil
}
