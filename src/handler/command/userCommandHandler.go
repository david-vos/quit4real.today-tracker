package command

import (
	"quit4real.today/src/model"
	"quit4real.today/src/repository"
)

type UserCommandHandlerImpl struct {
	UserRepository repository.UserRepository
}

func (handler *UserCommandHandlerImpl) Add(user model.User) error {
	var err = handler.UserRepository.Add(user)
	if err != nil {
		return err
	}
	return nil
}

func (handler *UserCommandHandlerImpl) Update(user model.User) error {
	var err = handler.UserRepository.Update(user)
	if err != nil {
		return err
	}
	return nil
}
