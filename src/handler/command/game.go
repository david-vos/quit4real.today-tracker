package command

import (
	"quit4real.today/src/repository"
)

type GameCommandHandler struct {
	GameRepository *repository.GameRepository
}

func (handler *GameCommandHandler) Add(id string, name string, platformId string) error {
	if handler.GameRepository.Exists(id, platformId) {
		return nil
	}
	return handler.GameRepository.Add(id, name, platformId)
}
