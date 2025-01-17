package command

import (
	"project/main/repository"
)

type TrackerCommandHandler struct {
	trackerRepository *repository.TrackerRepository
}

func (handler *TrackerCommandHandler) Add(platformId string, gameId string) error {
	var err = handler.trackerRepository.Add(platformId, gameId)
	if err != nil {
		return err
	}
	return nil
}
